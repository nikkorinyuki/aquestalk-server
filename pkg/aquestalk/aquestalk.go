package aquestalk

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>

// 関数ポインタの型定義
typedef unsigned char* (*SyntheFunc)(const char*, int, int*);
typedef void (*FreeWaveFunc)(unsigned char*);

// C側で動的に関数を呼び出すためのヘルパー
unsigned char* bridge_synthe(void* f, const char* koe, int speed, int* size) {
    return ((SyntheFunc)f)(koe, speed, size);
}
void bridge_free(void* f, unsigned char* wav) {
    ((FreeWaveFunc)f)(wav);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type AquesTalk struct {
	handle   unsafe.Pointer
	fnSynthe unsafe.Pointer
	fnFree   unsafe.Pointer
}

func New(libBasePath string, voice string) (*AquesTalk, error) {
	// 埋め込みSOのパスを構築
	libPath := fmt.Sprintf(libBasePath, voice)
	cPath := C.CString(libPath)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.dlopen(cPath, C.RTLD_NOW)
	if handle == nil {
		return nil, fmt.Errorf("%s voice load error: %s", voice, C.GoString(C.dlerror()))
	}

	return &AquesTalk{
		handle:   handle,
		fnSynthe: C.dlsym(handle, C.CString("AquesTalk_Synthe_Utf8")),
		fnFree:   C.dlsym(handle, C.CString("AquesTalk_FreeWave")),
	}, nil
}

func (a *AquesTalk) Close() error {
	if a.handle != nil {
		C.dlclose(a.handle)
		a.handle = nil
	}
	return nil
}

// 音声合成を実行
func (a *AquesTalk) Synthe(koe string, speed int) ([]byte, error) {
	pKoe := C.CString(koe)
	defer C.free(unsafe.Pointer(pKoe))

	var size C.int
	ptr := C.bridge_synthe(
		a.fnSynthe,
		pKoe,
		C.int(speed),
		&size,
	)

	if ptr == nil {
		return nil, fmt.Errorf("synthesis failed (code: %d)", size)
	}

	// WAVデータのコピーと解放
	defer C.bridge_free(a.fnFree, ptr)

	return C.GoBytes(unsafe.Pointer(ptr), size), nil
}
