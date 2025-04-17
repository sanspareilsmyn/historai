# historai
<img src="asset/historai2.png" alt="Historai Demo 2" width="400">

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Uses LLM](https://img.shields.io/badge/Uses-LLM%20API-blueviolet.svg)]()

<a href="https://www.buymeacoffee.com/sanspareilsmyn"><img src="https://img.buymeacoffee.com/button-api/?text=Buy me a coffee&emoji=&slug=sanspareilsmyn&button_colour=FFDD00&font_colour=000000&font_family=Cookie&outline_colour=000000&coffee_colour=ffffff" /></a>

**historai is an LLM-powered Go CLI Tool that helps you find commands in your shell history and suggests new commands based on your needs, using natural language.**

## 💻  Demo
![demo.gif](asset%2Fdemo.gif)

## 🤔 Why historai?
Your shell history (e.g., `~/.zsh_history`, `~/.bash_history`) is a personal knowledge base, but accessing its full potential can be difficult:

*   **Finding Past Commands:**
    *   **Keyword Limitations:** Standard `history | grep` or `Ctrl+R` rely on you remembering specific keywords. What if you only remember *what the command did*?
    *   **"Tip-of-the-Tongue" Problem:** You *know* you've run a command before for a specific task (e.g., that tricky `ffmpeg` conversion or a complex `git rebase` sequence), but recalling the precise syntax is hard.
    *   **Inefficient Scrolling:** Manually scrolling through potentially thousands of shell history entries is impractical.
*   **Discovering New Commands:**
    *   **Need a Command You Haven't Used?:** Sometimes you need a command for a task you haven't performed before, or you only vaguely recall parts of it. How do you find the *right* command?
    *   **Leveraging Past Experience:** Your history contains clues about the types of tasks you perform. Can this context help generate relevant *new* command suggestions?

**historai addresses this by offering two core functionalities:**

1.  **`find`:** Search your **shell history** based on *description*, not just keywords. It helps you locate the **exact commands you have run before**.
    *   **🧠 AI-Powered Search:** Leverages Large Language Models (LLMs) to understand your natural language description of the command you're looking for *within your shell history*.
    *   **💬 Describe, Don't Just Recall Keywords:** Ask questions like "`find` the command where I listed docker containers sorted by size" or "`find` how I connected to the production server via ssh last week".
    *   **💡 Unlock Your History's Value:** Makes your own command history a truly searchable and useful resource.

2.  **`suggest`:** Get AI-generated command suggestions for a task description. It can use your **shell history** for context but aims to provide **relevant commands, even if you haven't executed them previously**.
    *   **🧠 AI-Powered Suggestions:** Uses LLMs to generate potentially useful shell commands based on your task description.
    *   **📜 Context-Aware (Optional):** Can leverage your past history to understand the types of tools and patterns you typically use, potentially leading to more relevant suggestions.
    *   **🚀 Go Beyond History:** Helps you discover or construct commands for new tasks or variations of old ones. Ask things like "`suggest` how to convert a video to a gif" or "`suggest` a command to find all *.log files modified today".

**Overall Benefits:**

*   **✅ Natural Language Interface:** Interact with your shell history and command discovery using plain English.
*   **📦 Single Binary CLI Tool:** Built with Go for easy installation via Homebrew and execution.
*   **🔐 API Key Driven:** Requires **your own LLM API key** (initially Google AI Studio) and internet access. This allows leveraging powerful AI models without complex local setup.

By using AI, historai makes finding specific past commands (`find`) and discovering relevant new commands (`suggest`) much faster and more intuitive.

## ✨ Features

*   **CLI Interface:** Simple commands: `historai find "..."` and `historai suggest "..."`.
*   **Find Past Commands (`find`):** Search your shell history using natural language descriptions to locate commands you have previously executed.
*   **Suggest Commands (`suggest`):** Get AI-generated command suggestions for a task description. It can use your shell history for context but can propose commands you haven't run before, helping you discover or construct new commands.
*   **LLM Integration:** Connects to the Google AI Studio API (Gemini models) to interpret your query and generate responses.
*   **API Key Management:** Reads your Google AI Studio API Key securely from the `GOOGLE_API_KEY` environment variable (config file support - *TODO*).
*   **Shell History Context:** Reads your shell history file (initially `~/.zsh_history` for Zsh) to provide the search space (`find`) or contextual background (`suggest`). **(Note: Only Zsh is supported in the initial version. Support for Bash, Fish, etc., is planned).**

## 🏗️ Architecture

`historai` runs as a **local Go CLI application**. When you run `historai find "..."` or `historai suggest "..."`:

1.  It parses your natural language query.
2.  It reads relevant parts of your **shell history file** (currently `~/.zsh_history` for Zsh users).
3.  It retrieves your configured Google AI Studio API key from the environment variable.
4.  It constructs a prompt asking the LLM to either:
    *   (`find`) Find shell history entries matching your description *within the provided history context*.
    *   (`suggest`) Suggest relevant shell commands based on your task description, potentially using the provided history as context for the types of tools or patterns you use.
5.  It sends this prompt via HTTPS to the Google AI Studio API service.
6.  It receives the AI-identified relevant history entries (`find`) or command suggestions (`suggest`).
7.  It formats and displays these findings in your terminal.

The AI interpretation and generation happen on Google's servers.

## 🗺️ Roadmap

We plan to expand `historai` with more features and support:

*   **✨ Interactive Command Selection & Execution:** Parse LLM output into selectable command options. Use a TUI library (like `survey` or `bubbletea`) to allow users to choose a command with arrow keys and Enter. Provide helper shell functions/aliases (using `eval`, `print -z`, etc.) to allow executing the selected command directly in the user's current shell.
*   **🤖 Additional LLM Support:** Integrate with other LLM APIs (e.g., OpenAI, Anthropic) and potentially allow choosing the backend.
*   **🐚 Broader Shell Support:** Add support for other popular shells like **Bash (`~/.bash_history`)** and **Fish (`~/.local/share/fish/fish_history`)**. *(High Priority)*
*   **🔧 Configuration File:** Implement a configuration file for more persistent settings (API keys, default models, history file location, interactive mode preferences).
*   **🚀 Performance Optimizations:** Improve history parsing and API interaction speed, potentially using caching or more efficient reading methods.
*   **⚙️ Advanced Filtering & Context:** Allow more granular filtering of history (date ranges, directory context) before sending to LLM. Improve how context is used for suggestions.
*   **💡 Enhanced Suggestions:** Improve the relevance and safety of suggested commands, potentially adding confirmation steps, explanations, or allowing feedback.
---
## 🚀 Getting Started (Initial Version - Zsh Only)

This guide helps you install and run `historai` using Homebrew. **Please note that the current version only supports the Zsh shell and its history file (`~/.zsh_history`).** Support for other shells is planned (see Roadmap).

**1. Prerequisites:**
*   **Homebrew:** Must be installed on your macOS or Linux system. ([Installation Guide](https://brew.sh/))
*   **Zsh Shell (Currently Required):** You must be using Zsh as your primary shell for this initial version.
*   **Google AI Studio API Key:** You **MUST** obtain an API key from [Google AI Studio](https://aistudio.google.com/). This is required for `historai` to function.
*   Access to your **Zsh history file (`~/.zsh_history`)**.

**2. Install historai using Homebrew:**
*   Open your terminal and run the following commands:
    ```bash
    # 1. Add the historai Tap (only needed once)
    brew tap sanspareilsmyn/historai

    # 2. Install historai
    brew install historai
    ```
*   **Updating historai:** To get the latest version in the future, run:
    ```bash
    brew upgrade historai
    ```

**3. Set Up Your API Key:**
*   `historai` requires your Google AI Studio API key to communicate with the LLM. Make it available via an environment variable:
    ```bash
    export GOOGLE_API_KEY="YOUR_GOOGLE_API_KEY_HERE"
    ```
*   **Important:** Replace the placeholder with your actual key. For persistence across terminal sessions, add this `export` line to your Zsh configuration file (`~/.zshrc`) and restart your shell or run `source ~/.zshrc`.

**4. Run historai:**
*   Once installed and the API key is set, you can run `historai` directly:

*   **Using `find` (Searching your Zsh history):**
    ```bash
    # Example: Find a specific git command from your shell history
    historai find "the git command I used to stash changes with a message 'fixing bug 123'"
    ```
    ```bash
    # Example: Find a command based on its effect
    historai find "show me how I listed files sorted by size last month"
    ```
    *   `historai find` will search your current shell history file (initially `~/.zsh_history` for Zsh) using the Gemini API and display matching entries *you previously executed*.

*   **Using `suggest` (Getting command suggestions):**
    ```bash
    # Example: Get suggestions for a common task
    historai suggest "how to convert a video file input.mp4 to an animated gif output.gif"
    ```
    ```bash
    # Example: Ask for a command to find specific files
    historai suggest "a command to find all python files modified in the last 24 hours"
    ```
    *   `historai suggest` will use the Gemini API to generate relevant command suggestions. **Always review suggested commands before executing them.**

---

## 🙌 Contributing

We welcome contributions! Please see `CONTRIBUTING.md` (TODO: Create this file) for details on how to contribute, especially regarding the features outlined in the Roadmap (like adding Bash and Fish support!). If you are contributing, you might need Go and Git installed locally to build and test changes.

---

## 📄 License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
