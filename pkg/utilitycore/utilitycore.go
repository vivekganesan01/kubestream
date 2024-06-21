package utilitycore

import (
	"reflect"
	"runtime"
)

func GetFn(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
