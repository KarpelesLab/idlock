package idlock

import "sync"

type UintLock struct {
	mp map[uint]bool
	lk sync.Mutex
	cd *sync.Cond
}

// NewUint will instanciate a new UintLock instance and return it.
func NewUint() *UintLock {
	lk := &UintLock{
		mp: make(map[uint]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

// Lock will lock any number of uints and return on success. The method
// will not return until all locks can be acquired at the same time.
func (lk *UintLock) Lock(i ...uint) {
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
			// mark uints as locked
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
func (lk *UintLock) Unlock(i ...uint) {
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
