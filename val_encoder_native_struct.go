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
	offsets := make([]uintptr, len(f.offsets))
	for i := range f.offsets {
		offsets[i] = f.offsets[i].val
	}
	ef.list = append(ef.list, encoderFieldInfo{
		offsets:  offsets,
		rawField: rawField,
		quoted:   f.quoted,
		encoder:  enc,
	})
}

type encoderFieldInfo2 struct {
	index    []int
	rawField []byte
	quoted   bool
	encoder  ValEncoder2
}

type encoderFields2 struct {
	list []encoderFieldInfo2
}

func (ef *encoderFields2) init(size int) {
	ef.list = make([]encoderFieldInfo2, 0, size)
}

func (ef *encoderFields2) add(f *field, escapeHtml bool, enc ValEncoder2) {
	var rawField []byte
	if escapeHtml {
		rawField = encodeString(rawField, f.name, htmlSafeSet[:])
	} else {
		rawField = encodeString(rawField, f.name, safeSet[:])
	}
	ef.list = append(ef.list, encoderFieldInfo2{
		index:    f.index,
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
		s.RawField(fi.rawField)
		opt := EncOpts{
			Quoted: fi.quoted,
		}
		fi.encoder.Encode(curPtr, s, &opt)
		if s.Error != nil {
			return
		}
	}
	s.ObjectEnd()
}

type structEncoderBuilder2 struct {
	encoder *structEncoder2
	fields  structFields
}

func (enc *Encoder) newStructEncoder2(typ reflect.Type) *structEncoderBuilder2 {
	fields := describeStruct(typ, enc.tag, enc.onlyTaggedField)
	if fields.count() == 0 {
		return nil
	}
	return &structEncoderBuilder2{
		encoder: &structEncoder2{},
		fields:  fields,
	}
}

// struct encoder
type structEncoder2 struct {
	fields encoderFields2
}

func (enc *structEncoder2) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	s.ObjectStart()
OuterLoop:
	for i := range enc.fields.list {
		fi := &enc.fields.list[i]
		field := v
		l := len(fi.index) - 1
		for i := 0; i < l; i++ {
			field = field.Field(fi.index[i])
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					continue OuterLoop
				}
				field = field.Elem()
			}
		}
		field = field.Field(fi.index[l])
		s.RawField(fi.rawField)
		opt := EncOpts{
			Quoted: fi.quoted,
		}
		fi.encoder.Encode2(field, s, &opt)
		if s.Error != nil {
			return
		}
	}
	s.ObjectEnd()
}

// no fields to encoder
type emptyStructEncoder struct{}

func (*emptyStructEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.RawString("{}")
}

func (*emptyStructEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	s.RawString("{}")
}
