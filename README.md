# DNS Firewall Controller

## Background

Suppose we have a simple NAT network with a few devices connected to a gateway that provides access to the broader internet.
For simplicity, let us imagine the DNS server as a device on that network.
All clients connected to the network are expected to use the DNS server.

```mermaid
graph
  Internet

  subgraph NAT Network
    Gateway
    DNS_Server[DNS Server]
    Client_1[Client 1]
    Client_2[Client 2]

    DNS_Server-->Gateway
    Client_1-->Gateway
    Client_2-->Gateway
  end

  Gateway-->Internet
```

This topology is similar to that of most home networks, except that the router often includes its own DNS server.
Here you could imagine clients to be computers, IoT devices, phones, or any other networked consumer device.

### Benefits

Networks with a controlled DNS server have a few benefits:

1. Through a technique called a **DNS Sinkhole**, DNS servers can choose to block certain domains.
   Authoritative block-lists or even allow-lists can be used to prevent users from accessing malicious domains, or domains that are blocked for policy reasons.
1. DNS servers can be configured to keep logs of client connections, which could be important for security or compliance reasons.
1. DNS servers can be configured with additional security measures, such as DNSSEC.

### Drawbacks

The main drawback of this model is that devices connected to the network are under no obligation to honor the DNS server assigned to them via DHCP.
Firewalls can be used to block other DNS traffic, but there is nothing stopping a malicious client from resolving blocked domains via a side channel DNS server, such as through DNS over HTTPS.
Attackers could also sidestep the problem of DNS resolution completely using hardcoded IP addresses.

## Solution

One solution to this problem is to force clients of a network to announce their intention to connect to an address through the DNS server before giving them any permissions to access external addresses.
Honest clients respecting the network's preferred DNS server will be unaffected.

### Honest Client Makes a Successful Connection

```mermaid
sequenceDiagram
  autonumber

  participant Client
  participant DNS_Server as DNS Server
  participant Gateway

  Client->>+DNS_Server: resolve domain.example
  DNS_Server->>Gateway: allow Client to access X.X.X.X for 1000ms
  Note right of Gateway: Firewall opens X.X.X.X to the Client for 1000ms
  DNS_Server->>-Client: domain.example is X.X.X.X for 1000ms
  Client->>Gateway: access X.X.X.X

  Note right of Client: Client successfully makes a connection
```

### Malicious Client Attempts to use a Rogue DNS Server

If a **client** resolves an address through a rogue DNS server or side channel, they will still be denied by the **gateway** because the **preferred DNS server** did not report the request.

```mermaid
sequenceDiagram
  autonumber

  participant Rogue_DNS_Server as Rogue DNS Server
  participant Client
  participant Gateway as Gateway

  Client->>Rogue_DNS_Server: resolve domain.example
  Rogue_DNS_Server->>Client: domain.example is X.X.X.X for 1000ms
  Client-XGateway: access X.X.X.X

  Note right of Client: Client is denied
```

### Honest Client Attempts to Resolve a Block-listed Domain

```mermaid
sequenceDiagram
  autonumber

  participant Client
  participant DNS_Server as DNS Server
  participant Gateway

  Client->>+DNS_Server: resolve malicious-domain.example
  DNS_Server->>-Client: malicious-domain.example does not exist

  Note right of Client: Client did not receive a resolved address
```

