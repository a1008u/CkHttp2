# Http2


## 準備
## 1.Open sslで証明書の準備

```
openssl req -new -key rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
```

## 2.curlでhttp2を実行できる準備
```
brew install curl-openssl
```

## Goでhelloを返すサーバーを立てる
```go
// 起動だけ
go run main.go

// デバッグでの起動
env GODEBUG=http2debug=2 go run main.go
```



## http2で起動しているかチェック
```curl
curl  -v -XGET https://0.0.0.0:8000 --insecure
```
