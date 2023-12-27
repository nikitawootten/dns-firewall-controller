.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test ./...

.PHONY: codegen
codegen:
	go generate ./...

COREDNS_PORT := 5300

.PHONY: run-coredns
run-coredns:
	nix run .#coredns -- -conf support/TestCorefile -p $(COREDNS_PORT)

.PHONY: run-firewall-controller
run-firewall-controller:
	nix run .#firewall-controller -- server --address :8080

.PHONY: send-dns-request
send-dns-request:
	dig @127.0.0.1 -p $(COREDNS_PORT) lvh.me
