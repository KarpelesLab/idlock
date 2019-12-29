package idlock

import "sync"

type UintLock struct {
	mp map[uint]bool
	lk sync.Mutex
	cd *sync.Cond
}

func NewUint() *UintLock {
	lk := &UintLock{
		mp: make(map[uint]bool),
	}
	lk.cd = sync.NewCond(&lk.lk)

	return lk
}

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
