package varg

import (
	"fmt"
	"strings"
	"testing"
)

func (c *candidate) render() string {
	var render []string

	for _, line := range c.paragraph {
		render = append(render,
			strings.Join(line.words, " "))
	}

	return strings.Join(render, "\n")
}

func TestGreedyFill(t *testing.T) {
	raw := `read buffer names to search, one by one, ended with RET.  With a prefix argument, they ask for a regexp, and search in buffers whose names match the specified regexp.  Interactively 'multi-isearch-files' and 'multi-isearch-files-regexp' read file names to search, one by one, ended with RET.  With a prefix argument, they ask for a wildcard, and search in file buffers whose file names match the specified wildcard.`
	words := strings.Split(raw, " ")

	maxWidth := 72
	c := greedyFill(maxWidth, words)

	result := c.render()
	expected := `read buffer names to search, one by one, ended with RET.  With a prefix
argument, they ask for a regexp, and search in buffers whose names match
the specified regexp.  Interactively 'multi-isearch-files' and
'multi-isearch-files-regexp' read file names to search, one by one,
ended with RET.  With a prefix argument, they ask for a wildcard, and
search in file buffers whose file names match the specified wildcard.`

	if result != expected {
		t.FailNow()
	}
}

func TestPrettyFill(t *testing.T) {
	raw := `read buffer names to search, one by one, ended with RET.  With a prefix argument, they ask for a regexp, and search in buffers whose names match the specified regexp.  Interactively 'multi-isearch-files' and 'multi-isearch-files-regexp' read file names to search, one by one, ended with RET.  With a prefix argument, they ask for a wildcard, and search in file buffers whose file names match the specified wildcard.`
	words := strings.Split(raw, " ")

	maxWidth := 72
	c := prettyFill(maxWidth, words)
	result := c.render()
	expected := `read buffer names to search, one by one, ended with RET.  With a prefix
argument, they ask for a regexp, and search in buffers whose names
match the specified regexp.  Interactively 'multi-isearch-files' and
'multi-isearch-files-regexp' read file names to search, one by one,
ended with RET.  With a prefix argument, they ask for a wildcard, and
search in file buffers whose file names match the specified wildcard.`

	if result != expected {
		t.FailNow()
	}
}

func TestDegenerate(t *testing.T) {
	raw := `aaaa`
	words := strings.Split(raw, " ")
	maxWidth := 8

	p := prettyFill(maxWidth, words)
	g := greedyFill(maxWidth, words)

	if p.render() != raw {
		t.FailNow()
	}

	if g.render() != raw {
		t.FailNow()
	}
}

func TestFillSmall(t *testing.T) {
	raw := `aaa bb cc ddddd`
	words := strings.Split(raw, " ")

	maxWidth := 6
	p := prettyFill(maxWidth, words)
	g := greedyFill(maxWidth, words)

	pExpected := `aaa
bb cc
ddddd`

	gExpected := `aaa bb
cc
ddddd`

	if pExpected != p.render() {
		t.FailNow()
	}

	if gExpected != g.render() {
		t.FailNow()
	}
}

func newCandidateWithChecker(t *testing.T) (
	c *candidate, check func(expected string)) {
	c = newCandidate()

	check = func(expected string) {
		result := fmt.Sprintf("%v", c)

		if result != expected {
			t.Fatalf("Expected %s, but got %s", expected, result)
		}
	}

	return c, check
}

func TestAddWord(t *testing.T) {
	c, check := newCandidateWithChecker(t)

	check("&{0 [{[] 0}]}")
	c.addWord(4, "a")
	check("&{0 [{[a] 1}]}")
	c.addWord(4, "b")
	check("&{0 [{[a b] 3}]}")

	// Overrun
	c.addWord(4, "cc")
	check("&{2147483647 [{[a b cc] 6}]}")
}

func TestBreakAndAdd(t *testing.T) {
	c, check := newCandidateWithChecker(t)

	check("&{0 [{[] 0}]}")
	c.breakAndAdd(4, "a")
	check("&{16 [{[] 0} {[a] 1}]}")
	c.breakAndAdd(4, "b")
	check("&{25 [{[] 0} {[a] 1} {[b] 1}]}")
	c.breakAndAdd(4, "cc")
	check("&{34 [{[] 0} {[a] 1} {[b] 1} {[cc] 2}]}")
	c.breakAndAdd(4, "dddd")
	check("&{38 [{[] 0} {[a] 1} {[b] 1} {[cc] 2} {[dddd] 4}]}")

	// Overrun
	c.breakAndAdd(4, "eeeee")
	check("&{2147483647 [{[] 0} {[a] 1} {[b] 1} {[cc] 2} {[dddd] 4} {[eeeee] 5}]}")
}
