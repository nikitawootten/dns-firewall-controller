package firewall

import (
	"net"
)

type Firewall interface {
	AddRule(client net.IP, allowed []net.IP) error
	RemoveRule(client net.IP) error
}
