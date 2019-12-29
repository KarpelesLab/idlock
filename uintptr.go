package idlock

import "sync"

type UintptrLock struct {
	mp map[uintptr]bool
	lk sync.Mutex
	cd *sync.Cond
}

// NewUintptr will instanciate a new UintptrLock instance and return it.
func NewUintptr() *UintptrLock {
	lk := &UintptrLock{
		mp: make(map[uintptr]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

// Lock will lock any number of uintptrs and return on success. The method
// will not return until all locks can be acquired at the same time.
func (lk *UintptrLock) Lock(i ...uintptr) {
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
			// mark uintptrs as locked
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
func (lk *UintptrLock) Unlock(i ...uintptr) {
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
