package session

import "sync"

type cacheLoader[T any] struct {
	sync.RWMutex
	l      loader[T]
	cached *T
}

func NewCached(l SessionLoader) SessionLoader {
	return &cacheLoader[Session]{l: l}
}

func (l *cacheLoader[T]) Load() (res T, err error) {
	if c := l.lazyLoad(); c != nil {
		return *c, nil
	}

	l.Lock()
	defer l.Unlock()
	if res, err = l.l.Load(); err != nil {
		return res, err
	}
	l.cached = &res
	return res, nil
}

func (l *cacheLoader[T]) lazyLoad() *T { l.RLock(); defer l.RUnlock(); return l.cached }

func (l *cacheLoader[T]) Store(v T) error {
	l.Lock()
	defer l.Unlock()

	if err := l.l.Store(v); err != nil {
		return err
	}

	l.cached = &v
	return nil
}
