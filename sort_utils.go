package request

import (
	"sort"
)

// RequestSlice attaches the methods of sort.Interface to []Request, sorting in increasing order.
type RequestSlice []Request

func (p RequestSlice) Len() int { return len(p) }

func (p RequestSlice) Less(i, j int) bool {
	for k := 0; k < len(p[i]); k++ {
		if p[i][k] == p[j][k] {
			continue
		}
		return p[i][k] < p[j][k]
	}
	return false
}

func (p RequestSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// SortRequests sorts a slice of StoreService in increasing order.
func SortRequests(a []Request) { sort.Sort(RequestSlice(a)) }

// GetSortedKeys gets Views and transforms it`s keys to sorted slice ([]Request)
func GetSortedKeys(v Views) []Request {
	if v == nil {
		return nil
	}
	keys := make([]Request, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	SortRequests(keys)

	return keys
}
