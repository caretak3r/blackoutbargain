package llm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"blackoutbargain/game"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client manages interactions with the LLM service (Gemini)
type Client struct {
	APIKey         string
	GeminiClient   *genai.Client
	GeminiModel    *genai.GenerativeModel
	Context        context.Context
	Enabled        bool
	LastPromptSent string
}

// New initializes a new LLM client
func New() *Client {
	client := &Client{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Context: context.Background(),
		Enabled: false,
	}

	// Initialize LLM Client if key exists
	if client.APIKey != "" {
		geminiClient, err := genai.NewClient(client.Context, option.WithAPIKey(client.APIKey))
		if err != nil {
			log.Printf("Error creating LLM client: %v. LLM disabled.", err)
			return client
		}

		client.GeminiClient = geminiClient
		// Use a fast and capable model like Flash
		client.GeminiModel = geminiClient.GenerativeModel("gemini-2.0-flash")
		// Set safety settings to block harmful content
		client.GeminiModel.SafetySettings = []*genai.SafetySetting{
			{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
			{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
			{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockMediumAndAbove},
			{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockMediumAndAbove},
		}
		client.Enabled = true
		log.Println("LLM Client Initialized.")
	}

	return client
}

// Close cleans up the LLM client
func (c *Client) Close() error {
	if c.GeminiClient != nil {
		return c.GeminiClient.Close()
	}
	return nil
}

// GenerateResponse calls the LLM with the provided player input and game state
func (c *Client) GenerateResponse(playerInput string, gameState *game.GameState) (string, error) {
	if !c.Enabled || c.GeminiClient == nil || c.GeminiModel == nil {
		return "LLM support is not available. Using basic descriptions.", nil
	}

	// Construct the prompt
	prompt := c.buildPrompt(playerInput, gameState)
	c.LastPromptSent = prompt

	// Send the prompt to the LLM
	resp, err := c.GeminiModel.GenerateContent(c.Context, genai.Text(prompt))
	if err != nil {
		log.Printf("LLM API call error: %v", err)
		return "", fmt.Errorf("API request failed: %w", err)
	}

	// Process the response
	if len(resp.Candidates) > 0 &&
		resp.Candidates[0].Content != nil &&
		len(resp.Candidates[0].Content.Parts) > 0 {
		generatedText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
		// Simple cleanup: remove potential markdown emphasis added by LLM if not desired
		generatedText = strings.ReplaceAll(generatedText, "*", "")
		return generatedText, nil
	}

	// Handle cases where the LLM response is empty or malformed
	log.Printf("LLM Warning: Received empty or unexpected response format: %+v", resp)
	return "The situation doesn't seem to change.", nil
}

// buildPrompt creates the detailed prompt for the Gemini API
func (c *Client) buildPrompt(playerInput string, gameState *game.GameState) string {
	var sb strings.Builder

	// --- System Instructions / Context ---
	sb.WriteString("You are the narrator for 'Blackout Bargain', a text adventure game. The player is trapped in a dark Superstore after a power failure killed the lights and locked the doors. Dale, the security guard, was found dead (puncture wound, neck). The player is with Brenda (stocker) and Gary (manager). Goal: Escape.")
	sb.WriteString(" Core Puzzle Path: Find Dale (security station) -> Get Voucher (from Dale) & Scanner -> Use scanner code (8675309) on Locker -> Get Notebook -> Read notebook (mentions Brenda/Gary, 'OVERSTOCK' silent alarm, map to breaker panel needing manager key) -> Go to Manager's Office -> Get Emergency Card (mentions key in safe, code is 'OVERSTOCK' from inventory sheet) -> Get Inventory Sheet -> Use sheet ('OVERSTOCK' -> code 4711) on Safe -> Get Override Key -> Go to Loading Dock Breaker Panel (from map) -> Use Key & 'OVERSTOCK' (or 683778625) on panel -> Unlock door -> Confrontation (Gary is killer) -> Escape.")
	sb.WriteString(" Rules: Narrate atmospheric outcomes of player actions based on current state. Stick to the established items, characters, and puzzle path. Do NOT invent new major items, characters, bypasses, or solutions. If the player tries something irrelevant or impossible, explain why it fails or gently guide them back to relevant actions based on their known clues/location. Be concise but descriptive. Keep the tone tense/mysterious.")

	// --- Current Game State ---
	sb.WriteString("\n\n--- Current State ---")
	sb.WriteString(fmt.Sprintf("\nLocation: %s (%s)", gameState.GetLocationName(), gameState.GetLocationDescription()))

	// Inventory
	invItems := []string{}
	for item := range gameState.Inventory {
		invItems = append(invItems, string(item))
	}
	if len(invItems) == 0 {
		sb.WriteString("\nInventory: Empty")
	} else {
		sb.WriteString(fmt.Sprintf("\nInventory: %s", strings.Join(invItems, ", ")))
	}

	// Clues
	clueItems := []string{}
	for key, val := range gameState.Clues {
		clueItems = append(clueItems, fmt.Sprintf("%s: %s", key, val))
	}
	if len(clueItems) > 0 {
		sb.WriteString(fmt.Sprintf("\nKnown Clues: %s", strings.Join(clueItems, "; ")))
	}

	// --- Player's Action ---
	sb.WriteString("\n\n--- Player Action ---")
	sb.WriteString(fmt.Sprintf("\n%s", playerInput))

	// --- Instruction for LLM ---
	sb.WriteString("\n\n--- Narrator Response ---")
	sb.WriteString("\nDescribe the result: ") // Prompt LLM to generate the narrative

	return sb.String()
}
