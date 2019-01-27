package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/vishvananda/netlink"
)

var (
	cloneBr = "cloneBr"
	ip      = "10.10.10.1/24"
)

//CreateBr a bridge
func CreateBr() error {
	if _, err := net.InterfaceByName(cloneBr); err == nil {
		return nil
	}
	// create *netlink.Bridge object
	la := netlink.LinkAttrs{Name: cloneBr}
	br := &netlink.Bridge{LinkAttrs: la}
	//br := &netlink.Bridge{la}
	if err := netlink.LinkAdd(br); err != nil {
		log.Println("Error whiel creating BR")
		return err
	}
	//parse the Ip addr
	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		log.Println("Ip parsing error")
		return err
	}
	//next add the Ip to the device
	if err := netlink.AddrAdd(br, addr); err != nil {
		log.Println("Error adding IP")
		return err
	}
	//Once the br and ip is assoc, now bring up the device
	if err := netlink.LinkSetUp(br); err != nil {
		log.Println("Error bringing up the br")
		return err
	}
	return nil
}

//CreateVethPair a Veth pair with the PID ns and the Host-ns(bridge)
//The following function does the following:

//ip link add tap1 type veth peer name container_name-veth

//now move the container_name-veth into container.
//ip link set container_name-veth netns contianer_name

//bring up both veth pairs
//ip link set tap1 up
//ip netns exec container_name ip link set dev container_name-veth up

func CreateVethPair(pid int) error {
	// get bridge to set as master for one side of veth-pair
	br, err := netlink.LinkByName(cloneBr)
	if err != nil {
		return err
	}
	x1, x2 := rand.Intn(100000), rand.Intn(100000)
	hostInt := fmt.Sprintf("veth-%d", x1)
	peerInt := fmt.Sprintf("veth-%d", x2)

	//now create a linkAttrs for veth
	v := netlink.NewLinkAttrs()
	v.Name = hostInt
	//masterindex is always the index of the br, in our case always br
	v.MasterIndex = br.Attrs().Index
	vPair := &netlink.Veth{LinkAttrs: v, PeerName: peerInt}
	if err := netlink.LinkAdd(vPair); err != nil {
		log.Println("Error whiel creating BR")
		return err
	}
	//now get peer and put it in the peer NS.
	peer, err := netlink.LinkByName(peerInt)
	if err != nil {
		log.Println("Err while fetching peer")
		return err
	}
	if err = netlink.LinkSetNsPid(peer, pid); err != nil {
		log.Println("Error while attach veth peer to ns/pid")
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: netcreate <pid>")
		os.Exit(1)
	}
	pid, _ := strconv.Atoi(os.Args[1])

	if err := CreateBr(); err != nil {
		log.Fatalln("Error while creating br", err)
	}
	if err := CreateVethPair(pid); err != nil {
		log.Fatalln("Error while creating Veth Pair", err)
	}
}
