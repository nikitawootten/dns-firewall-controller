.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test ./...

COREDNS_PORT := 5300
COREDNS_COREFILE := support/TestCorefile

.PHONY: run-coredns
run-coredns:
	nix run nixpkgs#coredns -- -conf $(COREDNS_COREFILE) -p $(COREDNS_PORT)

.PHONY: run-dns-firewall-controller
run-dns-firewall-controller:
	nix run .#dns-firewall-controller

.PHONY: send-dns-request
send-dns-request:
	dig @127.0.0.1 -p $(COREDNS_PORT) lvh.me
