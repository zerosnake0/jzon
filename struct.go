/* see encoding/json
 * - some additional comments may be added
 * - some code may be slightly modified
 */
package jzon

import (
	"reflect"
	"sort"
)

func describeStruct(st reflect.Type, tagKey string, onlyTaggedField bool) structFields {
	// Anonymous fields to explore at the current level and the next.
	var current []field
	next := []field{{
		typ: st,
		offsets: []offset{{
			val: 0,
		}},
	}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		// move to next level
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			/* 1. First of all we can not have something like
			 *        type A struct {
			 *            A
			 *        }
			 *    or something like
			 *        type A struct {    type B struct {
			 *            B                  A
			 *        }                  }
			 *    we must have something like
			 *        type A struct {
			 *            *A
			 *        }
			 *    of course there can be more embedded levels, but at least one
			 *    pointer is required
			 *
			 * 2. Next, when we have a struct embedded by itself, according to the
			 *    field sorting which will be applied after, the embedded one will
			 *    have lower priority than the parent one, so we just skip it here
			 */
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.PkgPath != "" { // the field is not exported, i.e the field begins with lowercase
					if sf.Anonymous { // embedded field
						t := sf.Type
						if t.Kind() == reflect.Ptr {
							t = t.Elem()
						}
						/*
						 * We can only have
						 *     type A struct {
						 *         *B
						 *     }
						 * but not
						 *     type A struct {
						 *         **B
						 *     }
						 * (even if the **B is defined by type aliasing)
						 */
						if t.Kind() != reflect.Struct {
							continue
						}
					} else { // not embedded, just skip
						continue
					}
				}

				var (
					// `json:"<name>,<opts>"`
					name string
					opts tagOptions
				)
				tag, ok := sf.Tag.Lookup(tagKey)
				if ok {
					if tag == "-" {
						continue
					}
					name, opts = parseTag(tag)
					if !isValidTag(name) {
						name = ""
					}
				} else {
					if onlyTaggedField && !sf.Anonymous {
						continue
					}
				}

				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				/*
				 * When ft.Name() == "", we may have:
				 *   1. a pointer
				 *   2. an anonymous struct like
				 *      type A struct {
				 *          B struct { ... }
				 *      }
				 *   3. other case?
				 */
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}

				l := len(f.offsets)
				offsets := make([]offset, l, l+1)
				copy(offsets, f.offsets)
				offsets[l-1].val += sf.Offset

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					/* either:
					 *     1. the field has a json tag
					 *     2. the field is not embedded
					 *     3. the field type is not struct (or pointer of struct)
					 */
					ptrType := reflect.PtrTo(sf.Type)
					field := field{
						index:     index,
						offsets:   offsets,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),

						ptrType: ptrType,
						// rtype:   rtypeOfType(ptrType),
					}

					if name == "" {
						field.name = sf.Name
						field.tagged = false
					} else {
						field.name = name
						field.tagged = true
					}
					field.nameBytes = []byte(field.name)
					field.nameBytesUpper = toUpper(field.nameBytes, nil)
					field.equalFold = foldFunc(field.nameBytes)

					// Only strings, floats, integers, and booleans can be quoted.
					if opts.Contains("string") {
						switch ft.Kind() {
						case reflect.Bool,
							reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
							reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
							reflect.Float32, reflect.Float64,
							reflect.String:
							field.quoted = true
						}
					}

					fields = append(fields, field)
					if count[f.typ] > 1 {
						/* when we arrived here, we are inside a level, where there are at least
						 * two embedded field with a same (final) type
						 * but the embedded field type will be analysed only once (this is ensured
						 * by both:
						 *     1. the visited map above
						 *     2. the nextCount check below
						 *
						 * but there will be multiple fields with same:
						 *     1. json field name
						 *     2. field depth
						 *     3. field tagged or not
						 * but with different:
						 *     1. index (or index array)
						 *
						 * The fields are sorted by:
						 *     1. json field name
						 *     2. depth
						 *     3. tagged or not (tagged is less)
						 *     4. the index array ([0,0,0] < [0,0,1])
						 *
						 * For each json field name, the dominant field is selected between
						 * the two first elements with the same json field name:
						 *     1. the one with lesser depth wins, otherwise
						 *     2. the one with json tag wins, otherwise
						 *     3. there is no dominant field
						 *
						 * Hence, there will be no dominant one among these multiple similar fields.
						 * But we still add the same field once more, because in this case the embedded
						 * field with a same type and a deeper depth should also be ignored, for example:
						 *     type A struct {}
						 *     type B = A
						 *     type C struct { A }
						 *     type D struct {
						 *         C
						 *         B
						 *         A
						 *     }
						 * we can just add the same field again because they will be eliminated in the end
						 * we will only add once more because the same struct type is analysed only once
						 */
						fields = append(fields, field)
					}
					continue
				}

				nextCount[ft]++
				if nextCount[ft] == 1 {
					/* one example for nextCount[ft] > 1
					 *     type A struct {}
					 *     type B = A
					 *     type C struct {
					 *         B
					 *         A
					 *     }
					 * in this case the name is the name of type (A)
					 */
					if sf.Type.Kind() == reflect.Ptr {
						var elemRType rtype
						if sf.PkgPath == "" { // the field is exported
							elemRType = rtypeOfType(sf.Type.Elem())
						}
						offsets = append(offsets, offset{
							val:   0,
							rtype: elemRType,
						})
					}
					next = append(next, field{
						name:    ft.Name(),
						index:   index,
						offsets: offsets,
						typ:     ft,
					})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if x[i].name != x[j].name {
			return x[i].name < x[j].name
		}
		if len(x[i].index) != len(x[j].index) {
			return len(x[i].index) < len(x[j].index)
		}
		if x[i].tagged != x[j].tagged {
			return x[i].tagged
		}
		return byIndex(x).Less(i, j)
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	return fields
}
