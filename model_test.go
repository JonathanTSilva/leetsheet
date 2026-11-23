package main

import (
	"bytes"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func testProblems() []Problem {
	mk := func(title string, kws ...string) Problem {
		return Problem{
			ProblemTitle: title,
			Keywords:     kws,
			Whiteboard:   "wb of " + title,
			DryRun:       "dry of " + title,
			TestCases:    "cases of " + title,
			Complexity: Complexity{
				Time:  ComplexityDetail{Notation: "O(N)", Justification: "one pass"},
				Space: ComplexityDetail{Notation: "O(1)", Justification: "constant"},
			},
			ManualSolution: "def manual(): pass",
			IaSolution:     "def ia(): pass",
		}
	}
	return []Problem{
		mk("Two Sum", "#Array", "#HashTable"),
		mk("Valid Parentheses", "#String", "#Stack"),
		mk("Binary Search", "#Array", "#BinarySearch"),
	}
}

func newReadyModel(t *testing.T) *model {
	t.Helper()
	m := newModel(testProblems())
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return m
}

// send feeds messages to the model like the Bubble Tea runtime does.
func send(m *model, msgs ...tea.Msg) {
	for _, msg := range msgs {
		m.Update(msg)
	}
}

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "ctrl+j":
		return tea.KeyMsg{Type: tea.KeyCtrlJ}
	case "ctrl+k":
		return tea.KeyMsg{Type: tea.KeyCtrlK}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func typeText(m *model, s string) {
	for _, r := range s {
		send(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
}

func TestInit(t *testing.T) {
	if cmd := newModel(testProblems()).Init(); cmd != nil {
		t.Errorf("Init() = non-nil cmd, want nil")
	}
}

func TestNewModel(t *testing.T) {
	m := newModel(testProblems())
	if m.state != searchView || m.activePane != leftPane {
		t.Errorf("initial state = %v/%v, want searchView/leftPane", m.state, m.activePane)
	}
	if got := len(m.list.Items()); got != 3 {
		t.Errorf("list has %d items, want 3", got)
	}
	if m.ready || m.quitting {
		t.Error("new model must not be ready or quitting")
	}
}

func TestViewBeforeReady(t *testing.T) {
	if got := newModel(testProblems()).View(); got != "Initializing..." {
		t.Errorf("View() = %q, want Initializing...", got)
	}
}

func TestWindowSizeMakesModelReady(t *testing.T) {
	m := newReadyModel(t)
	if !m.ready || m.windowWidth != 120 || m.windowHeight != 40 {
		t.Fatalf("ready=%v size=%dx%d", m.ready, m.windowWidth, m.windowHeight)
	}
	view := m.View()
	for _, want := range []string{"LeetCode Finder", "v4.3", "Two Sum", "Valid Parentheses", "search", "quit"} {
		if !strings.Contains(view, want) {
			t.Errorf("search view missing %q", want)
		}
	}
	// A second resize reuses the existing viewports.
	send(m, tea.WindowSizeMsg{Width: 80, Height: 30})
	if m.windowWidth != 80 || m.leftViewport.Width != 80/2-activePaneStyle.GetHorizontalFrameSize() {
		t.Errorf("resize not applied: width=%d left=%d", m.windowWidth, m.leftViewport.Width)
	}
}

func TestCtrlCQuits(t *testing.T) {
	m := newReadyModel(t)
	_, cmd := m.Update(key("ctrl+c"))
	if !m.quitting {
		t.Error("quitting = false after ctrl+c")
	}
	if cmd == nil {
		t.Fatal("ctrl+c returned no command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("ctrl+c command did not produce tea.QuitMsg")
	}
	if got := m.View(); got != "Bye!\n" {
		t.Errorf("View() = %q, want Bye!", got)
	}
}

func TestSearchViewNavigation(t *testing.T) {
	m := newReadyModel(t)
	send(m, key("ctrl+j"))
	if m.list.Index() != 1 {
		t.Errorf("index after ctrl+j = %d, want 1", m.list.Index())
	}
	send(m, key("ctrl+j"), key("ctrl+k"))
	if m.list.Index() != 1 {
		t.Errorf("index after ctrl+j, ctrl+k = %d, want 1", m.list.Index())
	}
	send(m, key("ctrl+k"))
	if m.list.Index() != 0 {
		t.Errorf("index after ctrl+k = %d, want 0", m.list.Index())
	}
}

func TestEnterOpensSelectedProblem(t *testing.T) {
	m := newReadyModel(t)
	send(m, key("ctrl+j"), key("enter"))
	if m.state != problemView || m.activeProblem == nil || m.activeProblem.ProblemTitle != "Valid Parentheses" {
		t.Fatalf("state=%v active=%+v, want problemView on Valid Parentheses", m.state, m.activeProblem)
	}
	view := m.View()
	for _, want := range []string{
		"Valid Parentheses", "O(N) Time | O(1) Space", "WHITEBOARD", "wb of Valid Parentheses",
		"DRY RUN", "TEST CASES", "TIME COMPLEXITY", "O(N): one pass", "SPACE COMPLEXITY", "O(1): constant",
		"MANUAL SOLUTION", "switch pane", "toggle solution",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("problem view missing %q", want)
		}
	}
}

func TestEnterOnEmptyListStaysInSearch(t *testing.T) {
	m := newReadyModel(t)
	m.list.SetItems(nil)
	send(m, key("enter"))
	if m.state != searchView || m.activeProblem != nil {
		t.Errorf("state=%v active=%v, want searchView and no active problem", m.state, m.activeProblem)
	}
}

func TestProblemViewKeys(t *testing.T) {
	m := newReadyModel(t)
	send(m, key("enter"))

	send(m, key("tab"))
	if m.activePane != rightPane {
		t.Errorf("activePane after tab = %v, want rightPane", m.activePane)
	}
	if !strings.Contains(m.View(), "Two Sum") {
		t.Error("view lost problem title after tab")
	}
	send(m, key("tab"))
	if m.activePane != leftPane {
		t.Errorf("activePane after second tab = %v, want leftPane", m.activePane)
	}

	send(m, key("c"))
	if m.solutionViewIndex != 1 || !strings.Contains(m.View(), "IA SOLUTION") {
		t.Errorf("solutionViewIndex=%d, view should show IA SOLUTION", m.solutionViewIndex)
	}
	send(m, key("c"))
	if m.solutionViewIndex != 0 || !strings.Contains(m.View(), "MANUAL SOLUTION") {
		t.Errorf("solutionViewIndex=%d, view should show MANUAL SOLUTION", m.solutionViewIndex)
	}

	// Scrolling keys reach both viewports without leaving the problem view.
	send(m, key("down"))
	if m.state != problemView {
		t.Error("scrolling left the problem view")
	}

	// A resize while reading a problem rebuilds the viewports.
	send(m, tea.WindowSizeMsg{Width: 100, Height: 30})
	if m.state != problemView || !strings.Contains(m.View(), "WHITEBOARD") {
		t.Error("resize in problem view lost content")
	}
}

func TestLeaveProblemView(t *testing.T) {
	for _, k := range []string{"esc", "/"} {
		t.Run(k, func(t *testing.T) {
			m := newReadyModel(t)
			send(m, key("enter"), key(k))
			if m.state != searchView {
				t.Fatalf("state after %q = %v, want searchView", k, m.state)
			}
			if got := len(m.list.Items()); got != 3 {
				t.Errorf("list has %d items after leaving, want all 3", got)
			}
			if !strings.Contains(m.View(), "LeetCode Finder") {
				t.Error("search view not rendered after leaving problem view")
			}
		})
	}
}

func TestKeywordSearchFlow(t *testing.T) {
	m := newReadyModel(t)
	send(m, key("/"))
	if m.list.FilterState().String() != "filtering" {
		t.Fatalf("filter state = %v, want filtering", m.list.FilterState())
	}
	view := m.View()
	for _, want := range []string{"Search:", "navigate", "cancel"} {
		if !strings.Contains(view, want) {
			t.Errorf("filtering view missing %q", want)
		}
	}

	typeText(m, "#array")
	if got := listTitles(m.list); len(got) != 2 || got[0] != "Two Sum" || got[1] != "Binary Search" {
		t.Errorf("#array filtered to %v, want [Two Sum Binary Search]", got)
	}
	typeText(m, " hashtable")
	if got := listTitles(m.list); len(got) != 1 || got[0] != "Two Sum" {
		t.Errorf("#array hashtable filtered to %v, want [Two Sum]", got)
	}
}

func TestUpdateIgnoresUnrelatedMessages(t *testing.T) {
	m := newReadyModel(t)
	send(m, struct{}{})
	if m.state != searchView || m.quitting {
		t.Error("unrelated message changed state")
	}
}

func TestHeaderViewNarrowWindow(t *testing.T) {
	m := newReadyModel(t)
	m.windowWidth = 5 // narrower than the title and info boxes: the rule line must not panic
	if !strings.Contains(m.headerView(), "LeetCode Finder") {
		t.Error("header missing title in a narrow window")
	}
}

func TestRun(t *testing.T) {
	t.Run("quits on ctrl+c", func(t *testing.T) {
		var out bytes.Buffer
		if err := run(strings.NewReader("\x03"), &out, "testdata/valid.json"); err != nil {
			t.Fatalf("run returned %v", err)
		}
		if out.Len() == 0 {
			t.Error("run wrote nothing to the output")
		}
	})
	t.Run("uses the bundled database when the default file is missing", func(t *testing.T) {
		var out bytes.Buffer
		if err := run(strings.NewReader("\x03"), &out, t.TempDir()+"/missing.json"); err != nil {
			t.Fatalf("run returned %v", err)
		}
	})
	t.Run("loads the repository problems file", func(t *testing.T) {
		var out bytes.Buffer
		if err := run(strings.NewReader("\x03"), &out, "problems.json"); err != nil {
			t.Fatalf("run returned %v", err)
		}
	})
	t.Run("reports a malformed problems file", func(t *testing.T) {
		err := run(strings.NewReader(""), &bytes.Buffer{}, "testdata/malformed.json")
		if err == nil || !strings.Contains(err.Error(), "testdata/malformed.json") {
			t.Fatalf("error = %v, want one naming the file", err)
		}
	})
}
