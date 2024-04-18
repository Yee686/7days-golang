package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

type Int int

func (s Int) Len() int {
	cnt := 0
	for s > 0 {
		s /= 10
		cnt += 1
	}
	return cnt
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("k1"); ok || lru.Len() != 2 {
		t.Fatalf("Remove Oldest key failed")
	}

}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(k string, v Value) {
		keys = append(keys, k)
	}

	lru := New(int64(10), callback)
	lru.Add("key1", Int(643))
	lru.Add("k2", String("v2"))
	lru.Add("k3", String("v3"))
	lru.Add("k4", String("v4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}

	for _, v := range keys {
		fmt.Println(v)
	}
}
