package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

// FuzzSearch drives the "#keyword" filter with arbitrary query strings.
// Properties: filterProblems never panics, never changes allProblems, and
// every problem it keeps carries every queried keyword (compared
// case-insensitively, ignoring the leading "#" on the problem's keywords).
// An empty query (after trimming "#" and spaces) restores the full list.
func FuzzSearch(f *testing.F) {
	problems, err := loadProblems("testdata/search.json")
	if err != nil {
		f.Fatal(err)
	}

	for _, seed := range []string{
		"#array", "#ARRAY", "#array hashtable", "#array,string", "#string; stack",
		"##  stack ", "#binarysearch", "#graph", "#", "#   ", "", "array",
		"#\x00", "#\xff\xfe", "#ünïcödé", "#" + strings.Repeat("a", 4096),
	} {
		f.Add(seed)
	}
	// Seed with keywords taken from the real, bundled problem set.
	if bundled, err := loadProblems("problems.json"); err == nil {
		for _, p := range bundled[:min(len(bundled), 25)] {
			f.Add("#" + strings.Join(p.Keywords, " "))
		}
	}

	f.Fuzz(func(t *testing.T, query string) {
		l := list.New(nil, list.NewDefaultDelegate(), 80, 20)
		m := &model{allProblems: problems, list: l}
		m.filterProblems(query)

		if len(m.allProblems) != len(problems) {
			t.Fatalf("allProblems length changed for query %q", query)
		}

		searchStr := strings.TrimSpace(strings.TrimLeft(query, "#"))
		items := m.list.Items()
		if searchStr == "" {
			if len(items) != len(problems) {
				t.Fatalf("empty query %q returned %d of %d problems", query, len(items), len(problems))
			}
			return
		}

		want := parseSearchKeywords(searchStr)
		for _, it := range items {
			p, ok := it.(Problem)
			if !ok {
				t.Fatalf("query %q produced a non-Problem item %T", query, it)
			}
			have := map[string]bool{}
			for _, kw := range p.Keywords {
				have[strings.ToLower(strings.TrimLeft(kw, "#"))] = true
			}
			for _, w := range want {
				if !have[w] {
					t.Fatalf("query %q kept %q which lacks keyword %q", query, p.ProblemTitle, w)
				}
			}
		}
	})
}
