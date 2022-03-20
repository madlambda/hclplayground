playground:
	GOOS=js GOARCH=wasm go build -o www/playground.wasm

shareh:
	go install github.com/katcipis/shareh@latest

run: playground shareh
	shareh -dir ./www
