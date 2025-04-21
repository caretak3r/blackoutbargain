package game

import (
	"fmt"
	"strings"
)

// GetLocationName returns a short name for the current location.
func (gs *GameState) GetLocationName() string {
	switch gs.Location {
	case LocRegister:
		return "Near Register 4 (Front)"
	case LocSecurityStation:
		return "Security Station (Electronics)"
	case LocLockerArea:
		return "Employee Locker Area"
	case LocManagersOffice:
		return "Manager's Office"
	case LocLoadingDock:
		return "Loading Dock (Back)"
	case LocEscaped:
		return "Outside (Escaped)"
	default:
		return "Unknown Location"
	}
}

// GetLocationDescription provides the base description for the current location.
func (gs *GameState) GetLocationDescription() string {
	// These are base descriptions; LLM can elaborate when examining the area.
	switch gs.Location {
	case LocRegister:
		return "The Superstore is eerily dark, lit only by emergency signs. Thunder rattles the windows. You're near Register 4 with Brenda and Gary. The main doors are dead silent and locked."
	case LocSecurityStation:
		return "You're at the security station in the dimly lit electronics section. Dale's body is slumped against the dark monitors."
	case LocLockerArea:
		desc := "You are standing near the employee lockers. Dale's locker is here."
		if _, opened := gs.Clues["locker_opened"]; opened {
			desc += " It's open."
		} else {
			desc += " It looks locked."
		}
		return desc
	case LocManagersOffice:
		return "You are inside the cramped manager's office. There's a desk, a corkboard, and a small safe embedded in the wall."
	case LocLoadingDock:
		desc := "You've reached the loading dock area at the back of the store. The storm howls louder here."
		if _, unlocked := gs.Clues["door_unlocked"]; unlocked {
			desc += "\nThe heavy loading door stands slightly ajar, unlocked!"
		} else {
			desc += "\nA large breaker panel is on the wall next to the sealed loading door."
		}
		return desc
	case LocEscaped:
		return "You are outside in the raging storm." // Should be handled by gameOver view
	default:
		return "You are somewhere..."
	}
}

// GetVisibleItems lists items available to 'take' in the current location.
func (gs *GameState) GetVisibleItems() string {
	items := []string{}
	switch gs.Location {
	case LocSecurityStation:
		if !gs.Inventory[ItemVoucher] {
			items = append(items, string(ItemVoucher))
		}
		if !gs.Inventory[ItemScanner] {
			items = append(items, string(ItemScanner))
		}
	case LocLockerArea:
		if _, opened := gs.Clues["locker_opened"]; opened && !gs.Inventory[ItemNotebook] {
			items = append(items, string(ItemNotebook))
		}
	case LocManagersOffice:
		if !gs.Inventory[ItemCard] {
			items = append(items, string(ItemCard))
		}
		if !gs.Inventory[ItemInventory] {
			items = append(items, string(ItemInventory))
		}
		if _, opened := gs.Clues["safe_opened"]; opened && !gs.Inventory[ItemOverrideKey] {
			items = append(items, string(ItemOverrideKey))
		}
	}
	if len(items) > 0 {
		return "You see: " + strings.Join(items, ", ") + "."
	}
	return ""
}

// GetInventoryDescription lists items the player is carrying.
func (gs *GameState) GetInventoryDescription() string {
	if len(gs.Inventory) == 0 {
		return "Inventory: Empty."
	}
	items := []string{}
	for itm := range gs.Inventory {
		items = append(items, string(itm))
	}
	// Sort inventory for consistent display? (Optional)
	// sort.Strings(items)
	return "Inventory: " + strings.Join(items, ", ") + "."
}

// GetInputPrompt returns the input prompt string.
func (gs *GameState) GetInputPrompt() string {
	prompt := "> "
	if gs.InputRequired != "" {
		prompt = fmt.Sprintf("Enter %s: ", strings.ReplaceAll(gs.InputRequired, "_", " "))
	}
	return prompt + gs.CurrentInput // Display current typed input
}
