# ビルドステージ - Golang環境を使用してビルド
FROM golang:1.24 AS builder

# 作業ディレクトリの設定
WORKDIR /app

# キャッシュを効率的に使うためにgo.modとgo.sumを先にコピー
COPY go.* ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Makefileを使ってビルド
RUN make build

# 実行ステージ - 軽量なDebianベースのイメージ
FROM gcr.io/distroless/base-debian10

# ビルドされたバイナリをコピー
COPY --from=builder /app/aquestalk-server /usr/local/bin/main

# アプリケーションを実行
CMD ["main"]