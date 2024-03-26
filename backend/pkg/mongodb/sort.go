package mongodb

type SortOrder = int
type Sort = map[string]SortOrder

const (
	Ascending  SortOrder = 1
	Descending SortOrder = -1
)

// MakeSort base on sorted key and order
func MakeSort(sortedBy string, order SortOrder) Sort {
	sort := make(Sort)
	sort[sortedBy] = order
	return sort
}
