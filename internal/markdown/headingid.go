package markdown

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

// newGFMIDs returns a parser.IDs that generates heading IDs following the
// GitHub Flavored Markdown convention: Unicode letters and digits are kept
// (lowercased where applicable), spaces/hyphens become "-", others are
// dropped. This preserves non-ASCII characters such as Japanese kana/kanji,
// unlike goldmark's default which strips all multi-byte characters.
func newGFMIDs() parser.IDs {
	return &gfmIDs{values: map[string]bool{}, next: map[string]int{}}
}

type gfmIDs struct {
	values map[string]bool
	// next remembers the last suffix handed out for each base ID so that a
	// document with many identical headings (e.g. a conversation log with
	// hundreds of "User" / "Claude" headings) does not rescan "-1", "-2", ...
	// from the start every time.
	next map[string]int
}

func (s *gfmIDs) Generate(value []byte, kind ast.NodeKind) []byte {
	value = bytes.TrimSpace(value)
	var result []rune
	for _, r := range string(value) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			result = append(result, unicode.ToLower(r))
		case unicode.IsSpace(r) || r == '-' || r == '_':
			result = append(result, '-')
		}
	}
	res := []byte(string(result))
	if len(res) == 0 {
		if kind == ast.KindHeading {
			res = []byte("heading")
		} else {
			res = []byte("id")
		}
	}
	key := string(res)
	if !s.values[key] {
		s.values[key] = true
		return res
	}
	// Suffixes below next[key] are all taken (IDs are never released), so
	// starting there yields the same result as scanning from 1.
	for i := s.next[key] + 1; ; i++ {
		newKey := fmt.Sprintf("%s-%d", key, i)
		if !s.values[newKey] {
			s.values[newKey] = true
			s.next[key] = i
			return []byte(newKey)
		}
	}
}

func (s *gfmIDs) Put(value []byte) {
	s.values[string(value)] = true
}
