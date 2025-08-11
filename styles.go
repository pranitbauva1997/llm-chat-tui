package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 2).
			Align(lipgloss.Center)

	// Input styles
	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	// Spinner styles
	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	// Message bubble styles
	userMessageStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginRight(2).
				Align(lipgloss.Right)

	assistantMessageStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginLeft(2).
				Align(lipgloss.Left)
)

// centerHorizontally creates a style that centers content horizontally within the given width
func centerHorizontally(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center)
}

// centerHorizontallyWithPadding creates a style that centers content with left padding
func centerHorizontallyWithPadding(width, leftPadding int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		PaddingLeft(leftPadding)
}

// rightAlignInContainer creates a style that right-aligns content within the given width
func rightAlignInContainer(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Right)
}

// leftAlignInContainer creates a style that left-aligns content within the given width
func leftAlignInContainer(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left)
}

// createDivider creates a horizontal divider line with the specified width
func createDivider(width int) string {
	return lipgloss.NewStyle().Render(lipgloss.NewStyle().Width(width).Render("─"))
}

// createUserMessageBubble creates a styled user message bubble with dynamic width
func createUserMessageBubble(text string) lipgloss.Style {
	return userMessageStyle.Copy().Width(len(text) + 2)
}

// createAssistantMessageBubble creates a styled assistant message bubble with dynamic width
func createAssistantMessageBubble(text string) lipgloss.Style {
	return assistantMessageStyle.Copy().Width(len(text) + 2)
}

// withMarginTop adds top margin to any style
func withMarginTop(margin int) lipgloss.Style {
	return lipgloss.NewStyle().MarginTop(margin)
}

// withMarginBottom adds bottom margin to any style
func withMarginBottom(margin int) lipgloss.Style {
	return lipgloss.NewStyle().MarginBottom(margin)
}

// alignSpinnerWithInput creates a style for spinner positioning relative to input
func alignSpinnerWithInput(terminalWidth, inputWidth int) lipgloss.Style {
	leftPadding := (terminalWidth - inputWidth) / 2
	return lipgloss.NewStyle().
		Width(terminalWidth).
		PaddingLeft(leftPadding).
		Align(lipgloss.Left)
}
