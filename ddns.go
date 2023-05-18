package wakeonlan

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("wake_on_lan", parseCaddyfile)
}

// 结构体使用 wake-on-lan 技术，可以在 HTTP 请求时唤醒目标主机。
type Middleware struct {
	MAC              string `json:"mac,omitempty"`               // 目标主机的 MAC 地址，格式应与 net.ParseMAC 兼容。
	BroadcastAddress string `json:"broadcast_address,omitempty"` // 魔术包应发送到的广播地址（<ip>:<port>）。默认为 "255.255.255.255:9"。

	key             string
	logger          *zap.Logger
	pool            *caddy.UsagePool
	magicPacket     []byte
	broadcastSocket net.Conn
}

// 返回 Caddy 模块信息
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.wake_on_lan",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Prepare 函数用于准备唤醒信号和用于发送信号的 socket。
// 函数实现 caddy.Provisioner 接口
func (m *Middleware) Provision(ctx caddy.Context) error {
	m.key = fmt.Sprintf("wol-%s", m.MAC)
	m.logger = ctx.Logger(m)
	m.pool = caddy.NewUsagePool()

	mac, err := net.ParseMAC(m.MAC)
	if err != nil {
		return err
	}
	m.magicPacket = BuildMagicPacket(mac)
	if err != nil {
		return err
	}
	m.broadcastSocket, err = net.Dial("udp", m.BroadcastAddress)
	if err != nil {
		return err
	}
	return nil
}

// 函数用于发送准备好的唤醒信号，并透明地继续下一个 HTTP 处理程序。
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	_, throttled := m.pool.LoadOrStore(m.key, true)
	if throttled {
		_, err := m.pool.Delete(m.key)
		if err != nil {
			return err
		}
	} else {
		m.logger.Info("dispatched magic packet",
			zap.String("remote", r.RemoteAddr),
			zap.String("host", r.Host),
			zap.String("uri", r.RequestURI),
			zap.String("broadcast", m.BroadcastAddress),
			zap.String("mac", m.MAC),
		)
		_, err := m.broadcastSocket.Write(m.magicPacket)
		if err != nil {
			return err
		}
		time.AfterFunc(10*time.Minute, func() {
			_, _ = m.pool.Delete(m.key)
		})
	}
	return next.ServeHTTP(w, r)
}

// 函数用于解析 Caddyfile 配置文件中的参数。
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()

		switch len(args) {
		case 1:
			m.MAC, m.BroadcastAddress = args[0], "255.255.255.255:9"
		case 2:
			m.MAC, m.BroadcastAddress = args[0], args[1]
		default:
			return d.Err("unexpected number of arguments")
		}
	}
	return nil
}

// 建立一个魔术包
func BuildMagicPacket(mac net.HardwareAddr) []byte {
	macBytes := []byte(mac)
	mp := make([]byte, 6)
	for idx := 0; idx < 6; idx++ {
		mp[idx] = 0xFF
	}
	for idx := 0; idx < 16; idx++ {
		// TODO: refactor to pre-allocate the packet memory
		mp = append(mp, macBytes...)
	}
	return mp
}

// 函数关闭准备好的广播 socket。
func (m *Middleware) Cleanup() error {
	return m.broadcastSocket.Close()
}

// 用于将 Caddyfile 配置文件中的参数解析为 Middleware 结构体
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.CleanerUpper          = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
