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

func (c *candidate) copy() *candidate {
	var cpy candidate

	cpy.badness = c.badness
	cpy.paragraph = make([]line, len(c.paragraph))

	for i, _ := range c.paragraph {
		l := &cpy.paragraph[i]
		l.width = c.paragraph[i].width
		l.words = make([]string, len(c.paragraph[i].words))
		copy(l.words, c.paragraph[i].words)
	}

	return &cpy
}

func (c *candidate) addWord(targetWidth int, w string) {
	l := &c.paragraph[len(c.paragraph)-1]

	l.width += len(w)

	if len(l.words) != 0 {
		// On all words except first in the line, count the
		// space separation between words against the width of
		// the line.
		l.width += 1
	}

	l.words = append(l.words, w)

	// The last line doesn't count towards badness when it falls
	// into the target line length, but should it exceed that
	// length, it is not to be considered an acceptable solution.
	if l.width > targetWidth {
		c.badness = 2147483647
	}
}

func (c *candidate) breakAndAdd(targetWidth int, w string) {
	delta := targetWidth - c.paragraph[len(c.paragraph)-1].width

	c.badness += delta * delta
	c.paragraph = append(c.paragraph, *newLine())

	c.addWord(targetWidth, w)
}

func newLine() *line {
	var l line

	l.words = make([]string, 0)

	return &l
}

func newCandidate() *candidate {
	var c candidate

	c.paragraph = make([]line, 1, 30)
	c.paragraph[0] = *newLine()

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
		if c.paragraph[len(c.paragraph)-1].width+len(w) >= width {
			c.breakAndAdd(width, w)
		} else {
			c.addWord(width, w)
		}
	}

	return c
}

// Fills a paragraph with minimum raggedness.
//
// This searches for a minimal cost as determined by (targetWidth -
// lineWidth) ** 2.  Candidate solutions have more cost than the
// greedy fill are eliminated.
func prettyFill(width int, words []string) (best *candidate) {
	badnessMax := greedyFill(width, words).badness
	candidates := make([]candidate, 0, len(words))
	candidates = append(candidates, *newCandidate())

	for _, w := range words {
		for i, _ := range candidates {
			appendC := &candidates[i]

			// Skip if badness has already been exceeded
			if appendC.badness > badnessMax {
				continue
			}

			lineBreakC := appendC.copy()
			appendC.addWord(width, w)
			lineBreakC.breakAndAdd(width, w)

			candidates = append(candidates, *lineBreakC)
		}
	}

	// Locate minimum badness and return it
	minBadness := badnessMax

	for i, _ := range candidates {
		c := &candidates[i]

		if c.badness <= minBadness {
			best = c
			minBadness = c.badness
		}
	}

	return best
}
