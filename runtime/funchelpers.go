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
	return getFuncDetails(runtime.FuncForPC(reflect.ValueOf(f).Pointer()))
}

func GetFuncInfoForPc(pc uintptr) (funcName, funcFile string, funcLine int) {
	return getFuncDetails(runtime.FuncForPC(pc))
}

func GetStackFuncPointer(skipFrames int) uintptr {
	pc, _, _, ok := runtime.Caller(1 + skipFrames)
	if !ok {
		return 0
	}

	return pc
}

func getFuncDetails(f *runtime.Func) (funcName, funcFile string, funcLine int) {
	funcName = f.Name()
	funcFile, funcLine = f.FileLine(f.Entry())

	return
}
