package game

import (
	"strings"
	"testing"
)

func TestGetLocationName(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		expected string
	}{
		{
			name:     "register",
			location: LocRegister,
			expected: "Near Register 4 (Front)",
		},
		{
			name:     "security station",
			location: LocSecurityStation,
			expected: "Security Station (Electronics)",
		},
		{
			name:     "locker area",
			location: LocLockerArea,
			expected: "Employee Locker Area",
		},
		{
			name:     "manager's office",
			location: LocManagersOffice,
			expected: "Manager's Office",
		},
		{
			name:     "loading dock",
			location: LocLoadingDock,
			expected: "Loading Dock (Back)",
		},
		{
			name:     "escaped",
			location: LocEscaped,
			expected: "Outside (Escaped)",
		},
		{
			name:     "invalid location",
			location: Location(99),
			expected: "Unknown Location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameState{Location: tt.location}
			result := gs.GetLocationName()
			if result != tt.expected {
				t.Errorf("GetLocationName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetLocationDescription(t *testing.T) {
	tests := []struct {
		name     string
		state    *GameState
		expected string
	}{
		{
			name: "register location",
			state: &GameState{
				Location: LocRegister,
			},
			expected: "The Superstore is eerily dark, lit only by emergency signs.",
		},
		{
			name: "security station location",
			state: &GameState{
				Location: LocSecurityStation,
			},
			expected: "You're at the security station in the dimly lit electronics section.",
		},
		{
			name: "locker area - locked",
			state: &GameState{
				Location: LocLockerArea,
				Clues:    map[string]string{},
			},
			expected: "You are standing near the employee lockers. Dale's locker is here. It looks locked.",
		},
		{
			name: "locker area - unlocked",
			state: &GameState{
				Location: LocLockerArea,
				Clues:    map[string]string{"locker_opened": "true"},
			},
			expected: "You are standing near the employee lockers. Dale's locker is here. It's open.",
		},
		{
			name: "loading dock - locked",
			state: &GameState{
				Location: LocLoadingDock,
				Clues:    map[string]string{},
			},
			expected: "You've reached the loading dock area at the back of the store. The storm howls louder here.\nA large breaker panel is on the wall next to the sealed loading door.",
		},
		{
			name: "loading dock - unlocked",
			state: &GameState{
				Location: LocLoadingDock,
				Clues:    map[string]string{"door_unlocked": "true"},
			},
			expected: "You've reached the loading dock area at the back of the store. The storm howls louder here.\nThe heavy loading door stands slightly ajar, unlocked!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.GetLocationDescription()
			// Check that result contains the expected text (partial match)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("GetLocationDescription() = %q, want it to contain %q", result, tt.expected)
			}
		})
	}
}

func TestGetInventoryDescription(t *testing.T) {
	tests := []struct {
		name     string
		state    *GameState
		expected string
	}{
		{
			name: "empty inventory",
			state: &GameState{
				Inventory: map[Item]bool{},
			},
			expected: "Inventory: Empty.",
		},
		{
			name: "one item",
			state: &GameState{
				Inventory: map[Item]bool{
					ItemVoucher: true,
				},
			},
			expected: "Inventory: crumpled employee discount voucher.",
		},
		{
			name: "multiple items",
			state: &GameState{
				Inventory: map[Item]bool{
					ItemVoucher:   true,
					ItemScanner:   true,
					ItemNotebook:  true,
					ItemInventory: true,
				},
			},
			expected: "Inventory:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.GetInventoryDescription()

			if tt.name == "multiple items" {
				// For multiple items, just verify that all items are in the result
				items := []string{
					string(ItemVoucher),
					string(ItemScanner),
					string(ItemNotebook),
					string(ItemInventory),
				}

				for _, item := range items {
					if !strings.Contains(result, item) {
						t.Errorf("GetInventoryDescription() = %q, does not contain %q", result, item)
					}
				}
			} else {
				// For specific expected patterns, check exact match
				if result != tt.expected {
					t.Errorf("GetInventoryDescription() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

func TestGetVisibleItems(t *testing.T) {
	tests := []struct {
		name     string
		state    *GameState
		expected string
	}{
		{
			name: "no visible items",
			state: &GameState{
				Location:  LocRegister,
				Inventory: map[Item]bool{},
			},
			expected: "",
		},
		{
			name: "security station with no items taken",
			state: &GameState{
				Location:  LocSecurityStation,
				Inventory: map[Item]bool{},
			},
			expected: "You see: crumpled employee discount voucher, Dale's handheld scanner.",
		},
		{
			name: "security station with voucher taken",
			state: &GameState{
				Location: LocSecurityStation,
				Inventory: map[Item]bool{
					ItemVoucher: true,
				},
			},
			expected: "You see: Dale's handheld scanner.",
		},
		{
			name: "locker area with open locker",
			state: &GameState{
				Location:  LocLockerArea,
				Inventory: map[Item]bool{},
				Clues:     map[string]string{"locker_opened": "true"},
			},
			expected: "You see: small notebook.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.GetVisibleItems()
			if result != tt.expected {
				t.Errorf("GetVisibleItems() = %q, want %q", result, tt.expected)
			}
		})
	}
}
