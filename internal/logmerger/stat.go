package logmerger

type statistics struct {
	Total   int
	PerFile []perFileStatistics
}

type perFileStatistics struct {
	FileName string
	Count    int
	IsMain   bool
}

type mainThenFn []perFileStatistics

func (a mainThenFn) Len() int      { return len(a) }
func (a mainThenFn) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a mainThenFn) Less(i, j int) bool {
	if (a[i].IsMain || a[j].IsMain) && !(a[i].IsMain && a[j].IsMain) {
		return a[i].IsMain
	} else {
		return a[i].FileName < a[j].FileName
	}

}
