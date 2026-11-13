<!-- brand:header:start -->
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="assets/brand/header-5fd7cb2c-dark.svg">
  <img alt="leetsheet: A TUI-based LeetCode cheatsheet for quick problem walkthroughs in your terminal" src="assets/brand/header-5fd7cb2c-light.svg" width="100%">
</picture>
<!-- brand:header:end -->

<!-- badges:start -->
[![Status](https://img.shields.io/badge/status-active-brightgreen?style=flat-square)](https://www.repostatus.org/#active)
[![CI](https://img.shields.io/github/actions/workflow/status/JonathanTSilva/leetsheet/ci.yml?branch=main&style=flat-square&label=ci)](https://github.com/JonathanTSilva/leetsheet/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/JonathanTSilva/leetsheet/graph/badge.svg)](https://codecov.io/gh/JonathanTSilva/leetsheet)
[![Quality](https://sonarcloud.io/api/project_badges/measure?project=JonathanTSilva_leetsheet&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=JonathanTSilva_leetsheet)
[![Scorecard](https://img.shields.io/ossf-scorecard/github.com/JonathanTSilva/leetsheet?style=flat-square&label=scorecard)](https://scorecard.dev/viewer/?uri=github.com/JonathanTSilva/leetsheet)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)
[![Open in GitHub Codespaces](https://img.shields.io/badge/codespaces-open-181717?style=flat-square&logo=github)](https://codespaces.new/JonathanTSilva/leetsheet)
<!-- badges:end -->

## What it is

`leetsheet` is a command-line application for developers who want a fast, offline way to review LeetCode problems. Instead of navigating browser tabs, you get problem details, complexity analysis, whiteboard explanations and multiple solution implementations directly in your terminal, from a local `problems.json` database of over 1200 problems.

![leetsheet demo](assets/gif/demo.gif)

For each problem you get:

- **Whiteboard Explanation:** The high-level strategy and data structures involved.
- **Dry Run:** A step-by-step example to solidify understanding.
- **Test Cases:** Critical edge cases to consider.
- **Complexity Analysis:** Time and Space complexity with clear justifications.
- **Dual Solutions:** A manually crafted, commented solution and an AI-generated one for comparison.

Features:

- **Instant Search:** Filter through problems by title as you type.
- **Tag-Based Filtering:** Use `#tag` syntax (e.g., `#Array #HashTable`) to find problems by topic.
- **Comprehensive Problem View:** A two-pane layout with problem details on the left and code solutions on the right.
- **Toggleable Solutions:** Switch between the manual and AI-generated solution with the `c` key.
- **Vim-like Keybindings:** Navigate lists and scroll content efficiently.
- **Responsive TUI:** The layout adapts to your terminal window size.
- **Syntax Highlighting:** Code solutions are highlighted for readability.
- **Offline First:** All data is read from a local `problems.json` file.

## Quick start

You need Go (version 1.25 or newer, see `go.mod`).

```bash
git clone https://github.com/JonathanTSilva/leetsheet.git
cd leetsheet
go build
./leetsheet
```

`go build` creates an executable named `leetsheet` in the current directory and expects `problems.json` in the project root. To run it from anywhere, move the binary to a directory in your `PATH`:

```bash
sudo mv leetsheet /usr/local/bin/
```

## Usage

Run `leetsheet` from your terminal. The application is fully keyboard-driven.

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

### Adding new problems

All problem data is stored in `problems.json`. To add a problem, append a JSON object with this structure to the array in that file:

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

JSON requires strings to be escaped (newlines as `\n`, quotes as `\"`). The `code_to_json.py` utility script does this for you:

1. Optionally install `pyperclip` to copy the output to your clipboard automatically: `pip install pyperclip`.
2. Open `code_to_json.py` and paste your multi-line code solution between the triple quotes.
3. Run `python code_to_json.py`.
4. Paste the printed single-line string as the value of `ia_solution` or `manual_solution` in `problems.json`.

## Docs

The repository layout and the engineering standard it follows are described in the [handbook](https://github.com/JonathanTSilva/.github/blob/main/handbook/README.md). Instructions for AI coding agents are in [AGENTS.md](AGENTS.md).

## Built with

- [Go](https://golang.org/): the core language.
- [Bubble Tea](https://github.com/charmbracelet/bubbletea): the TUI framework.
- [Lipgloss](https://github.com/charmbracelet/lipgloss): terminal layouts.
- [Glamour](https://github.com/charmbracelet/glamour): Markdown and syntax-highlighted code rendering.

Thanks to the team at [Charm](https://charm.sh/) for this suite of tools.

## Contributing

See the [contributing guide](https://github.com/JonathanTSilva/.github/blob/main/CONTRIBUTING.md). Branches are `<type>/<issue>-<slug>`, commits and PR titles follow Conventional Commits, and PRs are squash-merged.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

<!-- brand:footer:start -->
<p align="right"><a href="https://www.jontobias.com"><picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/JonathanTSilva/.github/main/assets/brand/jontobias-dark.svg">
  <img alt="jontobias.com" src="https://raw.githubusercontent.com/JonathanTSilva/.github/main/assets/brand/jontobias-light.svg" width="200">
</picture></a></p>
<!-- brand:footer:end -->
