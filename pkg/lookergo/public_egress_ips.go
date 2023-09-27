package lookergo

import (
	"context"
)

type EgressIpAddresses struct {
	EgressIpAddresses *[]string `json:"egress_ip_addresses,omitempty"` // Egress IP addresses
}

const publicEgressIpsBasePath = "4.0/public_egress_ip_addresses"

type PublicEgressIpsResource interface {
	Get(ctx context.Context) (*EgressIpAddresses, *Response, error)
}

type PublicEgressIpsResourceOp struct {
	client *Client
}

var _PublicEgressIpsResource = &PublicEgressIpsResourceOp{}

func (s PublicEgressIpsResourceOp) Get(ctx context.Context) (*EgressIpAddresses, *Response, error) {
	return doGet(ctx, s.client, publicEgressIpsBasePath, new(EgressIpAddresses))
}
