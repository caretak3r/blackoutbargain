package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the lipgloss styles used by the TUI
type Styles struct {
	Base      lipgloss.Style
	Title     lipgloss.Style
	Location  lipgloss.Style
	Items     lipgloss.Style
	Inventory lipgloss.Style
	Message   lipgloss.Style
	Help      lipgloss.Style
	Prompt    lipgloss.Style
}

// NewStyles creates a new set of styles with default values
func NewStyles() Styles {
	s := Styles{}
	// Define styles using Lip Gloss
	s.Base = lipgloss.NewStyle().Padding(0, 1)                                                // Basic padding
	s.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).MarginBottom(1) // Purple, spacing
	s.Location = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))              // Light Blue/Cyan
	s.Items = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))                            // Cyan
	s.Inventory = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))                       // Orange
	s.Message = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Italic(true)            // Light Gray Italic
	s.Help = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))                            // Dark Gray
	s.Prompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))               // White prompt
	return s
}
