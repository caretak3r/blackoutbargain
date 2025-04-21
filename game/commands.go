package game

import (
	"fmt"
	"strings"
)

// HandleCommand processes a player's command and updates the game state.
// Returns true if the command was successfully processed.
func (gs *GameState) HandleCommand(input string) bool {
	input = strings.ToLower(input)

	// Handle specific input prompts first (codes)
	if gs.InputRequired != "" {
		switch gs.InputRequired {
		case "locker_code":
			if input == "8675309" {
				gs.Message = "Click! The locker swings open."
				gs.Clues["locker_opened"] = "true"
				if !gs.Inventory[ItemNotebook] { // Check if already have it somehow
					gs.Inventory[ItemNotebook] = true
					gs.Message += "\nYou find a " + string(ItemNotebook) + " inside and take it."
				} else {
					gs.Message += "\nIt's empty now."
				}
			} else {
				gs.Message = "Incorrect code. The lock doesn't budge."
			}
			gs.InputRequired = "" // Clear requirement
			return true
		case "safe_code":
			if input == "4711" {
				gs.Message = "Click! The safe door opens."
				gs.Clues["safe_opened"] = "true"
				if !gs.Inventory[ItemOverrideKey] { // Check if already have it
					gs.Inventory[ItemOverrideKey] = true // Give the key
					gs.Message += "\nYou find the " + string(ItemOverrideKey) + " inside and take it."
				} else {
					gs.Message += "\nIt's empty now."
				}
			} else {
				gs.Message = "Incorrect code. The safe remains locked."
			}
			gs.InputRequired = "" // Clear requirement
			return true
		case "breaker_code":
			// Keypad mapping: O=6, V=8, E=3, R=7, S=7, T=8, O=6, C=2, K=5 -> 683778625
			correctCode := input == "overstock" || input == "683778625"
			_, hasKey := gs.Inventory[ItemOverrideKey]

			if correctCode && hasKey {
				gs.Message = "CLUNK! A heavy sound echoes - the main magnetic door locks release.\nSuddenly, Gary lunges! 'You meddling kids!' From the shadows, Brenda appears, holding a wrench. 'Dale knew you were skimming, Gary!' she shouts.\nAfter a brief struggle, Gary is subdued near the loading dock door's manual release lever."
				gs.Clues["door_unlocked"] = "true"
				gs.Message += "\nYou can now 'escape' through the loading dock door."
			} else if !hasKey {
				gs.Message = "You need the Manual Override Key inserted to activate the panel."
			} else { // Has key but wrong code
				gs.Message = "Incorrect code entered on the keypad. Nothing happens."
			}
			gs.InputRequired = "" // Clear requirement
			return true
		}
	}

	// General command parsing (simple version)
	parts := strings.Fields(input)
	if len(parts) == 0 {
		gs.Message = "Please enter a command like 'look', 'go security', 'take voucher', 'use 8675309', 'inventory', or 'help'."
		return false
	}
	verb := parts[0]
	object := ""
	if len(parts) > 1 {
		object = strings.Join(parts[1:], " ") // Re-join object words
	}

	// Handle verbs managed directly by Go
	switch verb {
	case "go", "g":
		gs.handleGo(object)
	case "take", "t":
		gs.handleTake(object)
	case "use", "u":
		// This case should only be reached if isCriticalUse was true
		gs.handleCriticalUse(object, input) // Pass full input for code checks
	case "inventory", "i", "inv":
		gs.Message = gs.GetInventoryDescription() // Show inventory directly
	case "help", "h":
		gs.Message = "Commands: look (l), go [place] (g), examine [item/area] (x), take [item] (t), use [item/code] (u), inventory (i), help (h), escape. \nUse 'examine' or 'look' for more details (handled by AI if available)."
	case "escape":
		if gs.Location == LocLoadingDock && gs.Clues["door_unlocked"] == "true" {
			gs.GameOver = true // Trigger game end sequence in View()
		} else if gs.Location == LocLoadingDock {
			gs.Message = "You try the heavy loading dock door, but it's still magnetically locked."
		} else {
			gs.Message = "You can't escape from here. You need to reach the unlocked loading dock door."
		}
	case "look", "l", "examine", "x": // Basic fallback if LLM disabled
		gs.handleExamineFallback(object)
		return false // Indicate this should be handled by LLM if available
	default:
		gs.Message = fmt.Sprintf("I don't understand '%s'. Try 'help'.", verb)
		return false // Indicate this could be handled by LLM
	}
	return true
}

// HandleExamineFallback provides basic descriptions if the LLM is disabled
func (gs *GameState) handleExamineFallback(objectName string) {
	target := Item(strings.ToLower(objectName))

	// Check inventory first
	if _, have := gs.Inventory[target]; have {
		switch target {
		case ItemVoucher:
			gs.Message = "Voucher: Back says AISLE 13 // LAST SCAN."
		case ItemScanner:
			gs.Message = "Scanner: Frozen on Product ID: 8675309."
		case ItemNotebook:
			gs.Message = "Notebook: Mentions Brenda/Gary, OVERSTOCK alarm, Map."
		case ItemCard:
			gs.Message = "Card: Needs Key (safe) & Code ('OVERSTOCK' from Inventory)."
		case ItemInventory:
			gs.Message = "Inventory Sheet: OVERSTOCK -> Code: 4711."
		case ItemOverrideKey:
			gs.Message = "Key: Labeled 'Manual Override'."
		default:
			gs.Message = fmt.Sprintf("You look closely at the %s.", target)
		}
		return
	}

	// Check environment based on location (simplified)
	switch gs.Location {
	case LocRegister:
		gs.Message = "It's dark. Emergency lights glow. Main doors locked."
	case LocSecurityStation:
		gs.Message = "Dale's body is here. Monitors dark. Scanner nearby? Voucher clutched?"
	case LocLockerArea:
		gs.Message = "Dale's locker. Is it locked or open?"
	case LocManagersOffice:
		gs.Message = "Office: Desk, Corkboard (Card?), Safe."
	case LocLoadingDock:
		gs.Message = "Loading Dock: Breaker Panel, heavy door."
	default:
		gs.Message = "You look around."
	}

	// Add hints about visible items if not taken
	switch gs.Location {
	case LocSecurityStation:
		if !gs.Inventory[ItemVoucher] {
			gs.Message += " Dale clutches a voucher."
		}
		if !gs.Inventory[ItemScanner] {
			gs.Message += " A scanner lies nearby."
		}
	case LocLockerArea:
		if _, opened := gs.Clues["locker_opened"]; opened && !gs.Inventory[ItemNotebook] {
			gs.Message += " A notebook is inside the open locker."
		}
	case LocManagersOffice:
		if !gs.Inventory[ItemCard] {
			gs.Message += " A card is pinned to the board."
		}
		if !gs.Inventory[ItemInventory] {
			gs.Message += " An inventory sheet is on the desk."
		}
		if _, opened := gs.Clues["safe_opened"]; opened && !gs.Inventory[ItemOverrideKey] {
			gs.Message += " A key sits inside the open safe."
		}
	}
}

// IsCriticalUse determines if a 'use' command should be handled by Go logic.
func (gs *GameState) IsCriticalUse(input string) bool {
	lowerInput := strings.ToLower(input)
	// Check if using specific codes or keys at puzzle locations
	if gs.Location == LocLockerArea && strings.Contains(lowerInput, "8675309") {
		return true
	}
	if gs.Location == LocManagersOffice && strings.Contains(lowerInput, "4711") {
		return true
	}
	// For loading dock, check for key presence AND relevant input
	if gs.Location == LocLoadingDock {
		if _, hasKey := gs.Inventory[ItemOverrideKey]; hasKey {
			if strings.Contains(lowerInput, "key") || strings.Contains(lowerInput, "overstock") || strings.Contains(lowerInput, "683778625") {
				return true
			}
		}
	}
	return false
}

// handleGo moves the player to a new location if the destination is valid
func (gs *GameState) handleGo(destination string) {
	moved := false
	switch gs.Location {
	case LocRegister:
		if strings.Contains(destination, "security") || strings.Contains(destination, "electronics") || strings.Contains(destination, "back") || strings.Contains(destination, "scream") {
			gs.Location = LocSecurityStation
			gs.Message = "You hurry towards the back of the store, near the electronics section where the scream came from."
			moved = true
		}
	case LocSecurityStation:
		if strings.Contains(destination, "locker") {
			gs.Location = LocLockerArea
			gs.Message = "You move towards the nearby employee lockers, focusing on Dale's."
			moved = true
		} else if strings.Contains(destination, "office") || strings.Contains(destination, "manager") {
			gs.Location = LocManagersOffice
			gs.Message = "You head towards the small manager's office behind the customer service area."
			moved = true
		} else if strings.Contains(destination, "loading") || strings.Contains(destination, "dock") || strings.Contains(destination, "breaker") {
			if _, foundMap := gs.Clues["map_details"]; foundMap {
				gs.Location = LocLoadingDock
				gs.Message = "Following the crude map from Dale's notebook, you find the loading dock area."
				moved = true
			} else {
				gs.Message = "You aren't sure exactly where the loading dock or the specific breaker panel is."
			}
		} else if strings.Contains(destination, "register") || strings.Contains(destination, "front") {
			gs.Location = LocRegister
			gs.Message = "You head back towards the front registers."
			moved = true
		}
	case LocLockerArea:
		if strings.Contains(destination, "security") {
			gs.Location = LocSecurityStation
			gs.Message = "You step away from the lockers and back to the main security station area."
			moved = true
		}
	case LocManagersOffice:
		if strings.Contains(destination, "security") || strings.Contains(destination, "customer service") {
			gs.Location = LocSecurityStation // Assume security is closer/more relevant route back
			gs.Message = "You leave the manager's office, heading back towards the security station."
			moved = true
		} else if strings.Contains(destination, "loading") || strings.Contains(destination, "dock") || strings.Contains(destination, "breaker") {
			if _, foundMap := gs.Clues["map_details"]; foundMap {
				gs.Location = LocLoadingDock
				gs.Message = "You head from the office towards the loading dock, following the map's directions."
				moved = true
			} else {
				gs.Message = "You don't know the specific route to the loading dock from here without the map details."
			}
		}
	case LocLoadingDock:
		if strings.Contains(destination, "office") {
			gs.Location = LocManagersOffice
			gs.Message = "You head back towards the manager's office area."
			moved = true
		} else if strings.Contains(destination, "security") {
			gs.Location = LocSecurityStation
			gs.Message = "You move back towards the main security station area."
			moved = true
		}
	}

	if !moved && destination != "" {
		gs.Message = fmt.Sprintf("You can't find a way to '%s' from here, or you don't know where that is.", destination)
	} else if destination == "" {
		gs.Message = "Where do you want to go? (e.g., 'go security', 'go office')"
	}
	// If moved, message is set inside the switch case.
}

// handleTake attempts to take an item from the current location and add it to inventory
func (gs *GameState) handleTake(objectName string) {
	target := Item(strings.ToLower(objectName)) // Normalize object name
	taken := false
	alreadyHave := false

	// Check inventory first
	if _, have := gs.Inventory[target]; have {
		alreadyHave = true
	}

	if !alreadyHave {
		switch gs.Location {
		case LocSecurityStation:
			if target == ItemVoucher {
				gs.Inventory[ItemVoucher] = true
				taken = true
			} else if target == ItemScanner {
				gs.Inventory[ItemScanner] = true
				taken = true
			}
		case LocLockerArea:
			if target == ItemNotebook {
				if _, opened := gs.Clues["locker_opened"]; opened {
					gs.Inventory[ItemNotebook] = true
					taken = true
				} else {
					gs.Message = "The locker needs to be open first."
					return // Exit early
				}
			}
		case LocManagersOffice:
			if target == ItemCard {
				gs.Inventory[ItemCard] = true
				taken = true
			} else if target == ItemInventory {
				gs.Inventory[ItemInventory] = true
				taken = true
			} else if target == ItemOverrideKey {
				if _, opened := gs.Clues["safe_opened"]; opened {
					gs.Inventory[ItemOverrideKey] = true
					taken = true
				} else {
					gs.Message = "The safe needs to be open first."
					return // Exit early
				}
			}
		}
	}

	// Set feedback message
	if taken {
		gs.Message = fmt.Sprintf("You take the %s.", target)
	} else if alreadyHave {
		gs.Message = fmt.Sprintf("You already have the %s.", target)
	} else {
		gs.Message = fmt.Sprintf("You don't see a '%s' you can take here.", objectName)
	}
}

// handleCriticalUse handles 'use' commands identified as puzzle-critical
func (gs *GameState) handleCriticalUse(objectName, fullInput string) {
	lowerFullInput := strings.ToLower(fullInput) // Use lower case for comparisons

	switch gs.Location {
	case LocLockerArea:
		// Trying to use the scanner code on the locker
		// Check if the input string contains the code number
		if strings.Contains(lowerFullInput, "8675309") {
			gs.Message = "Enter the code for the locker:"
			gs.InputRequired = "locker_code"
			gs.CurrentInput = "" // Clear input buffer for code entry
		} else {
			// Guide the user if they typed 'use' but not the code
			gs.Message = "To use the code on the locker, try 'use 8675309'."
		}
	case LocManagersOffice:
		// Trying to use the safe code
		// Check if the input string contains the code number
		if strings.Contains(lowerFullInput, "4711") {
			gs.Message = "Enter the code for the safe:"
			gs.InputRequired = "safe_code"
			gs.CurrentInput = ""
		} else {
			// Guide the user if they typed 'use' but not the code
			gs.Message = "To use the code on the safe, try 'use 4711'."
		}
	case LocLoadingDock:
		// Trying to use the override key or the OVERSTOCK code
		hasKeyInInventory := gs.Inventory[ItemOverrideKey]
		mentionsKey := strings.Contains(lowerFullInput, "key")
		mentionsCode := strings.Contains(lowerFullInput, "overstock") || strings.Contains(lowerFullInput, "683778625")

		// Check if the player mentions the key or the code AND has the key in inventory
		if (mentionsKey || mentionsCode) && hasKeyInInventory {
			gs.Message = "You insert the Manual Override Key into the panel slot. Now, enter the activation code (OVERSTOCK or keypad numbers):"
			gs.InputRequired = "breaker_code"
			gs.CurrentInput = ""
		} else if (mentionsKey || mentionsCode) && !hasKeyInInventory {
			// Player tried to use key/code but doesn't have the key
			gs.Message = "You need the Manual Override Key first. Find it in the manager's safe and 'take' it."
		} else {
			// Player typed 'use' but didn't mention key or code specifically enough
			gs.Message = "To use the breaker panel, try 'use key' or 'use overstock' once you have the key."
		}
	default:
		// This case shouldn't be reached if IsCriticalUse is accurate
		gs.Message = fmt.Sprintf("You can't use '%s' in that specific way here.", objectName)
	}
}
