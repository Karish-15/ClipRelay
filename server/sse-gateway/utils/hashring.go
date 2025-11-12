package utils

import (
	"os"
	"strings"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash/v2"
)

type Gateway struct {
	Addr string
}

func (g Gateway) String() string {
	return g.Addr
}

type Hasher struct{}

func (h Hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

var (
	HashRing    *consistent.Consistent
	GatewayAddr string
)

func InitConsistentHashingRing() {
	GatewayAddr = os.Getenv("GATEWAY_ADDR")
	if GatewayAddr == "" {
		panic("GATEWAY_ADDR not set")
	}

	gatewaysEnv := os.Getenv("GATEWAYS")
	if gatewaysEnv == "" {
		panic("GATEWAYS not set")
	}

	addrs := strings.Split(gatewaysEnv, ",")
	members := make([]consistent.Member, len(addrs))
	for i, a := range addrs {
		members[i] = Gateway{Addr: strings.TrimSpace(a)}
	}

	cfg := consistent.Config{
		PartitionCount:    271,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            Hasher{},
	}

	HashRing = consistent.New(members, cfg)
}
