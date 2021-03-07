// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package utils

import (
	"sync"
)

type SyncSetInt struct {
	mutex sync.RWMutex
	m     map[int]null
}

func NewSyncSetInt() *SyncSetInt {
	return &SyncSetInt{m: make(map[int]null)}
}

func (s *SyncSetInt) Has(key int) bool {
	s.mutex.RLock()
	_, ok := s.m[key]
	s.mutex.RUnlock()
	return ok
}

func (s *SyncSetInt) Add(key int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.m[key]
	s.m[key] = null{}
	return !ok
}

func (s *SyncSetInt) Delete(key int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.m[key]
	delete(s.m, key)
	return ok
}

func (s *SyncSetInt) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m = make(map[int]null)
}
