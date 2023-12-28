{ pkgs, ... }@inputs:
let
  coredns-with-plugin = import ./buildCoreDnsWithPlugin.nix inputs;
  firewall-controller = import ./buildFirewallController.nix inputs;
in
pkgs.nixosTest {
  name = "firewall-controller-integration-test";
  # In this scenario, we have 3 machines in the "internal" network, and 1 "external" machine.
  # The internet network consists of a router and 2 clients ("good" and "bad").
  # The difference between the clients is that the upstream client does not respect the preferred DNS server.
  # The external network consists of a web server, and it used to verify the firewall rules.
  nodes = let
    base = { ... }: {
      system.stateVersion = "23.11";
    };
    clientBase = { nodes, ... }: {
      imports = [ base ];
      virtualisation.vlans = [ 1 ];
      networking.defaultGateway = (pkgs.lib.head nodes.router.networking.interfaces.eth2.ipv4.addresses).address;
    };
  in {
    router = { nodes, ... }: {
      imports = [ base ];
      # Configure NAT
      virtualisation.vlans = [ 2 1 ];
      networking.firewall.enable = true;
      networking.firewall.filterForward = true;
      networking.nftables.enable = true;
      networking.nat.enable = true;
      networking.nat.internalIPs = [ "192.168.1.0/24" ];
      networking.nat.externalInterface = "eth1";

      # Configure DNS server
      services.coredns = {
        enable = true;
        package = coredns-with-plugin;
        config = let
          hosts-file = ''
            ${(pkgs.lib.head nodes.server.networking.interfaces.eth1.ipv4.addresses).address} server
          '';
        in ''
          . {
            squawker 127.0.0.1:8080
            hosts ${pkgs.writeText "test-hosts" hosts-file}
          }
        '';
      };

      # Configure firewall controller
      systemd.services.firewall-controller = {
        wantedBy = [ "multi-user.target" ];
        after = [ "network.target" ];
        description = "Start the firewall controller";
        serviceConfig = {
          Type = "exec";
          User = "firewall";
          ExecStart = "${firewall-controller}/bin/firewall-controller server";
        };
      };

      users.users.firewall = {
        isSystemUser = true;
        group = "wheel";
      };
    };

    goodClient = { config, ... }: {
      imports = [ clientBase ];
      # The good client respects the preferred DNS server.
      networking.nameservers = [ config.networking.defaultGateway.address ];
    };

    badClient = { ... }: {
      imports = [ clientBase ];
    };

    server = { ... }: {
      imports = [ base ];
      virtualisation.vlans = [ 2 ];
    };
  };

  testScript = ''
    router.start()
    server.start()
    goodClient.start()
    badClient.start()

    router.wait_for_unit("firewall-controller")
    server.wait_for_unit("network.target")
    goodClient.wait_for_unit("network.target")
    badClient.wait_for_unit("network.target")

    goodClient.succeed("ping -c 1 server")
    badClient.fail("ping -c 1 server")
  '';
}
