package network

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

var privateCIDRs []*net.IPNet

func init() {
	cidrs := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	for _, cidr := range cidrs {
		_, ipNet, _ := net.ParseCIDR(cidr)
		privateCIDRs = append(privateCIDRs, ipNet)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, cidr := range privateCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// NewSafeClient returns an *http.Client with a 30-second timeout.
// When allowPrivate is false, connections to private/internal IP addresses
// are blocked to prevent SSRF attacks. DNS resolution results are used
// directly for dialing to prevent DNS rebinding.
func NewSafeClient(allowPrivate bool) *http.Client {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			if len(ips) == 0 {
				return nil, fmt.Errorf("no IP addresses found for %s", host)
			}

			ip := ips[0].IP
			if !allowPrivate && isPrivateIP(ip) {
				return nil, fmt.Errorf("access to private IP address %s (%s) is blocked", ip, host)
			}

			// Connect using resolved IP directly (prevents DNS rebinding)
			resolvedAddr := net.JoinHostPort(ip.String(), port)
			return dialer.DialContext(ctx, network, resolvedAddr)
		},
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}
