package idlock

import "sync"

type Int64Lock struct {
	mp map[int64]bool
	lk sync.Mutex
	cd *sync.Cond
}

func NewInt64() *Int64Lock {
	lk := &Int64Lock{
		mp: make(map[int64]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

func (lk *Int64Lock) Lock(i int64) {
	var f bool

	lk.lk.Lock()
	for {
		if _, f = lk.mp[i]; !f {
			// mark int as locked
			lk.mp[i] = true
			lk.lk.Unlock()
			return
		}

		// wait for an update
		lk.cd.Wait()
	}
}

func (lk *Int64Lock) Unlock(i int64) {
	lk.lk.Lock()
	delete(lk.mp, i)
	lk.cd.Broadcast()
	lk.lk.Unlock()
}
