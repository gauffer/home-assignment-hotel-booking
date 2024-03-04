package unitofwork

import "sync"

type UnitOfWork interface {
	Begin()
	Commit()
}

type MutexUnitOfWork struct {
	mutex *sync.Mutex
}

func NewMutexUnitOfWork(mutex *sync.Mutex) *MutexUnitOfWork {
	return &MutexUnitOfWork{mutex: mutex}
}

func (u *MutexUnitOfWork) Begin() {
	u.mutex.Lock()
}

func (u *MutexUnitOfWork) Commit() {
	u.mutex.Unlock()
}
