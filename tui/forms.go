package tui

import (
	"fmt"
	"strings"

	"blackoutbargain/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// FormSubmittedMsg is sent when a form is submitted
type FormSubmittedMsg struct {
	Value string // The submitted value
}

// CodeInputForm creates a huh form for entering codes
func CodeInputForm(prompt string, width int) (formModel, tea.Cmd) {
	var value string

	// Create the form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(prompt).
				Value(&value).
				Placeholder("Enter code..."),
		),
	)

	// Set form styles
	form = form.WithShowHelp(false).WithWidth(width / 2)

	// Create and return the model
	model := formModel{
		form:  form,
		value: &value,
	}

	return model, model.Init()
}

// formModel is a wrapper model for the huh form
type formModel struct {
	form  *huh.Form
	value *string
	err   error
}

func (m formModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle quit for Esc key
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		return m, tea.Quit
	}

	// Update the form
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Check if form is submitted (when form submission is done)
	if form, ok := form.(*huh.Form); ok && form.State == huh.StateCompleted {
		cmds = append(cmds, func() tea.Msg {
			return FormSubmittedMsg{Value: *m.value}
		})
	}

	return m, tea.Batch(cmds...)
}

func (m formModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	return m.form.View()
}

// UpdateModelWithHuhForm integrates huh form into the main model
func UpdateModelWithHuhForm(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FormSubmittedMsg:
		// When form is submitted, process the input code
		input := strings.TrimSpace(msg.Value)

		// Process the code according to what is required
		m.GameState.CurrentInput = "" // Reset current input
		m.GameState.HandleCommand(input)

		// Return to normal mode
		return *m, nil
	}

	// Handle any other form-related messages
	return *m, nil
}

// CreateCodeInputForm creates the appropriate huh form based on game state
func CreateCodeInputForm(gs *game.GameState, width int) tea.Cmd {
	var prompt string

	switch gs.InputRequired {
	case "locker_code":
		prompt = "Enter the locker code:"
	case "safe_code":
		prompt = "Enter the safe code:"
	case "breaker_code":
		prompt = "Enter the breaker activation code:"
	default:
		prompt = "Enter code:"
	}

	return func() tea.Msg {
		// Create a new form model
		model, _ := CodeInputForm(prompt, width)
		return model // Return the model directly as a message
	}
}
