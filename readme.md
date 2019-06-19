### php-go ###

使用go来编写PHP的扩展 (自用学习cGO编程) fork 来自 [kitech/php-go]()

### 环境 ###

* PHP 5.5+/7.x
* php-dev 安装 (`apt-get install php-dev`)
* go 版本 1.4+,  MacOS, go 1.8+

### 构建和安装 ###

- go get 方式:

    ```bash
    go get github.com/xiusin/php-go
    cd $GOPATH/src/github.com/xiusin/php-go
    // 个别系统可能出现php-config与php不一致的情况. 建议都设置PHPCFG
    PHPCFG=`which php-config` make
    ```

- git方式:

    ```
    mkdir -p $GOPATH/src/github.com/xiusin
    git clone https://github.com/xiusin/php-go.git $GOPATH/src/github.com/xiusin/php-go
    cd $GOPATH/src/github.com/xiusin/php-go
    make
    ls -lh php-go/hello.so
    php -d extension=./hello.so examples/hello.php
    ```

### 实例 ###

```go
// package main is required
package main

import "github.com/xiusin/php-go/phpgo"

func foo_in_go() {
}

type Bar struct{}
func NewBar() *Bar{
    return &Bar{}
}

func init() {
    phpgo.InitExtension("mymod", "1.0")
    phpgo.AddFunc("foo_in_php", foo_in_go)
    phpgo.AddClass("bar_in_php", NewBar)
}

// should not run this function
// required for go build though.
func main() { panic("wtf") }
```

### TODO ###

- [ ] 使用go get编译安装扩展
- [x] 改进php7支持
- [ ] 命名空间支持
- [ ] 多扩展支持
- [ ] 类成员权限访问支持
- [x] 不限数量 function/method/class 支持
- [x] 全局ini变量支持
- [ ] 填充phpinfo
- [ ] 支持函数回调
