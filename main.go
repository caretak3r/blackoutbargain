package main

import (
	"fmt"
	"log"
	"os"

	"blackoutbargain/llm"
	"blackoutbargain/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Main Function ---
func main() {
	// Set up logging
	logFile, err := os.OpenFile("blackout_bargain.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Starting Blackout Bargain...")

	// Initialize LLM client
	llmClient := llm.New()
	defer func() {
		if llmClient != nil {
			err := llmClient.Close()
			if err != nil {
				log.Printf("Error closing LLM client: %v", err)
			}
		}
	}()

	// Check if LLM initialization failed critically
	if os.Getenv("GEMINI_API_KEY") != "" && llmClient != nil && !llmClient.Enabled {
		fmt.Println("Error initializing LLM - check API key and permissions. Exiting.")
		log.Println("LLM client initialization failed critically.")
		os.Exit(1)
	}

	// Initialize the TUI model
	m := tui.New(llmClient)

	// Create and run the Bubble Tea program
	// Using AltScreen helps restore the terminal state on exit
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()) // Enable mouse if needed later

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		log.Printf("Error running Bubble Tea program: %v", err)
		os.Exit(1)
	}

	log.Println("Blackout Bargain finished.")
}
