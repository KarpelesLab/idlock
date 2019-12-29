package idlock

import "sync"

type IntptrLock struct {
	mp map[intptr]bool
	lk sync.Mutex
	cd *sync.Cond
}

func NewIntptr() *IntptrLock {
	lk := &IntptrLock{
		mp: make(map[intptr]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

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
