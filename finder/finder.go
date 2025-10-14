package finder

// Find reads passage as markdown source code, returns filtered.
func Find(passage string) []Candidate {
	var ret []Candidate
	for _, text := range FilterText(passage) {
		ret = append(ret, FilterWord(text)...)
	}
	return ret
}
