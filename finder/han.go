package finder

import (
	"strconv"
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
	// dot between numbers
	if (i-1 >= 0 && unicode.IsNumber(rune(s[i-1]))) || (i+1 < len(s) && unicode.IsNumber(rune(s[i+1]))) {
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
	// Considering the punctuation separating, only natual number would be possible here.
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
		if _, err := strconv.Atoi(item.Word); err == nil {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}
