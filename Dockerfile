FROM debian:11 AS builder

WORKDIR /app

COPY go.* ./

RUN apt-get update && apt-get upgrade -y && apt-get install -y make curl gcc
RUN ARCH=$(dpkg --print-architecture) && \
if [ "$ARCH" = "arm64" ]; then GOARCH="arm64"; else GOARCH="amd64"; fi && \
curl -OL https://go.dev/dl/go1.24.0.linux-${GOARCH}.tar.gz && \
tar -C /usr/local -xzf go1.24.0.linux-${GOARCH}.tar.gz && \
rm -rf go1.24.0.linux-${GOARCH}.tar.gz
ENV PATH $PATH:/usr/local/go/bin
ADD ./ .

RUN go mod download
RUN make build

# 実行ステージ
FROM gcr.io/distroless/base-debian11 AS prod

COPY --from=builder /usr/lib/x86_64-linux-gnu/libstdc++.so.6 /usr/lib/
COPY --from=builder /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/
COPY --from=builder /lib/x86_64-linux-gnu/libc.so.6 /lib/
COPY --from=builder /lib/x86_64-linux-gnu/libm.so.6 /lib/
COPY --from=builder /app/aquestalk-server /usr/local/bin/main

CMD ["main"]