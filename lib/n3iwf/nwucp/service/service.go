package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"

	"my5G-RANTester/lib/n3iwf/context"
	"my5G-RANTester/lib/n3iwf/ngap/message"

	log "github.com/sirupsen/logrus"
)

// Run setup N3IWF NAS for UE to forward NAS message
// to AMF
func Run() error {
	// N3IWF context
	n3iwfSelf := context.N3IWFSelf()
	tcpAddr := fmt.Sprintf("%s:%d", n3iwfSelf.IPSecGatewayAddress, n3iwfSelf.TCPPort)

	tcpListener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Errorf("Listen TCP address failed: %+v", err)
		return errors.New("Listen failed")
	}

	log.Tracef("Successfully listen %+v", tcpAddr)

	go listenAndServe(tcpListener)

	return nil
}

// listenAndServe handle TCP listener and accept incoming
// requests. It also stores accepted connection into UE
// context, and finally, call serveConn() to serve the messages
// received from the connection.
func listenAndServe(tcpListener net.Listener) {
	defer func() {
		err := tcpListener.Close()
		if err != nil {
			log.Errorf("Error closing tcpListener: %+v", err)
		}
	}()

	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			log.Error("TCP server accept failed. Close the listener...")
			return
		}

		log.Tracef("Accepted one UE from %+v", connection.RemoteAddr())

		// Find UE context and store this connection in to it, then check if
		// there is any cached NAS message for this UE. If yes, send to it.
		n3iwfSelf := context.N3IWFSelf()

		ueIP := strings.Split(connection.RemoteAddr().String(), ":")[0]
		ue, ok := n3iwfSelf.AllocatedUEIPAddressLoad(ueIP)
		if !ok {
			log.Errorf("UE context not found for peer %+v", ueIP)
			continue
		}

		// Store connection
		ue.TCPConnection = connection

		if ue.TemporaryCachedNASMessage != nil {
			// Send to UE
			if n, err := connection.Write(ue.TemporaryCachedNASMessage); err != nil {
				log.Errorf("Writing via IPSec signalling SA failed: %+v", err)
			} else {
				log.Trace("Forward NWu <- N2")
				log.Tracef("Wrote %d bytes", n)
			}
			// Clean the cached message
			ue.TemporaryCachedNASMessage = nil
		}

		go serveConn(ue, connection)
	}
}

// serveConn handle accepted TCP connection. It reads NAS packets
// from the connection and call forward() to forward NAS messages
// to AMF
func serveConn(ue *context.N3IWFUe, connection net.Conn) {
	defer func() {
		err := connection.Close()
		if err != nil {
			log.Errorf("Error closing connection: %+v", err)
		}
	}()

	data := make([]byte, 65535)
	for {
		n, err := connection.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				log.Warn("Connection close by peer")
				ue.TCPConnection = nil
				return
			} else {
				log.Errorf("Read TCP connection failed: %+v", err)
			}
		}
		log.Tracef("Get NAS PDU from UE:\nNAS length: %d\nNAS content:\n%s", n, hex.Dump(data[:n]))

		forwardData := make([]byte, n)
		copy(forwardData, data[:n])

		go forward(ue, forwardData)
	}
}

// forward forwards NAS messages sent from UE to the
// associated AMF
func forward(ue *context.N3IWFUe, packet []byte) {
	log.Trace("Forward NWu -> N2")
	message.SendUplinkNASTransport(ue.AMF, ue, packet)
}
