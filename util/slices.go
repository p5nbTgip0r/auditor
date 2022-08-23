package util

import (
	mapset "github.com/deckarep/golang-set/v2"
)

func SliceDiff[T comparable](old, new []T) (added, removed mapset.Set[T]) {
	oldSet := mapset.NewSet[T]()
	newSet := mapset.NewSet[T]()
	for _, t := range old {
		oldSet.Add(t)
	}
	for _, t := range new {
		newSet.Add(t)
	}

	added = newSet.Difference(oldSet)
	removed = oldSet.Difference(newSet)
	return
}

func SliceDiffIdentifier[I comparable, T interface{}](old, new []T, getId func(T) I) (added, removed map[I]T) {
	oldMap := make(map[I]T, len(old))
	newMap := make(map[I]T, len(new))
	added = make(map[I]T, 0)
	removed = make(map[I]T, 0)
	for _, t := range old {
		oldMap[getId(t)] = t
	}
	for _, t := range new {
		newMap[getId(t)] = t
	}

	for id, t := range newMap {
		_, ok := oldMap[id]
		if !ok {
			added[id] = t
		}
	}
	for id, t := range oldMap {
		_, ok := newMap[id]
		if !ok {
			removed[id] = t
		}
	}
	return
}
