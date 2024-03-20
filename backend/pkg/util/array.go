package util

// Combination combine for a1 and a2
func Combination[T interface{}](a1 []T, a2 []T) []T {
	combine := make([]T, len(a1)+len(a2))
	for _, ele := range a1 {
		combine = append(combine, ele)
	}
	for _, ele := range a2 {
		combine = append(combine, ele)
	}
	return combine
}
