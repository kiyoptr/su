package tagprovider

import "fmt"

func Constant(key string, value interface{}) Provider {
	return func() (string, string) {
		return key, fmt.Sprintf("%v", value)
	}
}
