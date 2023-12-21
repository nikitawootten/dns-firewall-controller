package cmd

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/nikitawootten/dns-firewall-controller/src/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run the firewall controller client",
	Args:  cobra.NoArgs,
}

func initClient(cmd *cobra.Command) (pb.FirewallControllerClient, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}

	return pb.NewFirewallControllerClient(conn), err
}

var clientWriteClientPolicyCmd = &cobra.Command{
	Use:   "write-client-policy",
	Short: "Send a client policy to the server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		grpc_client, err := initClient(cmd)
		if err != nil {
			log.Fatalf("failed to connect to server: %v", err)
		}

		// Parse client IP

		client_raw, err := cmd.Flags().GetString("client")
		if err != nil {
			log.Fatal(err)
		}

		client := net.ParseIP(client_raw)
		proto_client, err := pb.ToProtoIpAddress(client)
		if err != nil {
			log.Fatal(err)
		}

		// Parse allowed IPs

		allowed_raw, err := cmd.Flags().GetStringArray("allow")
		if err != nil {
			log.Fatal(err)
		}

		allowed := make([]net.IP, len(allowed_raw))
		for i, ip_raw := range allowed_raw {
			allowed[i] = net.ParseIP(ip_raw)
		}

		proto_allowed, err := pb.ToProtoIpAddresses(allowed)
		if err != nil {
			log.Fatal(err)
		}

		// Parse duration

		for_raw, err := cmd.Flags().GetString("for")
		if err != nil {
			log.Fatal(err)
		}

		for_duration, err := time.ParseDuration(for_raw)
		if err != nil {
			log.Fatal(err)
		}

		allow_until := time.Now().Add(for_duration)

		grpc_client.WriteClientPolicy(context.Background(), &pb.ClientPolicy{
			Client:     proto_client,
			AllowedIps: proto_allowed,
			AllowUntil: timestamppb.New(allow_until),
		})
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().String("address", defaultAddress, "Address to send requests to")
	// TODO: add TLS support

	clientCmd.AddCommand(clientWriteClientPolicyCmd)
	clientWriteClientPolicyCmd.Flags().String("client", "", "IP address of client")
	clientWriteClientPolicyCmd.MarkFlagRequired("client")

	clientWriteClientPolicyCmd.Flags().StringArray("allow", []string{}, "IP addresses to allow")

	clientWriteClientPolicyCmd.Flags().String("for", "", "Duration of policy, e.g. 100ms")
	clientWriteClientPolicyCmd.MarkFlagRequired("for")
}
