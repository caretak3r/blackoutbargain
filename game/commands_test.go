package game

import (
	"strings"
	"testing"
)

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		initialState   *GameState
		expectedState  *GameState
		expectedRetval bool
	}{
		{
			name:  "help command",
			input: "help",
			initialState: &GameState{
				Location:  LocRegister,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
			},
			expectedState: &GameState{
				Location:  LocRegister,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
				Message:   "Commands: look (l), go [place] (g), examine [item/area] (x), take [item] (t), use [item/code] (u), inventory (i), help (h), escape. \nUse 'examine' or 'look' for more details (handled by AI if available).",
			},
			expectedRetval: true,
		},
		{
			name:  "go to security station",
			input: "go security",
			initialState: &GameState{
				Location:  LocRegister,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
			},
			expectedState: &GameState{
				Location:  LocSecurityStation,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
				Message:   "You hurry towards the back of the store, near the electronics section where the scream came from.",
			},
			expectedRetval: true,
		},
		{
			name:  "take voucher",
			input: "take crumpled employee discount voucher",
			initialState: &GameState{
				Location:  LocSecurityStation,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
			},
			expectedState: &GameState{
				Location: LocSecurityStation,
				Inventory: map[Item]bool{
					ItemVoucher: true,
				},
				Clues:   make(map[string]string),
				Message: "You take the crumpled employee discount voucher.",
			},
			expectedRetval: true,
		},
		{
			name:  "examine fallback",
			input: "examine room",
			initialState: &GameState{
				Location:  LocSecurityStation,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
			},
			expectedState: &GameState{
				Location:  LocSecurityStation,
				Inventory: make(map[Item]bool),
				Clues:     make(map[string]string),
				Message:   "Dale's body is here. Monitors dark. Scanner nearby? Voucher clutched? Dale clutches a voucher. A scanner lies nearby.",
			},
			expectedRetval: false,
		},
		{
			name:  "correct locker code",
			input: "8675309",
			initialState: &GameState{
				Location:      LocLockerArea,
				Inventory:     make(map[Item]bool),
				Clues:         make(map[string]string),
				InputRequired: "locker_code",
			},
			expectedState: &GameState{
				Location: LocLockerArea,
				Inventory: map[Item]bool{
					ItemNotebook: true,
				},
				Clues: map[string]string{
					"locker_opened": "true",
				},
				InputRequired: "",
				Message:       "Click! The locker swings open.\nYou find a small notebook inside and take it.",
			},
			expectedRetval: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := tt.initialState
			result := gs.HandleCommand(tt.input)

			if result != tt.expectedRetval {
				t.Errorf("HandleCommand() returned %v, want %v", result, tt.expectedRetval)
			}

			// Check location
			if gs.Location != tt.expectedState.Location {
				t.Errorf("Location = %v, want %v", gs.Location, tt.expectedState.Location)
			}

			// Check inventory items
			for item, expected := range tt.expectedState.Inventory {
				actual, exists := gs.Inventory[item]
				if expected && (!exists || !actual) {
					t.Errorf("Inventory missing expected item %v", item)
				}
			}

			// Check clue state
			for clue, expectedVal := range tt.expectedState.Clues {
				actualVal, exists := gs.Clues[clue]
				if !exists || actualVal != expectedVal {
					t.Errorf("Clue %v = %v, want %v", clue, actualVal, expectedVal)
				}
			}

			// Check message (strip whitespace to focus on content)
			actualMsg := strings.TrimSpace(gs.Message)
			expectedMsg := strings.TrimSpace(tt.expectedState.Message)
			if actualMsg != expectedMsg {
				t.Errorf("Message = %q, want %q", actualMsg, expectedMsg)
			}

			// Check input required state
			if gs.InputRequired != tt.expectedState.InputRequired {
				t.Errorf("InputRequired = %v, want %v", gs.InputRequired, tt.expectedState.InputRequired)
			}
		})
	}
}

func TestIsCriticalUse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		state    *GameState
		expected bool
	}{
		{
			name:  "locker code at locker area",
			input: "use 8675309 on locker",
			state: &GameState{
				Location:  LocLockerArea,
				Inventory: make(map[Item]bool),
			},
			expected: true,
		},
		{
			name:  "locker code at wrong location",
			input: "use 8675309 on locker",
			state: &GameState{
				Location:  LocRegister,
				Inventory: make(map[Item]bool),
			},
			expected: false,
		},
		{
			name:  "safe code at manager office",
			input: "use 4711 on safe",
			state: &GameState{
				Location:  LocManagersOffice,
				Inventory: make(map[Item]bool),
			},
			expected: true,
		},
		{
			name:  "key at loading dock with key",
			input: "use key with panel",
			state: &GameState{
				Location: LocLoadingDock,
				Inventory: map[Item]bool{
					ItemOverrideKey: true,
				},
			},
			expected: true,
		},
		{
			name:  "key at loading dock without key",
			input: "use key with panel",
			state: &GameState{
				Location:  LocLoadingDock,
				Inventory: make(map[Item]bool),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.IsCriticalUse(tt.input)
			if result != tt.expected {
				t.Errorf("IsCriticalUse(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
