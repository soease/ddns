package ddns

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(Gizmo{})
}

// Gizmo只是一个例子；可以是你自己的类型
type Gizmo struct {
}

// 通过CaddyModule方法返回Caddy模块的信息
func (Gizmo) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "ddns",
		New: func() caddy.Module { return new(Gizmo) },
	}
}

// 实现Handler接口
func (g *Gizmo) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello, World!")
	return nil
}
