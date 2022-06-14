package levelDb

import (
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type harness struct {
	bloom     filter.Filter
	generator filter.FilterGenerator
	filter    []byte
}

func newHarness() *harness {
	bloom := filter.NewBloomFilter(10)
	return &harness{
		bloom:     bloom,
		generator: bloom.NewGenerator(),
	}
}

func (h *harness) add(key []byte) {
	h.generator.Add(key)
}

func (h *harness) addNum(key uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], key)
	h.add(b[:])
}

func (h *harness) build() {
	b := &util.Buffer{}
	h.generator.Generate(b)
	h.filter = b.Bytes()
}

func (h *harness) reset() {
	h.filter = nil
}

func (h *harness) filterLen() int {
	return len(h.filter)
}
