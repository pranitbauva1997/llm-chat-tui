package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/openai/openai-go"
    openai_option "github.com/openai/openai-go/option"
)

type Secrets struct {
    OpenaiApiKey string
}

type ChatMessages struct {
    Role string
    Content string
}

type Chat struct {
    Messages []ChatMessages
    Client   openai.Client
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
    // API key loaded successfully (not printing for security)

    // Initialize OpenAI client
    client := openai.NewClient(openai_option.WithAPIKey(secrets.OpenaiApiKey))

    chat := Chat{
        Client: client,
    }
    
    // Start the bubble tea program
    p := tea.NewProgram(chat)
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(1)
    }
}

func (c Chat) Init() tea.Cmd {
    return nil
}

func (c Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return c, tea.Quit
        }
    }
    return c, nil
}

func (c Chat) View() string {
    return "LLM Chat TUI\n\nPress Ctrl+C to quit"
}
