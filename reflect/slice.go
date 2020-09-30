package reflect

import (
	"reflect"

	"github.com/ShrewdSpirit/su/errors"
)

func ReverseSlice(input interface{}) (output interface{}, err error) {
	value := reflect.ValueOf(input)

	if value.Kind() != reflect.Slice {
		err = errors.Newf("expected slice, got %T", value)
		return
	}

	if value.Len() == 0 {
		return
	}

	result := reflect.MakeSlice(reflect.SliceOf(value.Index(0).Type()), 0, value.Cap())

	for i := value.Len() - 1; i >= 0; i-- {
		result = reflect.Append(result, value.Index(i))
	}

	output = result.Interface()

	return
}
