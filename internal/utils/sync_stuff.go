// Copyright (c) 2020-2021 KHS Films
//
// This file is a part of mtproto package.
// See https://github.com/xelaj/mtproto/blob/master/LICENSE for details

package utils

import (
	"reflect"
	"sync"

	"github.com/xelaj/mtproto/internal/encoding/tl"
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

type SyncIntObjectChan struct {
	mutex sync.RWMutex
	m     map[int]chan tl.Object
}

func NewSyncIntObjectChan() *SyncIntObjectChan {
	return &SyncIntObjectChan{m: make(map[int]chan tl.Object)}
}

func (s *SyncIntObjectChan) Has(key int) bool {
	s.mutex.RLock()
	_, ok := s.m[key]
	s.mutex.RUnlock()
	return ok
}

func (s *SyncIntObjectChan) Get(key int) (chan tl.Object, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *SyncIntObjectChan) Add(key int, value chan tl.Object) {
	s.mutex.Lock()
	s.m[key] = value
	s.mutex.Unlock()
}

func (s *SyncIntObjectChan) Keys() []int {
	keys := make([]int, 0, len(s.m))
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

func (s *SyncIntObjectChan) Delete(key int) bool {
	s.mutex.Lock()
	_, ok := s.m[key]
	delete(s.m, key)
	s.mutex.Unlock()
	return ok
}

type SyncIntReflectTypes struct {
	mutex sync.RWMutex
	m     map[int][]reflect.Type
}

func NewSyncIntReflectTypes() *SyncIntReflectTypes {
	return &SyncIntReflectTypes{m: make(map[int][]reflect.Type)}
}

func (s *SyncIntReflectTypes) Has(key int) bool {
	s.mutex.RLock()
	_, ok := s.m[key]
	s.mutex.RUnlock()
	return ok
}

func (s *SyncIntReflectTypes) Get(key int) ([]reflect.Type, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *SyncIntReflectTypes) Add(key int, value []reflect.Type) {
	s.mutex.Lock()
	s.m[key] = value
	s.mutex.Unlock()
}

func (s *SyncIntReflectTypes) Keys() []int {
	keys := make([]int, 0, len(s.m))
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

func (s *SyncIntReflectTypes) Delete(key int) bool {
	s.mutex.Lock()
	_, ok := s.m[key]
	delete(s.m, key)
	s.mutex.Unlock()
	return ok
}
