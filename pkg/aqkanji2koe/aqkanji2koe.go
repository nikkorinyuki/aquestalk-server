package aqkanji2koe

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <dlfcn.h>

// 関数ポインタの型定義
typedef void* (*type_AqKanji2Koe_Create)(const char*, int*);
typedef void (*type_AqKanji2Koe_Release)(void*);
typedef int (*type_AqKanji2Koe_SetDevKey)(const char*);
typedef int (*type_AqKanji2Koe_Convert)(void*, const char*, char*, int);
typedef int (*type_AqKanji2Koe_ConvRoman)(void*, const char*, char*, int);

// ヘルパー関数: 関数ポインタ経由で呼び出すためのラッパー
void* call_Create(void* f, const char* path, int* err) {
    return ((type_AqKanji2Koe_Create)f)(path, err);
}
void call_Release(void* f, void* handle) {
    ((type_AqKanji2Koe_Release)f)(handle);
}
int call_SetDevKey(void* f, const char* key) {
    return ((type_AqKanji2Koe_SetDevKey)f)(key);
}
int call_Convert(void* f, void* handle, const char* text, char* buf, int size) {
    return ((type_AqKanji2Koe_Convert)f)(handle, text, buf, size);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	ErrNone           = 0   // 成功
	ErrArgument       = 101 // 関数呼び出し時の引数が NULL になっている。
	ErrNotInitialized = 104 // 初期化されていない(初期化ルーチンが呼ばれていない)
	ErrTextTooLong    = 105 // 入力テキストが長すぎる
	ErrDicUnload      = 106 // システム辞書データが指定されていない
	ErrInvalidText    = 107 // 変換できない文字コードが含まれている
	ErrProcessing     = 100 // その他のエラー
)

// 最小バッファサイズ
const minBufSize = 256

type AqKanji2Koe struct {
	libHandle  unsafe.Pointer
	handle     unsafe.Pointer
	fCreate    unsafe.Pointer
	fRelease   unsafe.Pointer
	fSetDevKey unsafe.Pointer
	fConvert   unsafe.Pointer
	fConvRoman unsafe.Pointer
}

// New はライブラリを動的にロードして初期化します
func New(libPath string, dicPath string) (*AqKanji2Koe, error) {
	cLibPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cLibPath))

	// 1. ライブラリをオープン
	lib := C.dlopen(cLibPath, C.RTLD_LAZY)
	if lib == nil {
		return nil, fmt.Errorf("dlopen failed: %s", C.GoString(C.dlerror()))
	}

	a := &AqKanji2Koe{libHandle: lib}

	// 2. シンボルの取得
	symbols := []struct {
		name string
		ptr  *unsafe.Pointer
	}{
		{"AqKanji2Koe_Create", &a.fCreate},
		{"AqKanji2Koe_Release", &a.fRelease},
		{"AqKanji2Koe_SetDevKey", &a.fSetDevKey},
		{"AqKanji2Koe_Convert", &a.fConvert},
		{"AqKanji2Koe_ConvRoman", &a.fConvRoman},
	}

	for _, s := range symbols {
		cName := C.CString(s.name)
		*s.ptr = C.dlsym(lib, cName)
		C.free(unsafe.Pointer(cName))
		if *s.ptr == nil {
			C.dlclose(lib)
			return nil, fmt.Errorf("symbol not found: %s", s.name)
		}
	}

	// 3. インスタンス生成
	cDicPath := C.CString(dicPath)
	defer C.free(unsafe.Pointer(cDicPath))
	var cErr C.int

	a.handle = C.call_Create(a.fCreate, cDicPath, &cErr)
	if a.handle == nil {
		C.dlclose(lib)
		return nil, fmt.Errorf("aqk2k_create failed (code: %d)", cErr)
	}

	return a, nil
}

func (a *AqKanji2Koe) Close() error {
	if a.handle != nil {
		C.call_Release(a.fRelease, a.handle)
		a.handle = nil
	}
	if a.libHandle != nil {
		C.dlclose(a.libHandle)
		a.libHandle = nil
	}
	return nil
}

func (a *AqKanji2Koe) SetDevKey(devKey string) error {
	pDevKey := C.CString(devKey)
	defer C.free(unsafe.Pointer(pDevKey))
	ret := C.call_SetDevKey(a.fSetDevKey, pDevKey)
	if ret == 1 {
		return fmt.Errorf("invalid devKey")
	}
	return nil
}

func (a *AqKanji2Koe) Convert(text string) (string, error) {
	return a.callConvert(text, false)
}

func (a *AqKanji2Koe) ConvertRoman(text string) (string, error) {
	return a.callConvert(text, true)
}

func (a *AqKanji2Koe) callConvert(text string, roman bool) (string, error) {
	pText := C.CString(text)
	defer C.free(unsafe.Pointer(pText))

	bufSize := len(text) * 2
	if bufSize < minBufSize {
		bufSize = minBufSize
	}
	buf := make([]byte, bufSize)

	var ret C.int
	fTarget := a.fConvert
	if roman {
		fTarget = a.fConvRoman
	}

	ret = C.call_Convert(fTarget, a.handle, pText, (*C.char)(unsafe.Pointer(&buf[0])), C.int(bufSize))

	if ret != ErrNone {
		return "", fmt.Errorf("convert failed (code: %d)", ret)
	}

	return C.GoString((*C.char)(unsafe.Pointer(&buf[0]))), nil
}
