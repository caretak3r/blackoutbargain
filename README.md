# Blackout Bargain 🛒🔦

## Overview

Blackout Bargain is an interactive text-based adventure game built with Go and the Charm TUI libraries.

The lights flicker and die in the Superstore. Thunder cracks outside. You're trapped inside with a few colleagues, the electronic doors sealed shut by the power outage. When the security guard is found dead, seemingly murdered, you realize the blackout is the least of your worries. Explore the darkened store, gather clues, interact with items and the environment, solve puzzles, and uncover the truth to find a way out before the killer finds you.

The game follows the narrative and puzzle progression detailed in `prompt.md`.

## ✨ Features

*   **Rich Interactive Terminal UI:** Built using the Charm stack (Bubble Tea, Lipgloss, Huh) for an engaging text-based experience.
*   **Narrative-Driven Gameplay:** Unravel a mystery by following the story, examining clues, and making choices.
*   **Puzzle Solving:** Interact with items, decipher codes, and overcome obstacles to progress.
*   **Go Backend:** Written entirely in Go.
*   **LLM Integration:** Leverages the Google Gemini API (requires an API key). *Note: The specific role of the LLM in the gameplay loop might need further clarification.*
*   **Logging:** Tracks game progress and potential errors in `blackout_bargain.log`.

## 🛠️ Technologies Used

*   **Language:** Go
*   **TUI Framework:** [Charm](https://charm.sh/)
    *   [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Application Framework)
    *   [Lipgloss](https://github.com/charmbracelet/lipgloss) (Styling)
    *   [Huh](https://github.com/charmbracelet/huh) (Interactive Forms/Prompts)
    *   *(Potentially [Harmonica](https://github.com/charmbracelet/harmonica) for physics/animations)*
*   **LLM:** [Google Gemini API](https://ai.google.dev/)

## ⚙️ Setup & Installation

1.  **Prerequisites:**
    *   Go (check `go.mod` for specific version requirements if any, otherwise latest stable is recommended).
    *   A Google Gemini API Key.
2.  **Clone:**
    ```bash
    # Replace with your repository URL if applicable
    git clone https://your-repository-url/blackout-bargain.git
    cd blackout-bargain
    ```
3.  **Set API Key:**
    The application requires your Google Gemini API key to be available as an environment variable:
    ```bash
    export GEMINI_API_KEY="YOUR_API_KEY_HERE"
    ```
    *(The application checks for this key on startup and may exit if it's missing or invalid).*
4.  **Build:**
    ```bash
    go build -o blackoutbargain main.go
    ```
    *(This creates the `blackoutbargain` executable).*

## ▶️ How to Play

1.  Ensure the `GEMINI_API_KEY` environment variable is set.
2.  Run the compiled executable from your terminal:
    ```bash
    ./blackoutbargain
    ```
3.  The game will start in your terminal. Follow the narrative prompts.
4.  Use the menus and input fields provided by the TUI to interact with the game world, examine objects, talk to characters (if implemented), and solve the puzzles outlined in the story.
5.  Your objective is to solve Dale's murder and escape the Superstore.

---

This content provides a comprehensive overview. You can adjust the details, especially regarding the LLM's exact role, add licensing information, or include screenshots/gifs once the TUI is more developed.