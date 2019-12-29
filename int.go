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

func (lk *IntLock) Lock(i int) {
	var f bool

	lk.lk.Lock()
	for {
		if _, f = lk.mp[i]; !f {
			// mark int as locked
			lk.mp[i] = true
			lk.lk.Unlock()
			return
		}

		// wait for an update
		lk.cd.Wait()
	}
}

func (lk *IntLock) Unlock(i int) {
	lk.lk.Lock()
	delete(lk.mp, i)
	lk.cd.Broadcast()
	lk.lk.Unlock()
}
