package finder

// Find reads passage as markdown source code, returns filtered word candidates.
func Find(passage string) []string {
	lines := FilterText(passage)

	var words []string
	for _, line := range lines {
		words = append(words, FilterWord(line)...)
	}
	return words
}
