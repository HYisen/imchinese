package finder

import (
	"strings"
	"unicode"
)

type Candidate struct {
	Word string
	Line string
	Path string
}

// FilterWord reads line as a sentence, find word that are candidates.
// Words with kanji or punctuations are typically not candidates, while initialisms are.
func FilterWord(text Text) []Candidate {
	var ret []Candidate
	tables := []*unicode.RangeTable{unicode.Han, unicode.Punct}
	var sb strings.Builder
	for _, ch := range text.Item {
		if unicode.IsOneOf(tables, ch) {
			ret = trimAndAppendNonEmpty(ret, sb.String(), text)
			sb.Reset()
		} else {
			sb.WriteRune(ch)
		}
	}
	ret = trimAndAppendNonEmpty(ret, sb.String(), text)
	return ret
}

func trimAndAppendNonEmpty(slice []Candidate, s string, text Text) []Candidate {
	s = strings.TrimSpace(s)
	if s != "" {
		slice = append(slice, Candidate{
			Word: s,
			Line: text.Item,
			Path: text.Path,
		})
	}
	return slice
}
