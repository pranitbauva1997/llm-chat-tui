package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/openai/openai-go"
	openai_option "github.com/openai/openai-go/option"
	"os"
)

type Secrets struct {
	OpenaiApiKey string
}

type ChatMessages struct {
	Role    string
	Content string
}

type Chat struct {
	Messages  []ChatMessages
	Client    openai.Client
	TextInput textinput.Model
	Width     int
	Height    int
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

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Type your message here..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 50

	chat := Chat{
		Client:    client,
		TextInput: ti,
	}

	// Start the bubble tea program with full screen mode
	p := tea.NewProgram(chat, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func (c Chat) Init() tea.Cmd {
	return nil
}

func (c Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.Width = msg.Width
		c.Height = msg.Height
		// Update text input width to be responsive
		c.TextInput.Width = min(c.Width-4, 80) // Leave some margin, max 80 chars
		return c, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		case "enter":
			// Handle message submission
			if c.TextInput.Value() != "" {
				// Add message to chat history
				c.Messages = append(c.Messages, ChatMessages{
					Role:    "user",
					Content: c.TextInput.Value(),
				})
				// Clear the input
				c.TextInput.SetValue("")
			}
			return c, nil
		}
	}

	// Update the text input
	c.TextInput, cmd = c.TextInput.Update(msg)
	return c, cmd
}

func (c Chat) View() string {
	// Use full terminal dimensions if available, fallback to defaults
	width := c.Width
	height := c.Height
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}

	// Define styles for full screen
	containerStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	// Create the main content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"LLM Chat TUI",
		"",
		c.TextInput.View(),
		"",
		"Press Enter to send, Ctrl+C to quit",
	)

	// Display chat messages if any
	if len(c.Messages) > 0 {
		// Apply width restriction to the entire message history container
		var messageHistoryWidth int
		if width > 120 {
			// Use 80% of width for wide terminals
			messageHistoryWidth = int(float64(width) * 0.8)
		} else {
			// Use full width for narrow terminals
			messageHistoryWidth = width
		}

		var messageLines []string

		// Create message bubbles with proper alignment
		for _, msg := range c.Messages {
			messageText := msg.Content

			// Style messages based on role
			if msg.Role == "user" {
				// User messages on the right (blue background)
				userStyle := lipgloss.NewStyle().
					Background(lipgloss.Color("4")).
					Foreground(lipgloss.Color("15")).
					Padding(0, 1).
					MarginRight(2).
					Width(len(messageText) + 2).
					Align(lipgloss.Right)

				userMessage := userStyle.Render(messageText)
				// Right align the entire message within the constrained width
				rightAligned := lipgloss.NewStyle().
					Width(messageHistoryWidth).
					Align(lipgloss.Right).
					Render(userMessage)
				messageLines = append(messageLines, rightAligned)
			} else {
				// Assistant messages on the left (gray background)
				assistantStyle := lipgloss.NewStyle().
					Background(lipgloss.Color("8")).
					Foreground(lipgloss.Color("15")).
					Padding(0, 1).
					MarginLeft(2).
					Width(len(messageText) + 2).
					Align(lipgloss.Left)

				assistantMessage := assistantStyle.Render(messageText)
				// Left align the entire message within the constrained width
				leftAligned := lipgloss.NewStyle().
					Width(messageHistoryWidth).
					Align(lipgloss.Left).
					Render(assistantMessage)
				messageLines = append(messageLines, leftAligned)
			}

			// Add spacing between messages
			messageLines = append(messageLines, "")
		}

		messageHistoryContainer := lipgloss.NewStyle().
			Width(messageHistoryWidth).
			Align(lipgloss.Center)

		messageHistory := messageHistoryContainer.Render(
			lipgloss.JoinVertical(lipgloss.Left, messageLines...),
		)

		// Layout with messages at top, input at bottom
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			messageHistory,
			"",
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("LLM Chat TUI"),
			"",
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(c.TextInput.View()),
			"",
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("Press Enter to send, Ctrl+C to quit"),
		)
	}

	return containerStyle.Render(content)
}
