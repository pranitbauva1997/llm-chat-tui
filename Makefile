format:
	gofmt -w .

build:
	go build .

run:
	./llm-chat-tui

build-and-run: build run
