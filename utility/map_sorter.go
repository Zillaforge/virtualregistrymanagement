package utility

import "sort"

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int { return len(p) }

func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SortValues(m map[string]int) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}
