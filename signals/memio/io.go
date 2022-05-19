package memio

import "sync/atomic"

type MemIO struct {
	v *uint32
}

func newMemIO() *MemIO {
	var v = uint32(0)
	return &MemIO{
		v: &v,
	}
}

func (m *MemIO) Value() bool {
	return atomic.LoadUint32(m.v) == 1
}

func (m *MemIO) Set(v bool) {
	if v {
		atomic.StoreUint32(m.v, 1)
		return
	}
	atomic.StoreUint32(m.v, 0)
}
