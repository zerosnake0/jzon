package jzon

import (
	"reflect"
	"unsafe"
)

// var (
// 	ptrLikeCache = newPtrLike()
// )
//
// type ptrLikeMap map[rtype]bool
//
// type ptrLike struct {
// 	v  atomic.Value
// 	mu sync.Mutex
// }
//
// func newPtrLike() *ptrLike {
// 	var p ptrLike
// 	p.v.Store(ptrLikeMap{})
// 	return &p
// }
//
// func (p *ptrLike) ptrLike(rtype elemRType, obj interface{}) bool {
// 	cache := p.v.Load().(ptrLikeMap)
// 	like, ok := cache[rtype]
// 	if ok {
// 		return like
// 	}
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	cache = p.v.Load().(ptrLikeMap)
// 	like, ok = cache[rtype]
// 	if ok {
// 		return like
// 	}
// 	newCache := make(ptrLikeMap, len(cache))
// 	for k, v := range cache {
// 		newCache[k] = v
// 	}
// 	p.updateCache(newCache, reflect.TypeOf(obj))
// 	p.v.Store(newCache)
// 	return newCache[rtype]
// }
//
// func (p *ptrLike) updateCache(cache ptrLikeMap, typ reflect.Type) {
// 	queue := []reflect.Type{typ}
// 	idx := 0
// 	for idx >= 0 {
// 		typ := queue[idx]
// 		idx -= 1
// 		rType := rtypeOfType(typ)
// 		if _, ok := cache[rType]; ok {
// 			continue
// 		}
// 		switch typ.Kind() {
// 		case reflect.Ptr:
// 			cache[rType] = true
// 		default:
// 			cache[rType] = false
// 		}
// 	}
// }

func (s *Streamer) Value(obj interface{}) *Streamer {
	if s.Error != nil {
		return s
	}
	if obj == nil {
		s.Null()
		return s
	}
	eface := (*eface)(unsafe.Pointer(&obj))
	enc := s.encoder.getEncoderFromCache(eface.rtype)
	if enc == nil {
		typ := reflect.TypeOf(obj)
		enc = s.encoder.createEncoder(eface.rtype, typ)
	}
	enc.Encode(eface.data, s)
	return s
}
