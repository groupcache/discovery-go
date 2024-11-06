# discovery-go

This package provide the cluster discovery engine for [Groupcache](https://github.com/groupcache/groupcache-go)

# Features

- Discovery API to implement custom nodes discovery provider. See: [discovery](./discovery/provider.go)
- Comes bundled with some discovery providers that can help you hit the ground running:
    - [kubernetes](https://kubernetes.io/docs/home/) [api integration](./discovery/kubernetes) is fully functional
    - [nats](https://nats.io/) [integration](./discovery/nats) is fully functional
    - [static](./discovery/static) is fully functional and for demo purpose
    - [dns](./discovery/dnssd) is fully functional