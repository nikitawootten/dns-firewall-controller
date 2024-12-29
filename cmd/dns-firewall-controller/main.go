package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/coreos/go-iptables/iptables"
	"github.com/nikitawootten/dns-firewall-controller/controller"
	"github.com/nikitawootten/dns-firewall-controller/firewall"
	"github.com/urfave/cli/v2"
)

func createListener(c *cli.Context) (net.Listener, error) {
	address := c.String("address")

	log.Printf("Starting DNS Firewall Controller on address '%v'", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address '%v': %w", address, err)
	}

	return listener, nil
}

func main() {
	app := &cli.App{
		Name:        "dns-firewall-controller",
		Description: "Open firewall rules based on DNS responses",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "address",
				Usage: "Specify the address to listen on",
				Value: ":6000",
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "mock",
				Category: "backend",
				Usage:    "Start the controller with a mock (no-op) backend",
				Action: func(c *cli.Context) error {
					listener, err := createListener(c)
					if err != nil {
						return err
					}

					backend := firewall.NewMockBackend()
					return controller.Start(backend, listener)
				},
			},
			{
				Name:     "iptables",
				Category: "backend",
				Usage:    "Start the controller with an iptables backend",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "table",
						Usage: "Specify the iptables table to use",
						Value: "filter",
					},
					&cli.StringFlag{
						Name:  "chain",
						Usage: "Specify the iptables chain to use",
						Value: "INPUT",
					},
				},
				Action: func(c *cli.Context) error {
					table := c.String("table")
					chain := c.String("chain")

					listener, err := createListener(c)
					if err != nil {
						return err
					}

					iptables, err := iptables.New()
					if err != nil {
						return err
					}

					config := firewall.IPTablesFirewallBackend{
						Table:    table,
						Chain:    chain,
						IPTables: iptables,
					}

					backend := firewall.NewIPTablesBackend(&config)
					return controller.Start(backend, listener)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
