package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/openai/openai-go"
	openai_option "github.com/openai/openai-go/option"
	"os"
	"strings"
)

type Secrets struct {
	OpenaiApiKey string
}

type ChatMessages struct {
	Role    string
	Content string
}

// Custom message types for bubble tea
type streamChunkMsg struct {
	content string
}

type streamCompleteMsg struct{}

type streamErrorMsg struct {
	err error
}

type Chat struct {
	Messages       []ChatMessages
	Client         openai.Client
	TextInput      textinput.Model
	Viewport       viewport.Model
	Width          int
	Height         int
	isStreaming    bool
	currentMessage strings.Builder
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

	// Initialize viewport with welcome content
	vp := viewport.New(80, 20)

	chat := Chat{
		Client:         client,
		TextInput:      ti,
		Viewport:       vp,
		isStreaming:    false,
		currentMessage: strings.Builder{},
	}

	// Start the bubble tea program with full screen mode and mouse support
	p := tea.NewProgram(chat, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// sendMessage sends the chat messages to OpenAI and returns a command for streaming
func (c Chat) sendMessage() tea.Cmd {
	return tea.Batch(c.startStreaming())
}

// startStreaming creates multiple commands to handle real-time streaming
func (c Chat) startStreaming() tea.Cmd {
	return func() tea.Msg {
		// Convert ChatMessages to OpenAI format
		var messages []openai.ChatCompletionMessageParamUnion
		for _, msg := range c.Messages {
			if msg.Role == "user" {
				messages = append(messages, openai.UserMessage(msg.Content))
			} else if msg.Role == "assistant" {
				messages = append(messages, openai.AssistantMessage(msg.Content))
			}
		}

		// Create streaming chat completion request
		stream := c.Client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    openai.ChatModelGPT4oMini,
		})

		// Process stream and collect all content for now
		// TODO: Implement true real-time streaming with channels if needed
		var fullContent strings.Builder
		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				fullContent.WriteString(chunk.Choices[0].Delta.Content)
			}
		}

		if err := stream.Err(); err != nil {
			return streamErrorMsg{err: err}
		}

		// Return both the content and completion message
		if fullContent.Len() > 0 {
			// We need to send both the content and completion
			// But since we can only return one message, let's send the content
			// and let the handler send the completion
			return streamChunkMsg{content: fullContent.String()}
		}

		return streamCompleteMsg{}
	}
}

// renderMessages renders all messages and returns the content string
func (c Chat) renderMessages() string {
	var messageLines []string

	// Calculate available width for messages and center container
	terminalWidth := c.Viewport.Width
	if terminalWidth == 0 {
		terminalWidth = 80 // fallback width
	}

	// Set max message container width and calculate centering
	maxMessageWidth := 120
	messageHistoryWidth := terminalWidth
	if messageHistoryWidth > maxMessageWidth {
		messageHistoryWidth = maxMessageWidth
	}

	// Calculate padding for centering the message container
	leftPadding := 0
	if terminalWidth > messageHistoryWidth {
		leftPadding = (terminalWidth - messageHistoryWidth) / 2
	}

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

			// Center the entire message container
			centeredMessage := lipgloss.NewStyle().
				Width(terminalWidth).
				PaddingLeft(leftPadding).
				Render(rightAligned)
			messageLines = append(messageLines, centeredMessage)
		} else if msg.Role == "assistant" {
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

			// Center the entire message container
			centeredMessage := lipgloss.NewStyle().
				Width(terminalWidth).
				PaddingLeft(leftPadding).
				Render(leftAligned)
			messageLines = append(messageLines, centeredMessage)
		}

		// Add spacing between messages
		messageLines = append(messageLines, "")
	}

	// Show current streaming message if streaming
	if c.isStreaming {
		streamingText := c.currentMessage.String()
		if streamingText == "" {
			streamingText = "..."
		} else {
			streamingText += "▋" // Add cursor to show it's streaming
		}

		// Assistant messages on the left (gray background) with streaming indicator
		assistantStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("8")).
			Foreground(lipgloss.Color("15")).
			Padding(0, 1).
			MarginLeft(2).
			Width(len(streamingText) + 2).
			Align(lipgloss.Left)

		assistantMessage := assistantStyle.Render(streamingText)
		// Left align the entire message within the constrained width
		leftAligned := lipgloss.NewStyle().
			Width(messageHistoryWidth).
			Align(lipgloss.Left).
			Render(assistantMessage)

		// Center the entire message container
		centeredMessage := lipgloss.NewStyle().
			Width(terminalWidth).
			PaddingLeft(leftPadding).
			Render(leftAligned)
		messageLines = append(messageLines, centeredMessage)
		messageLines = append(messageLines, "")
	}

	// Join all message lines and return content
	return strings.Join(messageLines, "\n")
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
		// Update viewport dimensions - reserve space for input and status
		viewportHeight := c.Height - 6 // Reserve space for title, input, and status
		c.Viewport.Width = c.Width
		c.Viewport.Height = viewportHeight
		return c, nil
	case tea.MouseMsg:
		// Handle mouse wheel events for viewport scrolling using new API
		if msg.Action == tea.MouseActionPress && (msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown) {
			c.Viewport, cmd = c.Viewport.Update(msg)
			return c, cmd
		}
		// For other mouse events, let them pass through
		return c, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		case "enter":
			// Handle message submission
			inputValue := c.TextInput.Value()
			if inputValue != "" && !c.isStreaming {
				// Add user message to chat history
				c.Messages = append(c.Messages, ChatMessages{
					Role:    "user",
					Content: inputValue,
				})
				// Update viewport content immediately
				content := c.renderMessages()
				c.Viewport.SetContent(content)
				c.Viewport.GotoBottom()
				// Clear the input
				c.TextInput.SetValue("")
				// Set streaming state and send to OpenAI
				c.isStreaming = true
				c.currentMessage.Reset()
				// Update viewport to show streaming indicator immediately
				content = c.renderMessages()
				c.Viewport.SetContent(content)
				c.Viewport.GotoBottom()
				return c, c.sendMessage()
			}
			return c, nil
		case "up", "down", "pgup", "pgdown", "home", "end":
			// Handle viewport scrolling
			c.Viewport, cmd = c.Viewport.Update(msg)
			return c, cmd
		}

		// Let the text input handle other keys when not streaming
		if !c.isStreaming {
			c.TextInput, cmd = c.TextInput.Update(msg)
			return c, cmd
		}
	case streamChunkMsg:
		// Handle streaming response chunk
		c.currentMessage.WriteString(msg.content)
		// Update viewport content immediately for streaming
		content := c.renderMessages()
		c.Viewport.SetContent(content)
		c.Viewport.GotoBottom()
		// After receiving content, send completion message
		return c, func() tea.Msg { return streamCompleteMsg{} }
	case streamCompleteMsg:
		// Handle streaming completion
		if c.currentMessage.Len() > 0 {
			// Add the complete assistant message
			c.Messages = append(c.Messages, ChatMessages{
				Role:    "assistant",
				Content: c.currentMessage.String(),
			})
			c.currentMessage.Reset()
		}
		c.isStreaming = false
		// Update viewport content immediately to clear streaming indicator
		content := c.renderMessages()
		c.Viewport.SetContent(content)
		c.Viewport.GotoBottom()
		// Re-focus the text input after streaming is complete
		return c, c.TextInput.Focus()
	case streamErrorMsg:
		// Handle streaming error
		c.Messages = append(c.Messages, ChatMessages{
			Role:    "assistant",
			Content: fmt.Sprintf("Error: %v", msg.err),
		})
		c.isStreaming = false
		// Re-focus the text input after error
		return c, c.TextInput.Focus()
	}

	// For non-key messages (like window resize), always update the text input
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

	// Create status message based on streaming state
	statusMsg := "Press Enter to send, Ctrl+C to quit • Use ↑/↓ or Page Up/Down to scroll with mouse wheel"
	if c.isStreaming {
		statusMsg = "Streaming response... Please wait"
	}

	// Check if there are any messages (excluding streaming state)
	hasMessages := len(c.Messages) > 0

	var layout string
	if !hasMessages {
		// Center the input when no messages exist
		// Calculate vertical centering
		availableHeight := height - 6 // Account for title, input, status, and spacing
		topPadding := availableHeight / 2

		var topSpacing []string
		for i := 0; i < topPadding; i++ {
			topSpacing = append(topSpacing, "")
		}

		layout = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("LLM Chat TUI"),
			"",
			strings.Join(topSpacing, "\n"),
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(c.TextInput.View()),
			"",
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(statusMsg),
		)
	} else {
		// Keep input at bottom when messages exist
		layout = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render("LLM Chat TUI"),
			"",
			c.Viewport.View(), // Use viewport for scrollable messages
			"",
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(c.TextInput.View()),
			"",
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(statusMsg),
		)
	}

	return layout
}
