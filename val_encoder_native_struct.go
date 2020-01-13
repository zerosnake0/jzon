package jzon

import (
	"reflect"
	"unsafe"
)

type structEncoderBuilder struct {
	encoder *structEncoder
	fields  structFields
}

func (enc *Encoder) newStructEncoder(typ reflect.Type) *structEncoderBuilder {
	fields := describeStruct(typ, enc.tag, enc.onlyTaggedField)
	if len(fields) == 0 {
		return nil
	}
	return &structEncoderBuilder{
		encoder: &structEncoder{},
		fields:  fields,
	}
}

// encoder field info
type encoderFieldInfo struct {
	offsets   []uintptr
	rawField  []byte
	quoted    bool
	omitEmpty bool
	encoder   ValEncoder
}

type encoderFields struct {
	list []encoderFieldInfo
}

func (ef *encoderFields) init(size int) {
	ef.list = make([]encoderFieldInfo, 0, size)
}

func (ef *encoderFields) add(f *field, escapeHtml bool, enc ValEncoder) {
	var rawField []byte
	if escapeHtml {
		rawField = encodeString(rawField, f.name, htmlSafeSet[:])
	} else {
		rawField = encodeString(rawField, f.name, safeSet[:])
	}
	offsets := make([]uintptr, len(f.offsets))
	for i := range f.offsets {
		offsets[i] = f.offsets[i].val
	}
	ef.list = append(ef.list, encoderFieldInfo{
		offsets:   offsets,
		rawField:  rawField,
		quoted:    f.quoted,
		omitEmpty: f.omitEmpty,
		encoder:   enc,
	})
}

// struct encoder
type structEncoder struct {
	fields encoderFields
}

func (enc *structEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (enc *structEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	s.ObjectStart()
OuterLoop:
	for i := range enc.fields.list {
		fi := &enc.fields.list[i]
		curPtr := add(ptr, fi.offsets[0], "struct field")
		// ptr != nil => curPtr != nil
		for _, offset := range fi.offsets[1:] {
			curPtr = *(*unsafe.Pointer)(curPtr)
			if curPtr == nil {
				// the embedded pointer field is nil
				continue OuterLoop
			}
			curPtr = add(curPtr, offset, "struct field")
		}
		if fi.omitEmpty {
			if fi.encoder.IsEmpty(curPtr) {
				continue
			}
		}
		s.RawField(fi.rawField)
		opt := EncOpts{
			Quoted: fi.quoted,
		}
		fi.encoder.Encode(curPtr, s, opt.noescape())
		if s.Error != nil {
			return
		}
	}
	s.ObjectEnd()
}

// no fields to encoder
type emptyStructEncoder struct{}

func (*emptyStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}

func (*emptyStructEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.RawString("{}")
}
