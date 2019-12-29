package idlock

import "sync"

type IntLock struct {
	mp map[int]bool
	lk sync.Mutex
	cd *sync.Cond
}

func NewInt() *IntLock {
	lk := &IntLock{
		mp: make(map[int]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

func (lk *IntLock) Lock(i ...int) {
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

func (lk *IntLock) Unlock(i ...int) {
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
