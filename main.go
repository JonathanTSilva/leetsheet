package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// *========( STYLES )=========*
var (
	docStyle            = lipgloss.NewStyle().Margin(0, 2)
	inactivePaneStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
	activePaneStyle     = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("#AD58B4"))
	paneTitleStyleFirst = lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("252"))
	paneTitleStyle      = lipgloss.NewStyle().Bold(true).Underline(true).MarginTop(1).Foreground(lipgloss.Color("252"))

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).BorderForeground(lipgloss.Color("240"))
	}()
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	footerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("250")).
			Padding(0, 1)

	keyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Padding(0, 1).
			MarginLeft(1)
)

// *========( DATA MODELS )=========*
type Problem struct {
	ProblemTitle   string     `json:"title"`
	Link           string     `json:"link"`
	Keywords       []string   `json:"keywords"`
	Complexity     Complexity `json:"complexity"`
	Whiteboard     string     `json:"whiteboard"`
	DryRun         string     `json:"dry_run"`
	TestCases      string     `json:"test_cases"`
	IaSolution     string     `json:"ia_solution"`
	ManualSolution string     `json:"manual_solution"`
}

func (p Problem) FilterValue() string { return p.ProblemTitle }
func (p Problem) Title() string       { return p.ProblemTitle }
func (p Problem) Description() string { return strings.Join(p.Keywords, ", ") }

type Complexity struct {
	Time  ComplexityDetail `json:"time"`
	Space ComplexityDetail `json:"space"`
}
type ComplexityDetail struct {
	Notation      string `json:"notation"`
	Justification string `json:"justification"`
}

// *========( MAIN APP MODEL )=========*
type viewState int
type activePane int

const (
	searchView viewState = iota
	problemView
)
const (
	leftPane activePane = iota
	rightPane
)

type model struct {
	allProblems       []Problem
	list              list.Model
	activeProblem     *Problem
	leftViewport      viewport.Model
	rightViewport     viewport.Model
	glamour           *glamour.TermRenderer
	state             viewState
	activePane        activePane
	solutionViewIndex int
	ready             bool
	quitting          bool
	windowWidth       int
	windowHeight      int
}

func newModel(problems []Problem) *model {
	glamourRenderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(0))
	items := make([]list.Item, len(problems))
	for i, p := range problems {
		items[i] = p
	}
	problemList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	problemList.Title = ""
	problemList.Styles.Title = lipgloss.NewStyle()
	problemList.SetShowHelp(false)
	problemList.DisableQuitKeybindings()
	problemList.FilterInput.Prompt = "Search: "

	m := model{
		allProblems: problems,
		list:        problemList,
		state:       searchView,
		activePane:  leftPane,
		glamour:     glamourRenderer,
	}

	return &m
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.onWindowSizeChanged()
		m.ready = true
	}
	var cmd tea.Cmd
	switch m.state {
	case searchView:
		cmd = m.updateSearchView(msg)
	case problemView:
		cmd = m.updateProblemView(msg)
	}
	return m, cmd
}

func (m *model) updateSearchView(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Store the previous filter value to detect changes
	prevFilter := m.list.FilterValue()

	// Let the list component handle its own state updates first.
	m.list, cmd = m.list.Update(msg)

	query := m.list.FilterValue()

	// Check if the filter changed and starts with #
	if strings.HasPrefix(query, "#") && query != prevFilter {
		m.filterProblems(query)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+j":
			m.list.CursorDown()
			return nil
		case "ctrl+k":
			m.list.CursorUp()
			return nil
		case "enter":
			if selectedItem, ok := m.list.SelectedItem().(Problem); ok {
				m.activeProblem = &selectedItem
				m.state = problemView
				m.setupViewports()
			}
		}
	}
	return cmd
}

func (m *model) filterProblems(query string) {
	filteredItems := make([]list.Item, 0)
	searchStr := strings.TrimSpace(strings.TrimLeft(query, "#"))

	if searchStr == "" {
		allItems := make([]list.Item, len(m.allProblems))
		for i, p := range m.allProblems {
			allItems[i] = p
		}
		m.list.SetItems(allItems)
		m.list.Select(0)
		m.list.Update(tea.KeyMsg{})
		return
	}

	separators := strings.NewReplacer(",", " ", ";", " ", "  ", " ")
	searchStr = separators.Replace(searchStr)
	searchKeywords := strings.Fields(strings.ToLower(searchStr))

	for _, prob := range m.allProblems {
		problemKeywords := make(map[string]struct{})
		for _, kw := range prob.Keywords {
			problemKeywords[strings.ToLower(strings.TrimLeft(kw, "#"))] = struct{}{}
		}

		matchesAll := true
		for _, searchKw := range searchKeywords {
			if _, found := problemKeywords[searchKw]; !found {
				matchesAll = false
				break
			}
		}

		if matchesAll {
			filteredItems = append(filteredItems, prob)
		}
	}

	m.list.SetItems(filteredItems)
	if len(filteredItems) > 0 {
		m.list.Select(0)
		// Force the list to refresh its visual state
		m.list.Update(tea.KeyMsg{})
	}
}

func (m *model) updateProblemView(msg tea.Msg) tea.Cmd {
	var cmd, cmd2 tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "/", "esc":
			m.list.ResetFilter()
			allItems := make([]list.Item, len(m.allProblems))
			for i, p := range m.allProblems {
				allItems[i] = p
			}
			m.list.SetItems(allItems)
			m.state = searchView
			return nil
		case "tab":
			m.activePane = (m.activePane + 1) % 2
		case "c":
			m.solutionViewIndex = (m.solutionViewIndex + 1) % 2

			m.setupViewports()

		}
	}

	// Pass the message to both viewports to keep them in sync
	m.leftViewport, cmd = m.leftViewport.Update(msg)
	m.rightViewport, cmd2 = m.rightViewport.Update(msg)

	return tea.Batch(cmd, cmd2)
}
func (m *model) onWindowSizeChanged() {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMargin := headerHeight + footerHeight
	if !m.ready {
		m.leftViewport = viewport.New(0, 0)
		m.rightViewport = viewport.New(0, 0)
	}
	listHeight := m.windowHeight - verticalMargin
	m.list.SetSize(m.windowWidth-docStyle.GetHorizontalFrameSize(), listHeight)
	paneWidth := m.windowWidth/2 - activePaneStyle.GetHorizontalFrameSize()
	paneHeight := m.windowHeight - verticalMargin - activePaneStyle.GetVerticalFrameSize()
	m.leftViewport.Width = paneWidth
	m.leftViewport.Height = paneHeight
	m.rightViewport.Width = paneWidth
	m.rightViewport.Height = paneHeight
	if m.state == problemView && m.activeProblem != nil {
		m.setupViewports()
	}
}
func (m *model) setupViewports() {
	paneWidth := m.leftViewport.Width - 4

	wbTitle := paneTitleStyleFirst.Render(strings.ToUpper("Whiteboard"))
	wbContent := m.activeProblem.Whiteboard
	drTitle := paneTitleStyle.Render(strings.ToUpper("Dry Run"))
	drContent := m.activeProblem.DryRun
	ttTitle := paneTitleStyle.Render(strings.ToUpper("Test Cases"))
	ttContent := m.activeProblem.TestCases
	tcTitle := paneTitleStyle.Render(strings.ToUpper("Time Complexity"))
	tcContent := fmt.Sprintf("%s: %s", m.activeProblem.Complexity.Time.Notation, m.activeProblem.Complexity.Time.Justification)
	scTitle := paneTitleStyle.Render(strings.ToUpper("Space Complexity"))
	scContent := fmt.Sprintf("%s: %s", m.activeProblem.Complexity.Space.Notation, m.activeProblem.Complexity.Space.Justification)
	leftJoined := lipgloss.JoinVertical(lipgloss.Left,
		wbTitle, wbContent,
		drTitle, drContent,
		ttTitle, ttContent,
		tcTitle, tcContent,
		scTitle, scContent,
	)
	leftFinal := lipgloss.NewStyle().Width(paneWidth).Render(leftJoined)
	m.leftViewport.SetContent(leftFinal)
	m.leftViewport.GotoTop()

	var solutionTitle, solutionCode string
	if m.solutionViewIndex == 0 {
		solutionTitle = "Manual Solution"
		solutionCode = m.activeProblem.ManualSolution
	} else {
		solutionTitle = "IA Solution"
		solutionCode = m.activeProblem.IaSolution
	}

	csTitle := paneTitleStyleFirst.Render(strings.ToUpper(solutionTitle))
	commentedSolutionCode, _ := m.glamour.Render(fmt.Sprintf("```python\n%s\n```", solutionCode))
	rightJoined := lipgloss.JoinVertical(lipgloss.Left, csTitle, commentedSolutionCode)
	rightFinal := lipgloss.NewStyle().Width(paneWidth).Render(rightJoined)
	m.rightViewport.SetContent(rightFinal)
	m.rightViewport.GotoTop()
}
func (m *model) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	if !m.ready {
		return "Initializing..."
	}
	var mainContent string
	if m.state == searchView {
		mainContent = docStyle.Render(m.list.View())
	} else {
		var left, right lipgloss.Style
		if m.activePane == leftPane {
			left, right = activePaneStyle, inactivePaneStyle
		} else {
			left, right = inactivePaneStyle, activePaneStyle
		}
		left.Width(m.windowWidth / 2)
		right.Width(m.windowWidth - m.windowWidth/2)
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, left.Render(m.leftViewport.View()), right.Render(m.rightViewport.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.headerView(), mainContent, m.footerView())
}
func (m model) headerView() string {
	titleStr := "LeetCode Finder"
	if m.state == problemView && m.activeProblem != nil {
		titleStr = m.activeProblem.ProblemTitle
	}
	title := titleStyle.Render(titleStr)
	info := infoStyle.Render(m.infoView())
	line := strings.Repeat("─", max(0, m.windowWidth-lipgloss.Width(title)-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line, info)
}
func (m model) infoView() string {
	if m.state == problemView && m.activeProblem != nil {
		return fmt.Sprintf("%s Time | %s Space", m.activeProblem.Complexity.Time.Notation, m.activeProblem.Complexity.Space.Notation)
	}
	return "v4.3"
}

func (m model) footerView() string {
	var shortcuts []string
	if m.list.FilterState() == list.Filtering {
		shortcuts = []string{
			keyStyle.Render("ctrl+j/k") + " navigate",
			keyStyle.Render("enter") + " select",
			keyStyle.Render("esc") + " cancel",
		}
	} else if m.state == searchView {
		shortcuts = []string{
			keyStyle.Render("/") + " search",
			keyStyle.Render("enter") + " select",
			keyStyle.Render("ctrl+c") + " quit",
		}
	} else { // problemView
		shortcuts = []string{
			keyStyle.Render("tab") + " switch pane",
			keyStyle.Render("c") + " toggle solution",
			keyStyle.Render("/") + " search",
			keyStyle.Render("esc") + " back",
		}
	}
	return footerStyle.Width(m.windowWidth).Render(strings.Join(shortcuts, " "))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func main() {
	jsonFile, err := os.Open("problems.json")
	if err != nil {
		log.Fatalf("Error opening problems.json: %v", err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var problems []Problem
	if err := json.Unmarshal(byteValue, &problems); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	p := tea.NewProgram(newModel(problems), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
