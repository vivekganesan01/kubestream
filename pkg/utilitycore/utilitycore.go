package utilitycore

import (
	"reflect"
	"runtime"
)

func GetFn(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Int32Ptr(i int32) *int32 { return &i }
