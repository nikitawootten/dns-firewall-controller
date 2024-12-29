package controller

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikitawootten/dns-firewall-controller/dns"
	"github.com/nikitawootten/dns-firewall-controller/firewall"
)

func Start(backend firewall.FirewallBackend, listener net.Listener) error {
	receiver := dns.NewDNSTapReceiver(listener, func(response *dns.DNSResponse) {
		rules := firewall.FirewallRulesFromDNSResponse(response)
		for _, rule := range rules {
			backend.AddRule(rule)
		}
	})

	if err := backend.Start(); err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		backend.Stop()
		listener.Close()
		os.Exit(0)
	}()

	receiver.Start()
	return nil
}
