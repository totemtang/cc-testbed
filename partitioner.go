package testbed

import (
//"github.com/totemtang/cc-testbed/clog"
)

type Partitioner interface {
	GetPart(key Key) int
	GetPartKey(partIndex int, rank int64) Key
	GetWholeKey(rank int64) Key
	GetKeyArray() []int64
}

type HashPartitioner struct {
	NParts     int64
	pKeysArray []int64
	keyLen     int
	keyRange   []int64
	compKey    []OneKey
}

func NewHashPartitioner(NParts int64, keyRange []int64) *HashPartitioner {
	hp := &HashPartitioner{
		NParts:     NParts,
		pKeysArray: make([]int64, NParts),
		keyLen:     len(keyRange),
		keyRange:   make([]int64, len(keyRange)),
		compKey:    make([]OneKey, len(keyRange)+2*PADDINGINT64),
	}
	hp.compKey = hp.compKey[PADDINGINT64 : hp.keyLen+PADDINGINT64]

	var unit int64 = 1

	for i, v := range keyRange {
		hp.keyRange[i] = v
		if i != 0 {
			unit *= v
		}
	}

	r := hp.keyRange[0] / hp.NParts
	m := hp.keyRange[0] % hp.NParts

	for i := int64(0); i < hp.NParts; i++ {
		hp.pKeysArray[i] = r * unit
		if i < m {
			hp.pKeysArray[i] += unit
		}
	}

	return hp
}

func (hp *HashPartitioner) GetPart(key Key) int {
	k := ParseKey(key, 0)
	return int(int64(k) % hp.NParts)
}

func (hp *HashPartitioner) GetPartKey(partIndex int, rank int64) Key {
	for i := len(hp.keyRange) - 1; i >= 0; i-- {
		hp.compKey[i] = OneKey(rank % hp.keyRange[i])
		rank = rank / hp.keyRange[i]
	}

	hp.compKey[0] = OneKey(int64(hp.compKey[0])*hp.NParts + int64(partIndex))

	//clog.Info("CompKey %v\n", hp.compKey)

	return CKey(hp.compKey)
}

func (hp *HashPartitioner) GetWholeKey(rank int64) Key {
	return hp.GetPartKey(0, rank)
}

func (hp *HashPartitioner) GetKeyArray() []int64 {
	return hp.pKeysArray
}
