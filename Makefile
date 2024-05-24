buildp:
	@go build -o bin/pong ./cmd/pong

buildw:
	@GOOS=windows GOARCH=386 go build -o bin/pong-windows ./cmd/pong

run: build
	@./bin/pong

test:
	@go test ./... -v


