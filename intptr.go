package idlock

import "sync"

type IntptrLock struct {
	mp map[intptr]bool
	lk sync.Mutex
	cd *sync.Cond
}

// NewIntptr will instanciate a new IntptrLock instance and return it.
func NewIntptr() *IntptrLock {
	lk := &IntptrLock{
		mp: make(map[intptr]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

// Lock will lock any number of intptrs and return on success. The method
// will not return until all locks can be acquired at the same time.
func (lk *IntptrLock) Lock(i ...intptr) {
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
			// mark intptrs as locked
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
func (lk *IntptrLock) Unlock(i ...intptr) {
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
