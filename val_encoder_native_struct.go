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
	if fields.count() == 0 {
		return nil
	}
	return &structEncoderBuilder{
		encoder: &structEncoder{},
		fields:  fields,
	}
}

// encoder field info
type encoderFieldInfo struct {
	offsets  []uintptr
	rawField []byte
	quoted   bool
	encoder  ValEncoder
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
	ef.list = append(ef.list, encoderFieldInfo{
		offsets:  f.offsets,
		rawField: rawField,
		quoted:   f.quoted,
		encoder:  enc,
	})
}

// struct encoder
type structEncoder struct {
	fields encoderFields
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
	for i := range enc.fields.list {
		fi := &enc.fields.list[i]
		curPtr := add(ptr, fi.offsets[0], "struct field")
		// ptr != nil => curPtr != nil
		broken := false
		for _, offset := range fi.offsets[1:] {
			curPtr = *(*unsafe.Pointer)(curPtr)
			if curPtr == nil {
				broken = true // the embedded pointer field is nil
				break
			}
			curPtr = add(curPtr, offset, "struct field")
		}
		if !broken {
			s.RawField(fi.rawField)
			opt := EncOpts{
				Quoted: fi.quoted,
			}
			fi.encoder.Encode(curPtr, s, &opt)
			if s.Error != nil {
				return
			}
		}
	}
	s.ObjectEnd()
}

// no fields to encoder
type emptyStructEncoder struct {
}

func (enc *emptyStructEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.RawString("{}")
}
