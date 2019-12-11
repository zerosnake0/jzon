package jzon

import (
	"reflect"
	"unsafe"
)

/*
 * Interface decoder is special, when the object is not nil,
 * the internal type cannot be analysed by reflect.TypeOf,
 * the value must be used
 */
type efaceDecoder struct {
}

func (dec *efaceDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	ef := (*eface)(ptr)
	if ef.data == nil {
		// the pointer is a nil pointer
		// or the element is a nil typed pointer (kinda tricky here)
		o, err := it.Read()
		if err != nil {
			return err
		}
		*(*interface{})(ptr) = o
		return nil
	}
	pObj := (*interface{})(ptr)
	obj := *pObj
	typ := reflect.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		/*
		 * Example:
		 *   var o interface{} = 1
		 *   Unmarshal(`"string"`, &o)
		 */
		o, err := it.Read()
		if err != nil {
			return err
		}
		*pObj = o
		return nil
	}
	// obj is pointer
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	ptrElemType := typ.Elem()
	if c == 'n' {
		it.head += 1
		if err := it.expectBytes("ull"); err != nil {
			return err
		}
		// we have already check above so that
		// obj is not nil
		if ptrElemType.Kind() != reflect.Ptr {
			/*
			 * Example:
			 *   i := 1
			 *   var o interface{} = &i
			 *   Unmarshal(`null`, &o)
			 */
			*pObj = nil
			return nil
		}
		/*
		 * Example:
		 *   i := 1
		 *   pi := &i
		 *   var o interface{} = &pi
		 *   Unmarshal(`null`, &o)
		 */
		*pObj = reflect.New(ptrElemType).Interface()
		return nil
	}
	// when we arrive here, we have:
	// 1 obj is pointer
	// 2 obj != nil
	if err := it.ReadVal(obj); err != nil {
		return err
	}
	*pObj = obj
	return nil
}

type ifaceDecoder struct {
}

func (dec *ifaceDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		if err = it.expectBytes("ull"); err != nil {
			return err
		}
		*(*interface{})(ptr) = nil
		return nil
	}
	o := packIFace(ptr)
	if o == nil {
		return IFaceError
	}
	return it.ReadVal(o)
}
