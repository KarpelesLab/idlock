package idlock

import "sync"

type Uint64Lock struct {
	mp map[uint64]bool
	lk sync.Mutex
	cd *sync.Cond
}

// NewUint64 will instanciate a new Uint64Lock instance and return it.
func NewUint64() *Uint64Lock {
	lk := &Uint64Lock{
		mp: make(map[uint64]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

// Lock will lock any number of uint64s and return on success. The method
// will not return until all locks can be acquired at the same time.
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

// Unlock releases all the initially obtained locks. Always release acquired
// locks and never release a lock you didn't acquire.
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
