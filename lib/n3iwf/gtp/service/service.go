package service

import (
	"context"
	"errors"
	"net"
	"syscall"

	n3iwf_context "my5G-RANTester/lib/n3iwf/context"

	log "github.com/sirupsen/logrus"

	gtpv1 "github.com/wmnsk/go-gtp/v1"
	"golang.org/x/net/ipv4"
)

// SetupGTPTunnelWithUPF set up GTP connection with UPF
// return *gtpv1.UPlaneConn, net.Addr and error
func SetupGTPTunnelWithUPF(upfIPAddr string) (*gtpv1.UPlaneConn, net.Addr, error) {
	n3iwfSelf := n3iwf_context.N3IWFSelf()

	// Set up GTP connection
	upfUDPAddr := upfIPAddr + ":2152"

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", upfUDPAddr)
	if err != nil {
		log.Errorf("Resolve UDP address %s failed: %+v", upfUDPAddr, err)
		return nil, nil, errors.New("Resolve Address Failed")
	}

	n3iwfUDPAddr := n3iwfSelf.GTPBindAddress + ":2152"

	localUDPAddr, err := net.ResolveUDPAddr("udp", n3iwfUDPAddr)
	if err != nil {
		log.Errorf("Resolve UDP address %s failed: %+v", n3iwfUDPAddr, err)
		return nil, nil, errors.New("Resolve Address Failed")
	}

	context := context.TODO()

	// Dial to UPF
	userPlaneConnection, err := gtpv1.DialUPlane(context, localUDPAddr, remoteUDPAddr)
	if err != nil {
		log.Errorf("Dial to UPF failed: %+v", err)
		return nil, nil, errors.New("Dial failed")
	}

	return userPlaneConnection, remoteUDPAddr, nil

}

// ListenAndServe binds and listens raw socket on N3IWF N3 interface,
// catching GTP packets and send it to NWu interface
func ListenAndServe(userPlaneConnection *gtpv1.UPlaneConn) error {
	go listenGTP(userPlaneConnection)
	return nil
}

// listenGTP handle the gtpv1 UPlane connection. It reads packets(without
// GTP header) from the connection and call forward() to forward user data
// to NWu interface.
func listenGTP(userPlaneConnection *gtpv1.UPlaneConn) {
	defer func() {
		err := userPlaneConnection.Close()
		if err != nil {
			log.Errorf("userPlaneConnection Close failed: %+v", err)
		}
	}()

	payload := make([]byte, 65535)

	for {
		n, _, teid, err := userPlaneConnection.ReadFromGTP(payload)
		log.Tracef("Read %d bytes", n)
		if err != nil {
			log.Errorf("Read from GTP failed: %+v", err)
			return
		}

		forwardData := make([]byte, n)
		copy(forwardData, payload[:n])

		go forward(teid, forwardData)
	}

}

// forward forwards user plane packets from N3 to UE,
// with GRE header and new IP header encapsulated
func forward(ueTEID uint32, packet []byte) {
	// This is the IP header template for packets with GRE header encapsulated.
	// The remaining mandatory fields are Dst and TotalLen, which specified
	// the destination IP address and the packet total length.

	// Find UE information
	self := n3iwf_context.N3IWFSelf()
	ue, ok := self.AllocatedUETEIDLoad(ueTEID)
	if !ok {
		log.Error("UE context not found")
		return
	}

	ipHeader := &ipv4.Header{
		Version:  4,
		Len:      20,
		TOS:      0,
		Flags:    ipv4.DontFragment,
		FragOff:  0,
		TTL:      64,
		Protocol: syscall.IPPROTO_GRE,
	}

	// GRE header
	greHeader := []byte{0, 0, 8, 0}

	// UE IP
	ueInnerIP := net.ParseIP(ue.IPSecInnerIP)

	greEncapsulatedPacket := append(greHeader, packet...)
	packetTotalLength := 20 + len(greEncapsulatedPacket)

	ipHeader.Dst = ueInnerIP
	ipHeader.TotalLen = packetTotalLength

	n3iwfSelf := n3iwf_context.N3IWFSelf()
	rawSocket := n3iwfSelf.NWuRawSocket

	// Send to UE
	if err := rawSocket.WriteTo(ipHeader, greEncapsulatedPacket, nil); err != nil {
		log.Errorf("Write to raw socket failed: %+v", err)
		return
	} else {
		log.Trace("Forward NWu <- N3")
		log.Tracef("Wrote %d bytes", packetTotalLength)
	}
}
