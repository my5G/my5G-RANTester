package ike

import (
	"my5G-RANTester/lib/n3iwf/ike/handler"
	ike_message "my5G-RANTester/lib/n3iwf/ike/message"
	"net"

	log "github.com/sirupsen/logrus"
)

func Dispatch(udpConn *net.UDPConn, localAddr, remoteAddr *net.UDPAddr, msg []byte) {
	// As specified in RFC 7296 section 3.1, the IKE message send from/to UDP port 4500
	// should prepend a 4 bytes zero
	if localAddr.Port == 4500 {
		for i := 0; i < 4; i++ {
			if msg[i] != 0 {
				log.Warn(
					"[IKE] Received an IKE packet that does not prepend 4 bytes zero from UDP port 4500," +
						" this packet may be the UDP encapsulated ESP. The packet will not be handled.")
				return
			}
		}
		msg = msg[4:]
	}

	ikeMessage, err := ike_message.Decode(msg)
	if err != nil {
		log.Error(err)
		return
	}

	switch ikeMessage.ExchangeType {
	case ike_message.IKE_SA_INIT:
		handler.HandleIKESAINIT(udpConn, localAddr, remoteAddr, ikeMessage)
	case ike_message.IKE_AUTH:
		handler.HandleIKEAUTH(udpConn, localAddr, remoteAddr, ikeMessage)
	case ike_message.CREATE_CHILD_SA:
		handler.HandleCREATECHILDSA(udpConn, localAddr, remoteAddr, ikeMessage)
	default:
		log.Warnf("Unimplemented IKE message type, exchange type: %d", ikeMessage.ExchangeType)
	}

}
