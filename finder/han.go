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
			s := strings.TrimSpace(sb.String())
			if s != "" {
				ret = append(ret, s)
			}
			sb.Reset()
		} else {
			sb.WriteRune(ch)
		}
	}
	if sb.Len() > 0 {
		ret = append(ret, sb.String())
	}
	return ret
}
