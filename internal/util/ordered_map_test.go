package util_test

import (
	"testing"

	"github.com/assemblaj/ggpo/internal/util"
)

func TestOrderedMapLenZero(t *testing.T) {
	om := util.NewOrderedMap[int, uint32](0)
	omSize := om.Len()
	if omSize != 0 {
		t.Errorf("Expected Size 0, got %d", omSize)
	}

}
func TestOrderedMapSet(t *testing.T) {
	om := util.NewOrderedMap[int, uint32](2)
	key := 8
	var value uint32 = 5555
	om.Set(key, value)
	if om.Len() == 0 {
		t.Errorf("Expected to have value, instead was empty")
	}
	want := value
	got, ok := om.Get(key)
	if !ok {
		t.Errorf("Expected value for key %d, value not available", key)
	}

	if want != got {
		t.Errorf("Wanted %d got %d", want, got)
	}
}

func TestOrderedMapDelete(t *testing.T) {
	om := util.NewOrderedMap[int, uint32](2)
	key := 8
	var value uint32 = 5555
	om.Set(key, value)
	om.Delete(key)
	if om.Len() != 0 {
		t.Errorf("Expected to be empty, instead not empty")
	}
	if len(om.Keys()) != 0 {
		t.Errorf("Expected keys to be empty, instead not empty")
	}
}

func TestOrderedMapClear(t *testing.T) {
	om := util.NewOrderedMap[int, uint32](3)
	key1, key2 := 898, 343
	var value1, value2 uint32 = 32423423, 234234234
	om.Set(key1, value1)
	om.Set(key2, value2)
	om.Clear()
	if om.Len() != 0 {
		t.Errorf("Expected to be empty, instead not empty")
	}
	if len(om.Keys()) != 0 {
		t.Errorf("Expected keys to be empty, instead not empty")
	}
}

func TestOrderedMapGreatest(t *testing.T) {
	om := util.NewOrderedMap[int, uint32](3)
	iter := 20
	for i := 0; i <= iter; i++ {
		om.Set(i, uint32(i*100))
	}
	kv := om.Greatest()
	if kv.Key != iter {
		t.Errorf("Expected greatest key to be %d, insteat %d", iter, kv.Key)
	}
	if kv.Value != uint32(iter*100) {
		t.Errorf("Expected value for greatest key to be %d, was instead %d", iter*100, kv.Value)
	}
}

func TestOrderedMapGet(t *testing.T) {
	om := util.NewOrderedMap[int, uint32](3)
	key1, key2 := 898, 343
	var value1, value2 uint32 = 32423423, 234234234
	om.Set(key1, value1)
	om.Set(key2, value2)
	v1, ok := om.Get(key1)
	if !ok {
		t.Errorf("Value should be avaiable but not avaiable.")
	}
	if v1 != value1 {
		t.Errorf("Expected %d, got %d ", v1, value1)
	}
	v2, ok := om.Get(key2)
	if !ok {
		t.Errorf("Value should be avaiable but not avaiable.")
	}
	if v2 != value2 {
		t.Errorf("Expected %d, got %d ", v1, value1)
	}
}
