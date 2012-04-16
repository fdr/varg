// To assist with filling paragraphs of text.
package varg

type candidate struct {
	badness   int
	paragraph []line
}

type line struct {
	words []string
	width int
}

func (c *candidate) addWord(w string) {
	l := &c.paragraph[len(c.paragraph)-1]
	l.words = append(l.words, w)
	l.width += len(w)
}

func (c *candidate) breakLine(targetWidth int) {
	delta := targetWidth - c.paragraph[len(c.paragraph)-1].width
	c.badness += delta * delta
	c.paragraph = append(c.paragraph, line{})
}

func newCandidate() *candidate {
	var c candidate

	c.paragraph = make([]line, 1, 30)

	return &c
}

// Fills a paragraph using a greedy algorithm.
//
// This maximizees the length of each line from top to bottom.
// Badness scores (as ascertained by prettyFill) are computed so that
// this procedure can be used to cull candidate line breaks.
func greedyFill(width int, words []string) *candidate {
	c := newCandidate()

	for _, w := range words {
		if c.paragraph[len(c.paragraph)-1].width+len(w) > width {
			c.breakLine(width)
			c.addWord(w)
		} else {
			c.addWord(w)
		}
	}

	return c
}

// Fills a paragraph with minimum raggedness.
//
// This searches for a minimal cost as determined by (targetWidth -
// lineWidth) ** 2.  Candidate solutions have more cost than the
// greedy fill are eliminated.
func prettyFill(width int, words []string) *candidate {
	badnessMax := greedyFill(width, words).badness
	candidates := make([]candidate, 0, len(words))
	candidates = append(candidates, *newCandidate())

	for _, w := range words {
		for _, appendC := range candidates {
			// Skip if badness has already been exceeded
			if appendC.badness > badnessMax {
				continue
			}

			lineBreakC := appendC

			lineBreakC.breakLine(width)
			lineBreakC.addWord(w)
		}
	}

	return &candidates[0]
}
