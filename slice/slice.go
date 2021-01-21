package slice

import (
	"reflect"

	"github.com/kiyoptr/su/errors"
)

func Insert(slice interface{}, value interface{}, at int) (result interface{}, err error) {
	panic("not implemented")
	return
}

func Reverse(slice interface{}) (reversed interface{}, err error) {
	value := reflect.ValueOf(slice)

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

	reversed = result.Interface()

	return
}

func RemoveIntSlice(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
