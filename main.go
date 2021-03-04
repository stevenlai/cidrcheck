package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/apparentlymart/go-cidr/cidr"
)

func main() {
	networkArg := flag.String("cidr", "", "CIDR of network")
	ipArg := flag.String("ip", "", "IP to verify")

	flag.Parse()

	if *networkArg == "" {
		fmt.Println(fmt.Errorf("Must provide CIDR with -cidr"))
		os.Exit(2)
	}

	if *ipArg == "" {
		fmt.Println(fmt.Errorf("Must provide IP with -ip"))
		os.Exit(2)
	}

	_, network, err := net.ParseCIDR(*networkArg)
	if err != nil {
		fmt.Println(fmt.Errorf("Invalid CIDR"))
		os.Exit(2)
	}
	_, testNetwork, err := net.ParseCIDR(fmt.Sprintf("%s/32", *ipArg))

	if err != nil {
		fmt.Println(fmt.Errorf("Invalid IP given"))
		os.Exit(2)
	}

	fmt.Printf("%v in %v: %+v\n", *ipArg, *networkArg, IsInCIDR(network, testNetwork))
}

// IsInCIDR checks if a given IP is in a given CIDR range
func IsInCIDR(n *net.IPNet, testNetwork *net.IPNet) bool {
	firstIP, lastIP := cidr.AddressRange(n)
	currentIP := firstIP

	// If only 1 address in range and test IP is the same, then it's true.
	if isSingleHostRange(n) && equalIP(n.IP, testNetwork.IP) {
		return true
	} else {
		for !equalIP(currentIP, lastIP) {
			// We use the testNetwork.IP here instead of the IP because ParseCIDR returns
			// the IP as 16 bytes whereas the IP from *net.IPNet is 4 bytes.
			if equalIP(currentIP, testNetwork.IP) {
				return true
			}
			currentIP = cidr.Inc(currentIP)
		}
	}
	return false
}

// isSingleHostRange checks if the mask of the given CIDR is 32 and therefore, a single host
func isSingleHostRange(n *net.IPNet) bool {
	prefixLen, _ := n.Mask.Size()
	if prefixLen == 32 {
		return true
	}
	return false
}

// equalIP checks if two net.IP are equal
func equalIP(a, b net.IP) bool {
	return bytes.Equal(a, b)
}
