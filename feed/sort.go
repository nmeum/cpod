package feed

type ByDate []Item

func (b ByDate) Len() int {
	return len(b)
}

func (b ByDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByDate) Less(i, j int) bool {
	return b[i].Date.After(b[j].Date)
}
