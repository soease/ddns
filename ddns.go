package ddns

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(&DDNS{})
}

type DDNS struct {
}

// 通过CaddyModule方法返回Caddy模块的信息
func (DDNS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.hello_world",
		New: func() caddy.Module { return new(DDNS) },
	}
}

// 实现Handler接口
func (g *DDNS) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello, World!")
	return nil
}
