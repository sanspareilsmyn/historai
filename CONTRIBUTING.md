# Contributing to historai

First off, thank you for considering contributing to historai! 🎉 We welcome any help to make this project better by enhancing its ability to **find past commands** and **suggest new ones** using shell history and AI. Whether it's reporting a bug, proposing a new feature (especially from our Roadmap!), improving documentation, or writing code, your contribution is valuable.

This document provides guidelines for contributing to historai. Please take a moment to review it.

## Code of Conduct

While we don't have a formal Code of Conduct yet, we expect all contributors to interact respectfully and constructively. Please be kind and considerate in all discussions and contributions. Harassment or exclusionary behavior will not be tolerated.

## How Can I Contribute?

There are many ways to contribute:

*   **🐛 Reporting Bugs:** If you find a bug (e.g., incorrect history parsing (initially for Zsh), errors interacting with the Google AI API, unexpected `find` results, irrelevant `suggest` outputs), please open an issue on GitHub. Provide as much detail as possible, including the command you ran (`find` or `suggest`), the query, relevant parts of your shell history (if possible and safe to share), expected behavior, and actual behavior/error messages.
*   **✨ Suggesting Enhancements:** Open an issue to discuss new features or improvements, especially those listed in the [Roadmap](#-roadmap) section of the README. Prime examples include:
    *   **Adding support for other shells (Bash, Fish)**
    *   Integrating new LLM providers
    *   Improving the relevance or safety of `suggest` results
    *   Implementing a configuration file
*   **📝 Improving Documentation:** Pull Requests (PRs) for documentation improvements (README, code comments, setup instructions, examples for both `find` and `suggest`, etc.) are always welcome!
*   **💻 Writing Code:**
    1.  **Discuss First (for significant changes):** It's often best to open an issue to discuss major changes (like adding a new LLM provider or shell support) before starting work, to ensure alignment.
    2.  **Follow the Workflow:** See the "Contribution Workflow" section below.

## Getting Started (Development Setup)

To contribute code, you'll need a local development environment. Unlike the user installation via Homebrew, development requires Go and Git. `historai` runs locally, reads your shell history file, and interacts with the Google AI Studio API.

1.  **Prerequisites:**
    *   **Go:** Version 1.22+ recommended (check `go.mod`).
    *   **Git:** For cloning the repository and managing changes.
    *   **Zsh Shell:** While the project aims for broader support, **current development and testing heavily rely on Zsh and `~/.zsh_history`**. Having Zsh installed with a populated history file is essential for testing the initial implementation.
    *   **Google AI Studio API Key:** You **MUST** obtain an API key from [Google AI Studio](https://aistudio.google.com/). This is required for the core functionality (`find` and `suggest`).
    *   (Recommended) `pre-commit` tool (install via `pip install pre-commit` or `brew install pre-commit`) if you intend to commit code.

2.  **Fork & Clone:**
    *   Fork the repository on GitHub.
    *   Clone your fork locally:
        ```bash
        git clone https://github.com/<your-username>/historai.git # Replace with your repo path
        cd historai
        ```

3.  **Prepare Test Data & Environment:**
    *   **Set API Key:** Make your Google AI Studio API key available via an environment variable (this is how the app currently reads it):
        ```bash
        export GOOGLE_API_KEY="YOUR_GOOGLE_API_KEY_HERE"
        # Consider adding to ~/.zshrc for convenience during development
        ```
    *   **Prepare Shell History (Current Focus: Zsh):** Ensure you have a `~/.zsh_history` file that `historai` can read. For testing specific `find` and `suggest` scenarios, you might want to:
        *   Temporarily add specific commands to your history (`fc -R` might be needed in Zsh to reload it).
        *   Have a varied history to test how context influences `suggest`.
        *   *Carefully* create a copy or a sample `test_history` file if you don't want to use your live history or need specific reproducible test cases. You might need to modify the code temporarily or add a flag to point to a test file during development.

4.  **Set up Pre-commit Hooks:**
    *   If installed, run `pre-commit install` in the repository root.

5.  **Dependencies & Local Build:**
    *   Fetch Go dependencies:
        ```bash
        go mod tidy
        ```
    *   Build the `historai` binary **locally**:
        ```bash
        go build -o historai ./cmd/historai
        ```

6.  **Running & Testing:**
    *   **1. Ensure API Key is Set:** Verify the `GOOGLE_API_KEY` environment variable is exported in your current shell session (`echo $GOOGLE_API_KEY`).
    *   **2. Run historai Locally:** Execute the `historai` binary you built with the `find` or `suggest` command.
        ```bash
        # Example: Test finding a command
        ./historai find "how did I list docker images recently?"

        # Example: Test suggesting a command
        ./historai suggest "command to recursively find and delete node_modules folders"

        # Example: Test find with a more specific query
        ./historai find "the command where I used grep to find 'error' in log files"
        ```
    *   **3. Check Output:**
        *   **`find` Result:** You should see relevant command(s) from your shell history file (currently Zsh) displayed.
        *   **`suggest` Result:** You should see one or more command suggestions. Evaluate their relevance and correctness.
        *   **Check for Errors:** Look for any error messages related to API key issues, API rate limits, history file reading problems, or unexpected responses from the LLM.
    *   **4. Run Unit Tests:** Run Go unit tests (these should ideally mock external dependencies like the history file reader or the API client where possible).
        ```bash
        go test ./...
        ```

7.  **Stopping the Environment:**
    *   The `historai` CLI command typically exits after running.
    *   You might want to unset the API key environment variable if you set it temporarily: `unset GOOGLE_API_KEY`.

## Contribution Workflow

1.  **Fork & Branch:** Fork the repo and create a descriptive branch from `main`: `git checkout -b feat/describe-your-feature` or `fix/handle-bash-timestamps`.
2.  **Develop:** Make your code changes. Add or update unit tests. Update documentation (README, comments) accordingly.
3.  **Test:** Run `go test ./...`. Perform manual testing using the steps in "Running & Testing" with your API key and shell history (currently Zsh). Test both `find` and `suggest` where applicable.
4.  **Pre-commit:** Ensure pre-commit checks pass (`pre-commit run --all-files` if needed).
5.  **Commit:** Use Conventional Commits format (e.g., `git commit -m "feat: add suggest command logic"` or `fix(parser): correctly handle multiline zsh history entries"` or `feat(shell): add initial bash history parsing`).
6.  **Push:** Push your branch to your fork: `git push origin feat/your-branch-name`.
7.  **Open a Pull Request (PR):** Create a PR to the main repository's `main` branch. Provide a clear title/description, link related issues, explain your changes, and mention how you tested them (including manual tests for `find`/`suggest`).

## Pull Request Process

1.  **Review:** Maintainer review and feedback.
2.  **CI Checks:** Automated checks (linters, tests - TODO: Setup CI) must pass.
3.  **Discussion & Iteration:** Address feedback and update your branch by pushing new commits.
4.  **Approval & Merge:** Once approved and checks pass, the PR will be merged into `main`.

## Questions?

Feel free to open an issue on GitHub to ask questions or discuss potential contributions, especially regarding roadmap items like adding new shell support!

Thank you for contributing to historai! 🙏
