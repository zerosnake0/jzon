package jzon

import (
	"reflect"
	"unsafe"
)

type IfaceValEncoder interface {
	Encode(o interface{}, s *Streamer, opts *EncOpts)
}

type IfaceValEncoderConfig struct {
	Type    reflect.Type
	Encoder IfaceValEncoder
}

type ifaceValEncoder struct {
	isEmpty isEmptyFunc
	encoder IfaceValEncoder
	rtype   rtype
}

func (enc *ifaceValEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *ifaceValEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, ptr)
	enc.encoder.Encode(obj, s, opts)
}

type pointerIfaceValEncoder struct {
	encoder IfaceValEncoder
	rtype   rtype
}

func (pointerIfaceValEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(ptr) == nil
}

func (enc *pointerIfaceValEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	ptr = *(*unsafe.Pointer)(ptr)
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, ptr)
	enc.encoder.Encode(obj, s, opts)
}

type directIfaceValEncoder struct {
	isEmpty isEmptyFunc
	encoder IfaceValEncoder
	rtype   rtype
}

func (enc *directIfaceValEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *directIfaceValEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, *(*unsafe.Pointer)(ptr))
	enc.encoder.Encode(obj, s, opts)
}

type dynamicIfaceValEncoder struct {
	encoder IfaceValEncoder
	rtype   rtype
}

func (enc *dynamicIfaceValEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	obj := packIFace(ptr)
	return obj == nil
}

func (enc *dynamicIfaceValEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	// ptr is a pointer to the interface (*I)
	obj := packIFace(ptr)
	if obj == nil {
		s.Null()
		return
	}
	enc.encoder.Encode(obj, s, opts)
}
