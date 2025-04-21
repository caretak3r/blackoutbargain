package tui

import (
	"fmt"
	"strings"

	"blackoutbargain/game"
	"blackoutbargain/llm"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the BubbleTea TUI model
type Model struct {
	// Terminal state
	Width  int
	Height int

	// Game state
	GameState *game.GameState

	// UI state
	Styles      Styles
	LoadingLLM  bool      // Flag to show loading state
	ActiveForm  tea.Model // Currently active form, if any
	ShowingForm bool      // Flag to indicate if a form is active

	// LLM state
	LLMClient    *llm.Client
	LastLLMInput string // Store the input that triggered the LLM call
}

// New creates a new TUI model
func New(llmClient *llm.Client) Model {
	return Model{
		GameState:  game.NewGameState(),
		Styles:     NewStyles(),
		LLMClient:  llmClient,
		LoadingLLM: false,
	}
}

// --- Bubble Tea Messages ---

// LLMResponseMsg is used for receiving LLM responses
type LLMResponseMsg struct {
	Response string
	Err      error
}

// LLMErrorMsg is used for handling specific LLM API call errors
type LLMErrorMsg struct {
	Err error
}

// --- Bubble Tea Interface Implementation ---

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// No initial command needed for now
	return nil
}

// Update handles input and events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd // Collect multiple commands

	// Handle form-related updates if a form is active
	if m.ShowingForm && m.ActiveForm != nil {
		switch msg := msg.(type) {
		case formModel:
			// Initialize a new form
			m.ActiveForm = msg
			m.ShowingForm = true
			return m, nil

		case tea.KeyMsg:
			if msg.String() == "esc" {
				// Cancel form on escape
				m.ActiveForm = nil
				m.ShowingForm = false
				return m, nil
			}

			// Pass the key message to the active form
			newForm, cmd := m.ActiveForm.Update(msg)
			m.ActiveForm = newForm
			return m, cmd

		case FormSubmittedMsg:
			// Handle form submission
			input := strings.TrimSpace(msg.Value)
			m.GameState.HandleCommand(input)
			m.ActiveForm = nil
			m.ShowingForm = false
			return m, nil
		}

		// For any other message type, try updating the form
		if m.ActiveForm != nil {
			newForm, cmd := m.ActiveForm.Update(msg)
			m.ActiveForm = newForm
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	// --- Handle LLM Response ---
	case LLMResponseMsg:
		m.LoadingLLM = false // Turn off loading indicator
		m.LastLLMInput = ""  // Clear context
		if msg.Err != nil {
			// This case might be less common if errors are caught by LLMErrorMsg
			m.GameState.Message = fmt.Sprintf("LLM Response Error: %s", msg.Err)
		} else {
			m.GameState.Message = msg.Response
			// Optional: Parse LLM response for specific clues recognized by Go game state
			// e.g., if strings.Contains(strings.ToLower(msg.response), "code 4711") { m.GameState.Clues["safe_code_hint"] = "4711" }
		}
		return m, nil // No further command needed now

	case LLMErrorMsg:
		m.LoadingLLM = false
		m.LastLLMInput = ""
		m.GameState.Message = fmt.Sprintf("LLM API Error: %s", msg.Err) // Display the specific error
		return m, nil

	case tea.KeyMsg:
		// --- Handle Input While LLM Loading ---
		if m.LoadingLLM {
			// Allow quitting even while loading
			if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
				return m, tea.Quit
			}
			// Option: Display a message like "Please wait..."
			// Or simply ignore other keys
			return m, nil
		}

		// --- Handle Regular Key Input ---
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(m.GameState.CurrentInput)
			m.GameState.CurrentInput = "" // Reset input field
			m.GameState.Message = ""      // Clear previous message (LLM or Go)

			if input == "" {
				m.GameState.Message = "Please enter a command."
				return m, nil
			}

			// --- Input Routing: Go Logic vs. LLM ---
			if m.GameState.InputRequired != "" {
				// Show form for code input
				m.ShowingForm = true
				return m, CreateCodeInputForm(m.GameState, m.Width)
			} else {
				// --- Delegate general actions based on verb ---
				parts := strings.Fields(strings.ToLower(input))
				verb := ""
				if len(parts) > 0 {
					verb = parts[0]
				}

				switch verb {
				case "go", "g", "take", "t", "inventory", "i", "inv", "help", "h", "escape":
					// Handle these navigation/core actions directly with Go logic
					m.GameState.HandleCommand(input)
					return m, nil

				case "use", "u":
					// Use Go for critical puzzle items/codes, delegate others to LLM
					if m.GameState.IsCriticalUse(input) {
						m.GameState.HandleCommand(input) // Use Go logic
						return m, nil
					} else {
						// Delegate non-critical 'use' to LLM
						if m.LLMClient == nil || !m.LLMClient.Enabled {
							m.GameState.Message = "LLM is disabled. Cannot process this 'use' command flexibly."
							return m, nil
						} else {
							m.LoadingLLM = true
							m.LastLLMInput = input              // Store for loading message
							m.GameState.Message = "Thinking..." // Placeholder message
							return m, m.callLLM(input)
						}
					}

				case "examine", "x", "look", "l", "talk", "ask", "search": // Common verbs for LLM
					// Delegate descriptive/interactive actions to LLM
					if m.LLMClient == nil || !m.LLMClient.Enabled {
						// Fallback Go logic if LLM disabled
						m.GameState.HandleCommand(input)
						return m, nil
					} else {
						m.LoadingLLM = true
						m.LastLLMInput = input
						m.GameState.Message = "Thinking..."
						return m, m.callLLM(input)
					}
				default:
					// Handle unknown verbs - delegate to LLM if available
					if m.LLMClient == nil || !m.LLMClient.Enabled {
						m.GameState.Message = fmt.Sprintf("I don't understand '%s'. Try 'help'.", input)
						return m, nil
					} else {
						m.LoadingLLM = true
						m.LastLLMInput = input
						m.GameState.Message = "Thinking..."
						return m, m.callLLM(input)
					}
				}
			}

		case tea.KeyBackspace:
			if len(m.GameState.CurrentInput) > 0 {
				// Handle UTF-8 runes correctly if needed, but simple slice is ok for basic input
				m.GameState.CurrentInput = m.GameState.CurrentInput[:len(m.GameState.CurrentInput)-1]
			}
			return m, nil

		// Append typed characters to the input buffer
		default:
			// Check if the key press represents a printable character
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
				m.GameState.CurrentInput += string(msg.Runes)
				return m, nil
			}
		}
	}

	// Return the updated model and any commands generated by non-key messages
	return m, tea.Batch(cmds...)
}

// View renders the TUI
func (m Model) View() string {
	if m.GameState.GameOver {
		// Use Lip Gloss for final message too
		finalMsg := "You shove the heavy door open and slip out into the fierce storm. Sirens approach...\n\nYou escaped the Blackout Nightmare!"
		help := "\n\nPress Ctrl+C or Esc to exit."
		return m.Styles.Base.Render(m.Styles.Message.Render(finalMsg)+m.Styles.Help.Render(help)) + "\n"
	}

	if m.Width == 0 {
		return "Initializing terminal size..." // Avoid rendering before we have dimensions
	}

	// If a form is active, render it within our UI
	if m.ShowingForm && m.ActiveForm != nil {
		var s strings.Builder

		// Title
		title := m.Styles.Title.Render("--- Blackout Bargain ---")
		s.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, title))
		s.WriteString("\n\n")

		// Game context
		s.WriteString(m.Styles.Location.Render(m.GameState.GetLocationDescription()))
		s.WriteString("\n\n")

		// Form
		formView := m.ActiveForm.View()
		styledForm := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69")).
			Render(formView)

		s.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, styledForm))
		s.WriteString("\n\n")

		// Footer
		s.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Left,
			m.Styles.Help.Render("Enter to submit, Esc to cancel")))

		return s.String()
	}

	var s strings.Builder

	// --- Title ---
	title := m.Styles.Title.Render("--- Blackout Bargain ---")
	// Center title within available width
	s.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, title))
	s.WriteString("\n\n")

	// --- Main Content Area ---
	var mainContent strings.Builder

	if m.LoadingLLM {
		// --- Loading State ---
		loadingText := fmt.Sprintf("Processing '%s'...", m.LastLLMInput)
		// You could add a spinner here using charm/bubbles/spinner
		mainContent.WriteString(m.Styles.Message.Foreground(lipgloss.Color("220")).Render(loadingText)) // Yellowish message
	} else {
		// --- Normal Game State View ---
		// Location Description (Generated by Go)
		mainContent.WriteString(m.Styles.Location.Render(m.GameState.GetLocationDescription()))
		mainContent.WriteString("\n") // Add spacing

		// Visible Items (Generated by Go)
		visibleItems := m.GameState.GetVisibleItems()
		if visibleItems != "" {
			mainContent.WriteString(m.Styles.Items.Render(visibleItems))
			mainContent.WriteString("\n")
		}

		// Inventory (Generated by Go)
		mainContent.WriteString(m.Styles.Inventory.Render(m.GameState.GetInventoryDescription()))
		mainContent.WriteString("\n\n") // More spacing

		// Message Area (Go messages or LLM response)
		if m.GameState.Message != "" {
			mainContent.WriteString(m.Styles.Message.Render(m.GameState.Message))
			mainContent.WriteString("\n\n")
		}

		// Input Prompt
		mainContent.WriteString(m.getStyledInputPrompt())
	}

	// Combine content and apply base padding/styling
	// Ensure content doesn't exceed terminal height (basic wrapping)
	styledContent := m.Styles.Base.Width(m.Width).Render(mainContent.String())
	s.WriteString(styledContent)

	// --- Footer Help Text ---
	s.WriteString("\n\n")
	s.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Left, m.Styles.Help.Render("Ctrl+C or Esc to quit.")))

	return s.String()
}

// getStyledInputPrompt returns the styled input prompt string
func (m Model) getStyledInputPrompt() string {
	return m.Styles.Prompt.Render(m.GameState.GetInputPrompt()) + m.GameState.CurrentInput
}

// callLLM constructs a command to call the LLM service
func (m Model) callLLM(playerInput string) tea.Cmd {
	return func() tea.Msg {
		if m.LLMClient == nil || !m.LLMClient.Enabled {
			return LLMResponseMsg{
				Response: "LLM support is not available. Using basic descriptions.",
				Err:      nil,
			}
		}

		response, err := m.LLMClient.GenerateResponse(playerInput, m.GameState)
		if err != nil {
			return LLMErrorMsg{Err: err}
		}

		return LLMResponseMsg{Response: response, Err: nil}
	}
}
