package varg

import (
	"strings"
	"testing"
)

func TestGreedyFill(t *testing.T) {
	raw := `read buffer names to search, one by one, ended with RET.  With a prefix argument, they ask for a regexp, and search in buffers whose names match the specified regexp.  Interactively 'multi-isearch-files' and 'multi-isearch-files-regexp' read file names to search, one by one, ended with RET.  With a prefix argument, they ask for a wildcard, and search in file buffers whose file names match the specified wildcard.`
	words := strings.Split(raw, " ")

	maxWidth := 72
	c := greedyFill(maxWidth, words)

	for _, line := range c.paragraph {
		if line.width > maxWidth {
			t.FailNow()
		}
	}
}
