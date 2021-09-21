package hw03frequencyanalysis

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

// RegExp [0] - match whole string.
// RegExp [1] - match word without punctuation.
// RegExp [2] - match punctuation.
var (
	RegExp          = regexp.MustCompile(`^([^.|,|;|:|?|!]+)([.|,|;|:|?|!]+)$`)
	ErrUnsortedText = errors.New("included text fragments are not sorted yet")
)

type Chunk struct {
	count    uint
	fragment string
}

type Chunks struct {
	chunks   []Chunk
	raw      map[string]uint
	isSorted bool
}

func (c *Chunks) Initialize(text string) *Chunks {
	c.raw = make(map[string]uint)
	// trim dashes
	text = strings.ReplaceAll(text, "-", "")

	for _, fragment := range strings.Fields(text) {
		fragment := strings.ToLower(fragment)
		if fragments := RegExp.FindStringSubmatch(fragment); len(fragments) == 3 {
			fragment = fragments[1]
		}

		if _, existed := c.raw[fragment]; !existed {
			c.raw[fragment] = 1
			continue
		}
		c.raw[fragment]++
	}

	for fragment, count := range c.raw {
		c.chunks = append(c.chunks, Chunk{count: count, fragment: fragment})
	}

	return c
}

func (c *Chunks) Sort() *Chunks {
	sort.Slice(c.chunks, func(i, j int) bool {
		if c.chunks[i].count != c.chunks[j].count {
			return c.chunks[i].count > c.chunks[j].count
		}
		return strings.Compare(c.chunks[i].fragment, c.chunks[j].fragment) < 0
	})
	c.isSorted = true

	return c
}

func (c *Chunks) GetTop10() ([]string, error) {
	result := make([]string, 0, 10)

	if !c.isSorted {
		return result, ErrUnsortedText
	}

	for i, chunk := range c.chunks {
		if i >= 10 {
			break
		}
		result = append(result, chunk.fragment)
	}

	return result, nil
}

func Top10(text string) []string {
	c := Chunks{}
	result, _ := c.Initialize(text).Sort().GetTop10()

	return result
}
