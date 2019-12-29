package idlock

import "sync"

type Uint64Lock struct {
	mp map[uint64]bool
	lk sync.Mutex
	cd *sync.Cond
}

func NewUint64() *Uint64Lock {
	lk := &Uint64Lock{
		mp: make(map[uint64]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

func (lk *Uint64Lock) Lock(i ...uint64) {
	var f bool

	if len(i) == 0 {
		return
	}

	lk.lk.Lock()
	for {
		for _, n := range i {
			if _, f = lk.mp[n]; f {
				break
			}
		}
		if !f {
			// mark uint64s as locked
			for _, n := range i {
				lk.mp[n] = true
			}
			lk.lk.Unlock()
			return
		}

		// wait for an update
		lk.cd.Wait()
	}
}

func (lk *Uint64Lock) Unlock(i ...uint64) {
	if len(i) == 0 {
		return
	}

	lk.lk.Lock()

	for _, n := range i {
		delete(lk.mp, n)
	}

	lk.cd.Broadcast()
	lk.lk.Unlock()
}
