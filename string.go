package idlock

import "sync"

type StringLock struct {
	mp map[string]bool
	lk sync.Mutex
	cd *sync.Cond
}

func NewString() *StringLock {
	lk := &StringLock{
		mp: make(map[string]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

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
			// mark ints as locked
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
