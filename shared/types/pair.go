package types

import "iter"

type Pair[FirstType any, SecondType any] struct {
	First  FirstType
	Second SecondType
}

func PairsIterator[FirstType any, SecondType any](pairs []Pair[FirstType, SecondType]) iter.Seq2[FirstType, SecondType] {
	return func(yield func(FirstType, SecondType) bool) {
		for _, p := range pairs {
			if !yield(p.First, p.Second) {
				return
			}
		}
	}
}
