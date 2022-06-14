package levelDb

import (
	"fmt"
	"testing"
)

func TestBloomFilter_Small(t *testing.T) {
	h := newHarness()
	h.add([]byte("hello"))
	h.add([]byte("world"))
	h.build()
	h.add([]byte("aa"))
	h.build()
	fmt.Println(h.bloom.Contains(h.filter, []byte("hello")))
	fmt.Println(h.bloom.Contains(h.filter, []byte("aa")))
}
