set -e
GOOS=js GOARCH=wasm go build -o app.wasm ./wasm
go run  .
