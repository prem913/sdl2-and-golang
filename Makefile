buildpong:
	@go build -o bin/pong ./cmd/pong
buildgol:
	@go build -o bin/gameoflife ./cmd/gameoflife
buildfp:
	@go build -o bin/flappy_birb ./cmd/flappy_birb
runpong: buildpong
	@./bin/pong
rungol: 
	@go run ./cmd/gameoflife
runspi:
	@go run ./cmd/spaceinvaders/
runfp:
	@go run ./cmd/flappy_birb/
runwindow:
	@go run ./cmd/pongwindow/


test:
	@go test ./... -v


