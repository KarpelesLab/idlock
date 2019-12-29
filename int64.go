package idlock

import "sync"

type Int64Lock struct {
	mp map[int64]bool
	lk sync.Mutex
	cd *sync.Cond
}

// NewInt64 will instanciate a new Int64Lock instance and return it.
func NewInt64() *Int64Lock {
	lk := &Int64Lock{
		mp: make(map[int64]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

// Lock will lock any number of int64s and return on success. The method
// will not return until all locks can be acquired at the same time.
func (lk *Int64Lock) Lock(i ...int64) {
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
			// mark int64s as locked
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
func (lk *Int64Lock) Unlock(i ...int64) {
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
