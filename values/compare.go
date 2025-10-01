package values

import (
	"fmt"
	"reflect"
)

var (
	int64Type   = reflect.TypeOf(int64(0))
	float64Type = reflect.TypeOf(float64(0))
)

// Equal returns a bool indicating whether a == b after conversion.
func Equal(a, b any) bool { //nolint: gocyclo
	a, b = ToLiquid(a), ToLiquid(b)
	if a == nil || b == nil {
		return a == b
	}
	ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
	fmt.Println(ra.Kind(), rb.Kind())
	switch joinKind(ra.Kind(), rb.Kind()) {
	case reflect.Array, reflect.Slice:
		if ra.Len() != rb.Len() {
			return false
		}
		for i := range ra.Len() {
			if !Equal(ra.Index(i).Interface(), rb.Index(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Bool:
		return ra.Bool() == rb.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ra.Convert(int64Type).Int() == rb.Convert(int64Type).Int()
	case reflect.Float32, reflect.Float64:
		return ra.Convert(float64Type).Float() == rb.Convert(float64Type).Float()
	case reflect.String:
		return ra.String() == rb.String()
	case reflect.Ptr:
		if rb.Kind() == reflect.Ptr && (ra.IsNil() || rb.IsNil()) {
			return ra.IsNil() == rb.IsNil()
		}
		return a == b
	case reflect.Map:
		return mapsAreEqual(ra.Interface().(map[string]any), rb.Interface().(map[string]any))
	default:
		return a == b
	}
}

// Less returns a bool indicating whether a < b.
func Less(a, b any) bool {
	a, b = ToLiquid(a), ToLiquid(b)
	if a == nil || b == nil {
		return false
	}
	ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
	switch joinKind(ra.Kind(), rb.Kind()) {
	case reflect.Bool:
		return !ra.Bool() && rb.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ra.Convert(int64Type).Int() < rb.Convert(int64Type).Int()
	case reflect.Float32, reflect.Float64:
		return ra.Convert(float64Type).Float() < rb.Convert(float64Type).Float()
	case reflect.String:
		return ra.String() < rb.String()
	default:
		return false
	}
}

func joinKind(a, b reflect.Kind) reflect.Kind { //nolint: gocyclo
	if a == b {
		return a
	}
	switch a {
	case reflect.Array, reflect.Slice:
		if b == reflect.Array || b == reflect.Slice {
			return reflect.Slice
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isIntKind(b) {
			return reflect.Int64
		}
		if isFloatKind(b) {
			return reflect.Float64
		}
	case reflect.Float32, reflect.Float64:
		if isIntKind(b) || isFloatKind(b) {
			return reflect.Float64
		}
	}
	return reflect.Invalid
}

func isIntKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isFloatKind(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// mapsAreEqual checks if two maps have the same keys and values.
func mapsAreEqual(a, b map[string]any) bool {
	// 1. Check if they have the same number of keys.
	if len(a) != len(b) {
		return false
	}

	// 2. Check if all keys and values from 'a' exist and are equal in 'b'.
	for key, valA := range a {
		valB, ok := b[key]
		// 2a. Does the key from 'a' even exist in 'b'?
		if !ok {
			return false
		}
		// 2b. Are the values for that key equal?
		// (This uses reflect.DeepEqual to handle nested objects)
		if !reflect.DeepEqual(valA, valB) {
			return false
		}
	}

	// If all checks pass, the maps are equal.
	return true
}
