package util

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type OrderedMap[K constraints.Ordered, V any] struct {
	keys []K
	m    map[K]V
}

type KeyValuePair[K constraints.Ordered, V any] struct {
	Key   K
	Value V
}

func NewOrderedMap[K constraints.Ordered, V any](capacity int) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		keys: make([]K, 0, capacity),
		m:    make(map[K]V),
	}
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	// Add the key to the slice if it doesn't already exist
	if _, ok := om.m[key]; !ok {
		om.keys = append(om.keys, key)
	}
	om.m[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := om.m[key]
	return v, ok
}

func (om *OrderedMap[K, V]) Delete(key K) {
	delete(om.m, key)
	// Find the index of the key in the slice and remove it
	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			break
		}
	}
}

func (om *OrderedMap[K, V]) Len() int {
	return len(om.m)
}

func (om *OrderedMap[K, V]) Keys() []K {
	return om.keys
}

func (om *OrderedMap[K, V]) Greatest() KeyValuePair[K, V] {
	sort.Slice(om.keys, func(i, j int) bool {
		return om.keys[i] > om.keys[j]
	})

	return KeyValuePair[K, V]{
		Key:   om.keys[0],
		Value: om.m[om.keys[0]],
	}
}

func (om *OrderedMap[K, V]) Clear() {
	om.keys = om.keys[:0]
	for k := range om.m {
		delete(om.m, k)
	}
}
