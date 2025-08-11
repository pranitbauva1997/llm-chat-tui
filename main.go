package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
)

type Secrets struct {
    OpenaiApiKey string
}

type ChatMessages struct {
    Role string
    Content string
}

type Chat struct {
    messages []ChatMessages
}

func main() {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("OPENAI_API_KEY environment variable not set")
        return
    }
    secrets := Secrets{
        OpenaiApiKey: apiKey,
    }
    fmt.Printf("OpenAI API Key: %s\n", secrets.OpenaiApiKey)

    chat := Chat{}
}

func (c Chat) Init() tea.Cmd {
    return nil
}
