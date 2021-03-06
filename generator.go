package testbed

import (
	"flag"
	"math/rand"
	"time"

	//"github.com/totemtang/cc-testbed/clog"
)

var CrossPercent = flag.Float64("cr", 0.0, "percentage of cross-partition transactions")

type TxnGen struct {
	padding1    [64]byte
	ID          int
	TXN         int
	nKeys       int64
	nParts      int
	partIndex   int
	rr          float64
	txnLen      int
	maxParts    int
	isPartition bool
	rnd         *rand.Rand
	zk          *ZipfKey
	q           *Query
	numAccess   int
	padding2    [64]byte
}

func NewTxnGen(ID int, TXN int, rr float64, txnLen int, maxParts int, zk *ZipfKey) *TxnGen {
	txnGen := &TxnGen{
		ID:          ID,
		TXN:         TXN,
		nKeys:       zk.nKeys,
		nParts:      zk.nParts,
		partIndex:   zk.partIndex,
		rr:          rr,
		txnLen:      txnLen,
		maxParts:    maxParts,
		isPartition: zk.isPartition,
		zk:          zk,
	}
	q := &Query{
		TXN:         txnGen.TXN,
		txnLen:      txnGen.txnLen,
		isPartition: txnGen.isPartition,
		partitioner: txnGen.zk.hp,
		accessParts: make([]int, 0, txnGen.maxParts+32),
		rKeys:       make([]Key, 0, txnGen.txnLen+16),
		wKeys:       make([]Key, 0, txnGen.txnLen+16),
	}
	q.rKeys = q.rKeys[8:8]
	q.wKeys = q.wKeys[8:8]
	q.accessParts = q.accessParts[16:16]
	txnGen.q = q

	//txnGen.local_seed = uint32(rand.Intn(10000000))
	txnGen.rnd = rand.New(rand.NewSource(time.Now().Unix() / int64(ID+1)))

	var numAccess int
	if txnGen.maxParts > *NumPart {
		txnGen.maxParts = *NumPart
	}
	if txnGen.maxParts < txnGen.txnLen {
		numAccess = txnGen.maxParts
	} else {
		numAccess = txnGen.txnLen
	}
	txnGen.numAccess = numAccess

	return txnGen
}

//Determine a read or a write operation
func insertRWKey(q *Query, k Key, rr float64, rnd *rand.Rand) {
	x := float64(rnd.Int63n(100))
	if x < rr {
		q.rKeys = append(q.rKeys, k)
	} else {
		q.wKeys = append(q.wKeys, k)
	}
}

func (tg *TxnGen) GenOneQuery() *Query {
	q := tg.q
	q.rKeys = q.rKeys[:0]
	q.wKeys = q.wKeys[:0]

	// Generate keys for different CC
	if tg.isPartition {
		//x := float64(RandN(&tg.local_seed, 100))
		x := float64(tg.rnd.Int63n(100))
		if x < *CrossPercent && tg.numAccess > 1 {
			// Generate how many partitions this txn will touch; more than 1

			//numAccess := tg.rnd.Intn(tg.numAccess-1) + 2
			numAccess := 2
			//q.accessParts = make([]int, numAccess)
			q.accessParts = q.accessParts[:numAccess]

			// Generate partitions this txn will touch
			var remotePart int
			remotePart = (tg.partIndex + tg.rnd.Intn(tg.nParts-1) + 1) % tg.nParts
			if remotePart > tg.partIndex {
				q.accessParts[0] = tg.partIndex
				q.accessParts[1] = remotePart
			} else {
				q.accessParts[0] = remotePart
				q.accessParts[1] = tg.partIndex
			}

			/*
				if tg.partIndex+numAccess <= tg.nParts {
					for i := 0; i < numAccess; i++ {
						q.accessParts[i] = tg.partIndex + i
					}
				} else {
					for i := 0; i < numAccess; i++ {
						tmp := tg.partIndex + i
						if tmp >= tg.nParts {
							q.accessParts[tmp-tg.nParts] = tmp - tg.nParts
						} else {
							q.accessParts[numAccess-tg.nParts+tg.partIndex+i] = tmp
						}
					}
				}
			*/
			/*for i := 0; i < numAccess; i++ {
				if tg.partIndex+numAccess <= tg.nParts {
					q.accessParts[i] = tg.partIndex + i
				} else {
					tmp := tg.partIndex + i
					if tmp >= tg.nParts {
						q.accessParts[tmp-tg.nParts] = tmp - tg.nParts
					} else {
						q.accessParts[numAccess-tg.nParts+tg.partIndex+i] = tmp
					}
				}
			}*/

			var j int = 0
			for i := 0; i < tg.txnLen; i++ {
				insertRWKey(q, tg.zk.GetOtherKey(q.accessParts[j]), tg.rr, tg.rnd)
				j = (j + 1) % numAccess
			}

		} else {
			//q.accessParts = make([]int, 1)
			q.accessParts = q.accessParts[:1]
			q.accessParts[0] = tg.partIndex

			for i := 0; i < tg.txnLen; i++ {
				insertRWKey(q, tg.zk.GetSelfKey(), tg.rr, tg.rnd)
			}
		}

	} else {
		// Generate random keys
		for i := 0; i < tg.txnLen; i++ {
			insertRWKey(q, tg.zk.GetKey(), tg.rr, tg.rnd)
		}
	}

	// Generate values according to the transaction type
	q.GenValue(tg.rnd)

	return q
}
