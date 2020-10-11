package runtime

import "testing"

func TestGetFuncName(t *testing.T) {
	t.Log(GetFuncName(TestGetFuncName, true))
}

func TestGetFuncInfo(t *testing.T) {
	t.Log(GetFuncInfo(TestGetFuncInfo))
}

func TestGetStackFuncPointer(t *testing.T) {
	t.Log(GetFuncInfoForPc(GetStackFuncPointer(1)))
}
