package cmd

import (
	"log"
	"net"

	firewall_controller "github.com/nikitawootten/dns-firewall-controller/src/firewall_controller/server"
	pb "github.com/nikitawootten/dns-firewall-controller/src/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the firewall controller server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			log.Fatal(err)
		}

		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to bind to address: %v", err)
		}

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		firewallController := firewall_controller.NewFirewallControllerServer()
		pb.RegisterFirewallControllerServer(grpcServer, firewallController)
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("address", defaultAddress, "Address to bind to")
}
