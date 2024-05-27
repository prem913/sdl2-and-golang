buildpong:
	@go build -o bin/pong ./cmd/pong
buildgol:
	@go build -o bin/gameoflife ./cmd/gameoflife
runpong: buildpong
	@./bin/pong
rungol: buildgol
	@go build -o bin/gameoflife ./gameoflife/main.go

test:
	@go test ./... -v


