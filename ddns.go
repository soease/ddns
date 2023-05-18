package HelloWorld

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(&HelloWorld{})
}

// Gizmo只是一个例子；可以是你自己的类型
type HelloWorld struct {
}

// 通过CaddyModule方法返回Caddy模块的信息
func (HelloWorld) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.hello_world",
		New: func() caddy.Module { return new(HelloWorld) },
	}
}

// 实现Handler接口
func (g *HelloWorld) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello, World!")
	return nil
}

