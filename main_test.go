package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

func titles(ps []Problem) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.ProblemTitle
	}
	return out
}

func listTitles(l list.Model) []string {
	items := l.Items()
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.(Problem).ProblemTitle
	}
	return out
}

func TestLoadProblems(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantTitles []string
		wantErr    string
	}{
		{"valid file", "testdata/valid.json", []string{"Two Sum", "Valid Parentheses"}, ""},
		{"empty array", "testdata/empty.json", []string{}, ""},
		{"missing fields default to zero values", "testdata/partial.json", []string{"Minimal"}, ""},
		{"missing file", "testdata/does-not-exist.json", nil, "no such file"},
		{"malformed json", "testdata/malformed.json", nil, "unmarshalling JSON"},
		{"wrong top-level shape", "testdata/wrong_shape.json", nil, "unmarshalling JSON"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadProblems(tt.path)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(titles(got), tt.wantTitles) {
				t.Errorf("titles = %v, want %v", titles(got), tt.wantTitles)
			}
		})
	}
}

func TestLoadProblemsParsesAllFields(t *testing.T) {
	got, err := loadProblems("testdata/valid.json")
	if err != nil {
		t.Fatal(err)
	}
	want := Problem{
		ProblemTitle: "Two Sum",
		Link:         "https://leetcode.com/problems/two-sum/",
		Keywords:     []string{"#Array", "#HashTable"},
		Complexity: Complexity{
			Time:  ComplexityDetail{Notation: "O(N)", Justification: "One pass."},
			Space: ComplexityDetail{Notation: "O(N)", Justification: "Hash map."},
		},
		Whiteboard:     "Use a hash map.",
		DryRun:         "nums = [2, 7]",
		TestCases:      "[2,7], 9",
		IaSolution:     "def f(): pass",
		ManualSolution: "def g(): pass",
	}
	if !reflect.DeepEqual(got[0], want) {
		t.Errorf("first problem = %+v, want %+v", got[0], want)
	}
}

// The bundled data file must stay loadable and well formed.
func TestBundledProblemsJSON(t *testing.T) {
	got, err := loadProblems("problems.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("problems.json is empty")
	}
	seen := map[string]bool{}
	for i, p := range got {
		if p.ProblemTitle == "" {
			t.Errorf("problem %d has no title", i)
		}
		if seen[p.ProblemTitle] {
			t.Errorf("duplicate title %q", p.ProblemTitle)
		}
		seen[p.ProblemTitle] = true
		if !strings.HasPrefix(p.Link, "https://leetcode.com/problems/") {
			t.Errorf("%q: unexpected link %q", p.ProblemTitle, p.Link)
		}
		if len(p.Keywords) == 0 {
			t.Errorf("%q: no keywords", p.ProblemTitle)
		}
	}
}

func TestProblemListItem(t *testing.T) {
	tests := []struct {
		name     string
		problem  Problem
		wantDesc string
	}{
		{"several keywords", Problem{ProblemTitle: "Two Sum", Keywords: []string{"#Array", "#HashTable"}}, "#Array, #HashTable"},
		{"one keyword", Problem{ProblemTitle: "X", Keywords: []string{"#Stack"}}, "#Stack"},
		{"no keywords", Problem{ProblemTitle: "X"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.problem.Description(); got != tt.wantDesc {
				t.Errorf("Description() = %q, want %q", got, tt.wantDesc)
			}
			if tt.problem.Title() != tt.problem.ProblemTitle || tt.problem.FilterValue() != tt.problem.ProblemTitle {
				t.Errorf("Title/FilterValue should both return the problem title")
			}
		})
	}
}

func TestParseSearchKeywords(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"array", []string{"array"}},
		{"Array HashTable", []string{"array", "hashtable"}},
		{"array,hashtable", []string{"array", "hashtable"}},
		{"array; hashtable", []string{"array", "hashtable"}},
		{"  array ,;  stack  ", []string{"array", "stack"}},
		{"", []string{}},
		{" , ; ", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := parseSearchKeywords(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSearchKeywords(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestHasAllKeywords(t *testing.T) {
	p := Problem{Keywords: []string{"#Array", "#HashTable", "string"}}
	tests := []struct {
		name string
		kws  []string
		want bool
	}{
		{"no keywords matches everything", nil, true},
		{"single match", []string{"array"}, true},
		{"match without hash prefix on problem side", []string{"string"}, true},
		{"all match", []string{"array", "hashtable"}, true},
		{"one missing", []string{"array", "stack"}, false},
		{"partial word does not match", []string{"hash"}, false},
		{"query keyword is expected lowercase", []string{"Array"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAllKeywords(p, tt.kws); got != tt.want {
				t.Errorf("hasAllKeywords(%v) = %v, want %v", tt.kws, got, tt.want)
			}
		})
	}
	if hasAllKeywords(Problem{}, []string{"array"}) {
		t.Error("problem without keywords must not match a keyword")
	}
}

func TestFilterProblems(t *testing.T) {
	problems, err := loadProblems("testdata/search.json")
	if err != nil {
		t.Fatal(err)
	}
	all := titles(problems)

	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"single keyword", "#array", []string{"Two Sum", "Group Anagrams", "Binary Search"}},
		{"case insensitive", "#ARRAY", []string{"Two Sum", "Group Anagrams", "Binary Search"}},
		{"keywords are ANDed", "#array hashtable", []string{"Two Sum", "Group Anagrams"}},
		{"comma separated", "#array,string", []string{"Group Anagrams"}},
		{"semicolon separated", "#string; stack", []string{"Valid Parentheses"}},
		{"extra # and spaces trimmed", "##  stack ", []string{"Valid Parentheses"}},
		{"matches keyword stored without #", "#binarysearch", []string{"Binary Search"}},
		{"no match", "#graph", []string{}},
		{"unmatched alongside matched", "#array graph", []string{}},
		{"empty query restores all", "#", all},
		{"whitespace-only query restores all", "#   ", all},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := list.New(nil, list.NewDefaultDelegate(), 80, 20)
			m := &model{allProblems: problems, list: l}
			m.filterProblems(tt.query)
			if got := listTitles(m.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterProblems(%q) = %v, want %v", tt.query, got, tt.want)
			}
			if len(m.allProblems) != len(problems) {
				t.Error("filterProblems must not modify allProblems")
			}
		})
	}
}

func TestFilterProblemsResetsSelection(t *testing.T) {
	problems, err := loadProblems("testdata/search.json")
	if err != nil {
		t.Fatal(err)
	}
	items := make([]list.Item, len(problems))
	for i, p := range problems {
		items[i] = p
	}
	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
	l.Select(3)
	m := &model{allProblems: problems, list: l}
	m.filterProblems("#array")
	if m.list.Index() != 0 {
		t.Errorf("cursor index = %d, want 0 after filtering", m.list.Index())
	}
}

func TestMax(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{1, 2, 2}, {2, 1, 2}, {3, 3, 3}, {-1, 0, 0}, {-5, -2, -2},
	}
	for _, tt := range tests {
		if got := max(tt.a, tt.b); got != tt.want {
			t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestInfoView(t *testing.T) {
	p := &Problem{Complexity: Complexity{
		Time:  ComplexityDetail{Notation: "O(N)"},
		Space: ComplexityDetail{Notation: "O(1)"},
	}}
	tests := []struct {
		name string
		m    model
		want string
	}{
		{"search view shows version", model{state: searchView}, "v4.3"},
		{"problem view shows complexity", model{state: problemView, activeProblem: p}, "O(N) Time | O(1) Space"},
		{"problem view without problem falls back", model{state: problemView}, "v4.3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.infoView(); got != tt.want {
				t.Errorf("infoView() = %q, want %q", got, tt.want)
			}
		})
	}
}
