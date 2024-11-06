# discovery-go

This package provide the cluster discovery engine for [Groupcache](https://github.com/groupcache/groupcache-go)

# Features

- Discovery API to implement custom nodes discovery provider. See: [discovery](./discovery/provider.go)
- Comes bundled with some discovery providers that can help you hit the ground running:
    - [kubernetes](https://kubernetes.io/docs/home/) [api integration](./discovery/kubernetes) is fully functional
    - [nats](https://nats.io/) [integration](./discovery/nats) is fully functional
    - [static](./discovery/static) is fully functional and for demo purpose
    - [dns](./discovery/dnssd) is fully functional

# Built-in providers

## DNS Provider Setup

This provider performs nodes discovery based upon the domain name provided. This is very useful when doing local development
using docker.

To use the DNS discovery provider one needs to provide the following:

- `DomainName`: the domain name
- `IPv6`: it states whether to lookup for IPv6 addresses.

```go
package main

import "github.com/groupcache/discovery-go/discovery/dnssd"

const domainName = "accounts"

// define the discovery options
config := dnssd.Config{
    dnssd.DomainName: domainName,
    dnssd.IPv6:       false,
}
// instantiate the dnssd discovery provider
disco := dnssd.NewDiscovery(&config)

// pass the service discovery when enabling cluster mode in the actor system
```