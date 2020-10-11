package runtime

import (
	"path"
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(f interface{}, getFullName bool) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if getFullName {
		return fullName
	}

	name := fullName

	if strings.Contains(fullName, "/") {
		name = path.Base(fullName)
	}

	if strings.Contains(name, ".") {
		splited := strings.Split(name, ".")
		name = splited[len(splited)-1]
	}

	return name
}

func GetFuncInfo(f interface{}) (funcName, funcFile string, funcLine int) {
	funcInfo := runtime.FuncForPC(reflect.ValueOf(f).Pointer())

	funcName = funcInfo.Name()
	funcFile, funcLine = funcInfo.FileLine(funcInfo.Entry())

	return
}
