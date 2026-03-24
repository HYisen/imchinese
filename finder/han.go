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

func isDotInsideWord(i int, s string) bool {
	if s[i] != '.' {
		return false
	}
	if i+1 == len(s) {
		return false
	}
	// dot as end of sentence
	if i+2 < len(s) && s[i+1] == ' ' && unicode.IsUpper(rune(s[i+2])) {
		return false
	}
	return true
}

// FilterWord reads line as a sentence, find word that are candidates.
// Words with kanji or punctuations are typically not candidates, while initialisms are.
func FilterWord(text Text) []Candidate {
	var candidates []Candidate
	tables := []*unicode.RangeTable{unicode.Han, unicode.Punct}
	var sb strings.Builder
	for i, ch := range text.Item {
		if unicode.IsOneOf(tables, ch) && !isDotInsideWord(i, text.Item) {
			candidates = trimAndAppendNonEmpty(candidates, sb.String(), text)
			sb.Reset()
		} else {
			sb.WriteRune(ch)
		}
	}
	candidates = trimAndAppendNonEmpty(candidates, sb.String(), text)
	return DropNumber(candidates)
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

func DropNumber(items []Candidate) []Candidate {
	var ret []Candidate
	for _, item := range items {
		if !isNumber(item.Word) {
			ret = append(ret, item)
		}
	}
	return ret
}

func isNumber(s string) bool {
	// not [strconv.Atoi] to allow decimal separator
	for _, ch := range s {
		if !unicode.IsDigit(ch) && ch != '.' && ch != '-' {
			return false
		}
	}
	return true
}
