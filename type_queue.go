package jzon

import "reflect"

type typeQueue []reflect.Type

func (tq *typeQueue) push(t reflect.Type) {
	*tq = append(*tq, t)
}

func (tq *typeQueue) pushAlsoPtr(t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		*tq = append(*tq, t)
	} else {
		*tq = append(*tq, reflect.PtrTo(t), t)
	}
}

func (tq *typeQueue) pop() (t reflect.Type) {
	q := *tq
	l := len(q)
	if l == 0 {
		return nil
	}
	t = q[l-1]
	*tq = q[:l-1]
	return
}
