format:
	gofmt -w .

build:
	go build .

run:
	./llm-chat-tui

vet:
	go vet

build-and-run: build run
