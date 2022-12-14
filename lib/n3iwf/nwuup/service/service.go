package service

import (
	"errors"
	"net"

	"my5G-RANTester/lib/n3iwf/context"

	log "github.com/sirupsen/logrus"

	gtpv1 "github.com/wmnsk/go-gtp/v1"
	"golang.org/x/net/ipv4"
)

// Run bind and listen raw socket on N3IWF NWu interface
// with UP_IP_ADDRESS, catching GRE encapsulated packets and forward
// to N3 interface.
func Run() error {
	// Local IPSec address
	n3iwfSelf := context.N3IWFSelf()
	listenAddr := n3iwfSelf.IPSecGatewayAddress

	// Setup raw socket
	// This raw socket will only capture GRE encapsulated packet
	connection, err := net.ListenPacket("ip4:gre", listenAddr)
	if err != nil {
		log.Errorf("Error setting listen socket on %s: %+v", listenAddr, err)
		return errors.New("ListenPacket failed")
	}
	rawSocket, err := ipv4.NewRawConn(connection)
	if err != nil {
		log.Errorf("Error opening raw socket on %s: %+v", listenAddr, err)
		return errors.New("NewRawConn failed")
	}

	n3iwfSelf.NWuRawSocket = rawSocket
	go listenAndServe(rawSocket)

	return nil
}

// listenAndServe read from socket and call forward() to
// forward packet.
func listenAndServe(rawSocket *ipv4.RawConn) {
	defer func() {
		err := rawSocket.Close()
		if err != nil {
			log.Errorf("Error closing raw socket: %+v", err)
		}
	}()

	buffer := make([]byte, 65535)

	for {
		ipHeader, ipPayload, _, err := rawSocket.ReadFrom(buffer)
		log.Tracef("Read %d bytes", len(ipPayload))
		if err != nil {
			log.Errorf("Error read from raw socket: %+v", err)
			return
		}

		forwardData := make([]byte, len(ipPayload[4:]))
		copy(forwardData, ipPayload[4:])

		go forward(ipHeader.Src.String(), forwardData)
	}

}

// forward forwards user plane packets from NWu to UPF
// with GTP header encapsulated
func forward(ueInnerIP string, packet []byte) {
	// Find UE information
	self := context.N3IWFSelf()
	ue, ok := self.AllocatedUEIPAddressLoad(ueInnerIP)
	if !ok {
		log.Error("UE context not found")
		return
	}

	var pduSession *context.PDUSession

	for _, pduSession = range ue.PduSessionList {
		break
	}

	if pduSession == nil {
		log.Error("This UE doesn't have any available PDU session")
		return
	}

	gtpConnection := pduSession.GTPConnection

	userPlaneConnection := gtpConnection.UserPlaneConnection

	n, err := userPlaneConnection.WriteToGTP(gtpConnection.OutgoingTEID, packet, gtpConnection.UPFUDPAddr)
	if err != nil {
		log.Errorf("Write to UPF failed: %+v", err)
		if err == gtpv1.ErrConnNotOpened {
			log.Error("The connection has been closed")
			// TODO: Release the GTP resource
		}
		return
	} else {
		log.Trace("Forward NWu -> N3")
		log.Tracef("Wrote %d bytes", n)
		return
	}
}
