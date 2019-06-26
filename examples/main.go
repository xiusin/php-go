package main

import (
	"github.com/xiusin/php-go/phpgo"
	"reflect"
)

func init() {

	phpgo.InitExtension("hello", "1.0.1")
	_ = phpgo.AddConstant("PHP_GO_MODULE_NAME", "hello-module")

	phpgo.NewFuncEntry("Hello", "name", func() string {
		return "hello name"
	})

	_ = phpgo.AddFunc("add_call_back", func(call interface{}) string {
		return reflect.TypeOf(call).Kind().String()
	})

}

// should not run this function
func main() { panic("wtf") }
