package idlock

import "sync"

type StringLock struct {
	mp map[string]bool
	lk sync.Mutex
	cd *sync.Cond
}

// NewString will instanciate a new StringLock instance and return it.
func NewString() *StringLock {
	lk := &StringLock{
		mp: make(map[string]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

// Lock will lock any number of strings and return on success. The method
// will not return until all locks can be acquired at the same time.
func (lk *StringLock) Lock(i ...string) {
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
			// mark strings as locked
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
func (lk *StringLock) Unlock(i ...string) {
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
