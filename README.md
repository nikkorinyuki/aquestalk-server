# AquesTalk TTS API サーバー (OpenAI Compatible)

このリポジトリは、AquesTalkの動的ライブラリ (.so) を呼び出して、日本語音声合成を提供するHTTP APIサーバーを実装したものです。APIはOpenAIの音声合成APIに互換しており、主に「ゆっくりボイス」の音声を生成します。

## 特徴

- **超軽量な日本語音声合成**: AquesTalkを利用した日本語音声合成。
- **OpenAI互換API**: OpenAIクライアントライブラリで利用可能な音声合成API。

## 使用上の注意

- **プラットフォーム**: Linux 64bit(x86_64)を前提としています。
- **AquesTalk SDKについて**: SDKは付属していません。ご自身で別途用意ください。
- **ライセンス**: 本プロジェクトはGPL 3.0ライセンスの下で提供されます。AquesTalkの使用にはライセンスを読み、ライセンスに従ってください。

# AquesTalk SDKについて

当リポジトリに付属していません。ご自身で用意し、環境変数にて場所を指定してください。

なお`言語処理ライブラリ AqKanji2Koe Linux Ver. 4.1`と`規則音声合成ライブラリ AquesTalk1 Linux Ver.1.7`を想定して作成しています。

環境変数（例）
- `AqKanji2Koe_LibPath`: `./aquestalk/aqk2k_lnx/lib/libAqKanji2Koe.so.4.1`
- `AqKanji2Koe_DicPath`: `./aquestalk/aqk2k_lnx/aq_dic`
- `AqKanji2Koe_DevKey`: `(開発ライセンスキーを入力)`
- `AquesTalk_LibPath`: `./aquestalk/aqtk1-lnx/lib64/%s/libAquesTalk.so`

## APIエンドポイント

- **POST** `/v1/audio/speech`
  
  音声合成を行うエンドポイントです。

### リクエスト例 (curl)

```bash
curl http://localhost:8080/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{
    "model": "tts-1",
    "input": "今日はいい天気ですね。",
    "voice": "f1",
    "response_format": "wav",
    "speed": 1.0
  }' \
  --output speech.wav
```

### パラメータ

- `model` (string): 固定値 `tts-1`。`tts-1-hd` は使用できません。
- `voice` (string): 使用する音声の種類。`dvd`, `f1`, `f2`, `imd1`, `jgr`, `m1`, `m2`, `r1` のいずれか。
- `input` (string): 合成するテキスト。
- `isKana` (bool?): inputで指定されたテキストがカナ音声記号列かを指定します。
- `response_format` (string): 固定値 `wav`。それ以外はエラーとなります。
- `speed` (float): 音声の速度。0.5から3.0の間の値で指定します。

- **POST** `/pronunciation`
  
  音声合成を行うエンドポイントです。

### リクエスト例 (curl)

```bash
curl http://localhost:8080/pronunciation \
  -H "Content-Type: application/json" \
  -d '{
    "input": "今日はいい天気ですね。"
  }' \
  --output output.txt
```

### パラメータ

- `input` (string): 変換するテキスト。

## 内部テキスト変換処理

AquesTalkは独自の音声記号列を入力として受け付ける仕様となっています。そのため、本APIでは、漢字、数字、英語が混在する通常の文章を、AquesTalkが認識可能な音声記号列に変換する処理を内部で行っています。

テキストから音声記号列への変換には、OpenJTalkの辞書を活用した独自実装の外部モジュール [AqKanji2Koe-OpenJTalk](https://github.com/Lqm1/aqkanji2koe-openjtalk) を使用しています。これは株式会社アクエストの商用製品「AqKanji2Koe」とは異なる独自実装です。これにより、ユーザーは通常の文章を入力するだけで、AquesTalkが要求する音声記号列へと自動変換され、スムーズな音声合成が可能となります。

## サンプルコード (Python)

以下のサンプルコードでは、OpenAIクライアントライブラリを使ってAPIを呼び出します。

```python
from openai import OpenAI

client = OpenAI(api_key="a", base_url="http://localhost:8080/v1")
response = client.audio.speech.create(
    model="tts-1",
    voice="f1",
    input="吾輩は猫である。名前はまだ無い。\nどこで生れたかとんと見当がつかぬ。何でも薄暗いじめじめした所でニャーニャー泣いていた事だけは記憶している。吾輩はここで始めて人間というものを見た。しかもあとで聞くとそれは書生という人間中で一番獰悪な種族であったそうだ。この書生というのは時々我々を捕えて煮て食うという話である。しかしその当時は何という考もなかったから別段恐しいとも思わなかった。ただ彼の掌に載せられてスーと持ち上げられた時何だかフワフワした感じがあったばかりである。掌の上で少し落ちついて書生の顔を見たのがいわゆる人間というものの見始であろう。この時妙なものだと思った感じが今でも残っている。第一毛をもって装飾されべきはずの顔がつるつるしてまるで薬缶だ。その後猫にもだいぶ逢ったがこんな片輪には一度も出会わした事がない。のみならず顔の真中があまりに突起している。そうしてその穴の中から時々ぷうぷうと煙を吹く。どうも咽せぽくて実に弱った。これが人間の飲む煙草というものである事はようやくこの頃知った。",
)
```

## 免責事項

- 本リポジトリの開発者は、ライセンス違反に対して一切の責任を負いません。

## ライセンス

本プロジェクトはGPL 3.0ライセンスの下で提供されています。AquesTalkを使用する際には、ライセンスを必ず確認し、ライセンスに従ってください。

## 参考リンク

- [AquesTalk公式ブログ](http://blog-yama.a-quest.com/?eid=970181)
- [AquesTalk FAQ](https://www.a-quest.com/faq.html)
