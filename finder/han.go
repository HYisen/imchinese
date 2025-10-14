package finder

import (
	"strings"
	"unicode"
)

// FilterWord reads line as a sentence, find word that are candidates.
// Words with kanji or punctuations are typically not candidates, while initialisms are.
func FilterWord(line string) []string {
	var ret []string
	tables := []*unicode.RangeTable{unicode.Han, unicode.Punct}
	var sb strings.Builder
	for _, ch := range line {
		if unicode.IsOneOf(tables, ch) {
			ret = trimAndAppendNonEmpty(ret, sb.String())
			sb.Reset()
		} else {
			sb.WriteRune(ch)
		}
	}
	ret = trimAndAppendNonEmpty(ret, sb.String())
	return ret
}

func trimAndAppendNonEmpty(slice []string, s string) []string {
	s = strings.TrimSpace(s)
	if s != "" {
		slice = append(slice, s)
	}
	return slice
}
