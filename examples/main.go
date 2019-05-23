package main

import (
	"github.com/xiusin/php-go/phpgo"
)

func init() {
	phpgo.InitExtension("hello", "1.0.1")
	phpgo.AddConstant("PHP_GO_MODULE_NAME", "hello-module")
	phpgo.NewFuncEntry("Hello","name", func() string {
		return "hello name"
	})
}

// should not run this function
func main() { panic("wtf") }
