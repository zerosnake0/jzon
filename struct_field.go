/* see encoding/json
 * - some additional comments may be added
 * - some code may be slightly modified
 */
package jzon

import (
	"reflect"
)

type offset struct {
	val   uintptr
	rtype rtype
}

type field struct {
	typ     reflect.Type
	index   []int
	offsets []offset

	name           string
	nameBytes      []byte // []byte(name)
	nameBytesUpper []byte
	equalFold      func(s, t []byte) bool // bytes.EqualFold or equivalent

	tagged    bool
	quoted    bool
	omitEmpty bool

	ptrType reflect.Type
	// rtype   rtype
}

type structFields []field

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

func dominantField(fields []field) (field, bool) {
	// The fields are sorted in increasing index-length order, then by presence of tag.
	// That means that the first field is the dominant one. We need only check
	// for error cases: two fields at top level, either both tagged or neither tagged.
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tagged == fields[1].tagged {
		return field{}, false
	}
	return fields[0], true
}
