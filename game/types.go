package game

// --- Game State Definitions ---

// Location represents a distinct area within the game world
type Location int

const (
	LocRegister Location = iota
	LocSecurityStation
	LocLockerArea
	LocManagersOffice
	LocLoadingDock
	LocEscaped
)

// Item represents a collectible object in the game
type Item string

const (
	ItemVoucher     Item = "crumpled employee discount voucher"
	ItemScanner     Item = "Dale's handheld scanner"
	ItemNotebook    Item = "small notebook"
	ItemCard        Item = "laminated emergency procedure card"
	ItemInventory   Item = "daily inventory printout"
	ItemOverrideKey Item = "Manual Override Key"
)

// GameState represents the current state of the game
type GameState struct {
	// Game state
	Location      Location
	Inventory     map[Item]bool
	Clues         map[string]string // Store discovered codes/facts
	GameOver      bool
	Message       string // Feedback/narrative display
	CurrentInput  string
	InputRequired string // Specific input needed: "locker_code", "safe_code", "breaker_code"
}

// NewGameState initializes and returns a new game state
func NewGameState() *GameState {
	return &GameState{
		Location:  LocRegister,
		Inventory: make(map[Item]bool),
		Clues:     make(map[string]string),
		GameOver:  false,
	}
}
