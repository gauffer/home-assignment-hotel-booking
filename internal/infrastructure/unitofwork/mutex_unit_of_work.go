package unitofwork

import "sync"

type UnitOfWork interface {
	Begin()
	Commit()
}

type mutexUnitOfWork struct {
	mu *sync.Mutex
}

func NewMutexUnitOfWork(mutex *sync.Mutex) *mutexUnitOfWork {
	return &mutexUnitOfWork{mu: mutex}
}

func (u *mutexUnitOfWork) Begin() {
	u.mu.Lock()
}

func (u *mutexUnitOfWork) Commit() {
	u.mu.Unlock()
}
