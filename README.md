# leetsheet

<div align="center">
<br>
<a href="#installation"><kbd> <br> Installation <br> </kbd></a>&ensp;&ensp;
<a href="#updating"><kbd> <br> Updating <br> </kbd></a>&ensp;&ensp;
<a href="#themes"><kbd> <br> Themes <br> </kbd></a>&ensp;&ensp;
<a href="#styles"><kbd> <br> Styles <br> </kbd></a>&ensp;&ensp;
<a href="#keybindings"><kbd> <br> Styles <br> </kbd></a>&ensp;&ensp;
<a href="CONTRIBUTING.md"><kbd> <br> Contributing <br> </kbd></a>&ensp;&ensp;
</div><br><br>

**A TUI-based LeetCode Cheatsheet for Quick Problem Walkthroughs in Your Terminal.**

`leetsheet` is a command-line application designed for developers who want a fast, offline, and focused way to review LeetCode problems. Instead of navigating browser tabs, get instant access to problem details, complexity analysis, whiteboard explanations, and multiple solution implementations directly in your terminal.

## Demo

<p align="center"><img src="/assets/gif/demo.gif?raw=true"/></p>

## About The Project

This isn't just another interview prep tool; it's a personal, offline-first, and blazingly fast cheatsheet. The goal is to minimize context switching and keep you in the "flow state" of your terminal. With a local JSON database of over 1200 problems, `leetsheet` provides a structured and comprehensive walkthrough for each, helping you internalize patterns and solutions without the distractions of the web.

Key information for each problem includes:

- **Whiteboard Explanation:** The high-level strategy and data structures involved.
- **Dry Run:** A step-by-step example to solidify understanding.
- **Test Cases:** Critical edge cases to consider.
- **Complexity Analysis:** Time and Space complexity with clear justifications.
- **Dual Solutions:** A manually crafted, commented solution and an AI-generated one for comparison.

## Features

- **Instant Search:** Filter through problems by title as you type.
- **Tag-Based Filtering:** Use `#tag` syntax (e.g., `#Array #HashTable`) to find problems by topic.
- **Comprehensive Problem View:** A clean, two-pane layout displaying problem details on the left and code solutions on the right.
- **Toggleable Solutions:** Instantly switch between the manually-written solution and an AI-generated one with the `c` key.
- **Vim-like Keybindings:** Navigate lists and scroll content efficiently.
- **Responsive TUI:** The layout adapts to your terminal window size.
- **Syntax Highlighting:** Code solutions are beautifully highlighted for readability.
  \-- **Offline First:** All data is read from a local `problems.json` file, making it available anywhere, anytime.

## Installation

### Prerequisites

You need to have Go (version 1.18 or newer) installed on your system.

```bash
# Check if Go is installed
go version
```

### Steps

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/leetsheet.git
    ```

2.  **Navigate to the project directory:**

    ```bash
    cd leetsheet
    ```

3.  **Get the data:**
    Ensure you have the `problems.json` file in the root of the project directory.

4.  **Build the binary:**

    ```bash
    go build
    ```

    This will create an executable file named `leetsheet` in the current directory.

5.  **(Optional) Add to your PATH:**
    For easy access from anywhere, move the binary to a directory in your system's PATH.

    ```bash
    sudo mv leetsheet /usr/local/bin/
    ```

## Usage

Simply run the executable from your terminal:

```bash
./leetsheet
```

If you moved it to your PATH, you can just run:

```bash
leetsheet
```

### Keybindings

The application is fully keyboard-driven.

| View        | Key             | Action                                |
| ----------- | --------------- | ------------------------------------- |
| **Search**  | `ctrl+c`        | Quit the application                  |
| (Default)   | `/`             | Start filtering (enter search mode)   |
|             | `ctrl+j/k`      | Navigate up/down the problem list     |
|             | `enter`         | Select a problem and view its details |
|             | `esc`           | Clear the current filter/search       |
|             | `#tag`          | Type `#` followed by a tag to filter  |
|             |                 |                                       |
| **Problem** | `esc` / `/`     | Go back to the search view            |
| (Details)   | `tab`           | Switch focus between left/right panes |
|             | `c`             | Toggle between Manual and IA solution |
|             | `↑`/`↓`/`j`/`k` | Scroll the content in the active pane |

## Adding New Problems

All problem data is stored in `problems.json`. To add a new problem, you need to append a new JSON object to the array in this file.

### Data Structure

Each problem object must follow this structure:

```json
{
  "title": "Your Problem Title",
  "link": "https://leetcode.com/problems/your-problem-name/",
  "keywords": ["#Tag1", "#Tag2"],
  "complexity": {
    "time": {
      "notation": "O(N)",
      "justification": "Explanation for time complexity."
    },
    "space": {
      "notation": "O(1)",
      "justification": "Explanation for space complexity."
    }
  },
  "whiteboard": "A high-level explanation of the approach.",
  "dry_run": "A step-by-step walk-through with an example.",
  "test_cases": "A list of important edge cases to consider.",
  "ia_solution": "The AI-generated code solution string.",
  "manual_solution": "The manually written and commented code solution string."
}
```

### Utility Script for Code Formatting

JSON requires strings to be properly escaped (e.g., newlines as `\n`, quotes as `\"`). To make it easy to add new code solutions, you can use the provided `code_to_json.py` utility script.

**Usage:**

1.  **(Optional) Install `pyperclip` to automatically copy the output to your clipboard:**

    ```bash
    pip install pyperclip
    ```

2.  **Paste your code:**
    Open `code_to_json.py` and paste your multi-line code solution inside the triple quotes.

    ```python
    import json
    import pyperclip

    # Paste your code between the triple quotes
    code = """class Solution:
        def yourCode(self, s: str) -> int:
            # Your implementation here
            return 0
    """

    # ... rest of the script
    ```

3.  **Run the script:**

    ```bash
    python code_to_json.py
    ```

4.  **Paste the output:**
    The script will print a JSON-safe, single-line string. Paste this string as the value for `ia_solution` or `manual_solution` in your `problems.json` file. If `pyperclip` is installed, it will also be in your clipboard.

## Built With

- [Go](https://golang.org/) - The core language.
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The TUI framework that makes this possible.
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For beautiful, stylish terminal layouts.
- [Glamour](https://github.com/charmbracelet/glamour) - For rendering Markdown and syntax-highlighted code blocks.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

A huge thank you to the team at [Charm](https://charm.sh/) for creating an amazing suite of tools for building beautiful TUIs.
