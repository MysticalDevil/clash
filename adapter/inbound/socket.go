package inbound

import (
	"net"
	"net/netip"

	C "github.com/MysticalDevil/clash/constant"
	"github.com/MysticalDevil/clash/context"
	"github.com/MysticalDevil/clash/transport/socks5"
)

// NewSocket receive TCP inbound and return ConnContext
func NewSocket(target socks5.Addr, conn net.Conn, source C.Type) *context.ConnContext {
	metadata := parseSocksAddr(target)
	metadata.NetWork = C.TCP
	metadata.Type = source
	if ip, port, err := parseAddr(conn.RemoteAddr()); err == nil {
		metadata.SrcIP = ip
		metadata.SrcPort = C.Port(port)
	}
	if addrPort, err := netip.ParseAddrPort(conn.LocalAddr().String()); err == nil {
		metadata.OriginDst = addrPort
	}
	return context.NewConnContext(conn, metadata)
}
