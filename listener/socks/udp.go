package socks

import (
	"net"

	"github.com/MysticalDevil/clash/adapter/inbound"
	"github.com/MysticalDevil/clash/common/pool"
	"github.com/MysticalDevil/clash/common/sockopt"
	C "github.com/MysticalDevil/clash/constant"
	"github.com/MysticalDevil/clash/log"
	"github.com/MysticalDevil/clash/transport/socks5"
)

type UDPListener struct {
	packetConn net.PacketConn
	addr       string
	closed     bool
}

// RawAddress implements C.Listener
func (l *UDPListener) RawAddress() string {
	return l.addr
}

// Address implements C.Listener
func (l *UDPListener) Address() string {
	return l.packetConn.LocalAddr().String()
}

// Close implements C.Listener
func (l *UDPListener) Close() error {
	l.closed = true
	return l.packetConn.Close()
}

func NewUDP(addr string, in chan<- *inbound.PacketAdapter) (C.Listener, error) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}

	if err := sockopt.UDPReuseaddr(l.(*net.UDPConn)); err != nil {
		log.Warnln("Failed to Reuse UDP Address: %s", err)
	}

	sl := &UDPListener{
		packetConn: l,
		addr:       addr,
	}
	go func() {
		for {
			buf := pool.Get(pool.UDPBufferSize)
			n, remoteAddr, err := l.ReadFrom(buf)
			if err != nil {
				pool.Put(buf)
				if sl.closed {
					break
				}
				continue
			}
			handleSocksUDP(l, in, buf[:n], remoteAddr)
		}
	}()

	return sl, nil
}

func handleSocksUDP(pc net.PacketConn, in chan<- *inbound.PacketAdapter, buf []byte, addr net.Addr) {
	target, payload, err := socks5.DecodeUDPPacket(buf)
	if err != nil {
		// Unresolved UDP packet, return buffer to the pool
		pool.Put(buf)
		return
	}
	packet := &packet{
		pc:      pc,
		rAddr:   addr,
		payload: payload,
		bufRef:  buf,
	}
	select {
	case in <- inbound.NewPacket(target, pc.LocalAddr(), packet, C.SOCKS5):
	default:
	}
}
