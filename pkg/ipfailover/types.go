package ipfailover

import (
	"github.com/projectatomic/appinfra-next/pkg/cmd/util/variable"
)

const (
	// DefaultName is the default IP Failover resource name.
	DefaultName = "ipfailover"

	// DefaultType is the default IP Failover type.
	DefaultType = "keepalived"

	// DefaultServicePort is the default service port.
	DefaultServicePort = 1985

	// DefaultWatchPort is the default IP Failover watched port number.
	DefaultWatchPort = 80

	// DefaultSelector is the default resource selector.
	DefaultSelector = "ipfailover=<name>"

	// DefaultInterface is the default network interface.
	DefaultInterface = "eth0"
)

// IPFailoverConfigCmdOptions are options supported by the IP Failover admin command.
type IPFailoverConfigCmdOptions struct {
	Type           string
	ImageTemplate  variable.ImageTemplate
	Credentials    string
	ServicePort    int
	Selector       string
	Create         bool
	ServiceAccount string

	//  Failover options.
	VirtualIPs       string
	NetworkInterface string
	WatchPort        int
	Replicas         int
}
