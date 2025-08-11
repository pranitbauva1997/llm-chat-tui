format:
	gofmt -w .

build:
	go build .

run:
	./llm-chat-tui

vet:
	go vet

clean:
	go clean -cache
	rm llm-chat-tui

build-and-run: build run
