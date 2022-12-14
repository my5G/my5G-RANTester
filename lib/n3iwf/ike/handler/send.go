package handler

import (
	ike_message "my5G-RANTester/lib/n3iwf/ike/message"
	"net"

	log "github.com/sirupsen/logrus"
)

func SendIKEMessageToUE(udpConn *net.UDPConn, srcAddr, dstAddr *net.UDPAddr, message *ike_message.IKEMessage) {
	log.Trace("[IKE] Send IKE message to UE")
	log.Trace("[IKE] Encoding...")
	pkt, err := ike_message.Encode(message)
	if err != nil {
		log.Errorln(err)
		return
	}
	// As specified in RFC 7296 section 3.1, the IKE message send from/to UDP port 4500
	// should prepend a 4 bytes zero
	if srcAddr.Port == 4500 {
		prependZero := make([]byte, 4)
		pkt = append(prependZero, pkt...)
	}

	log.Trace("[IKE] Sending...")
	n, err := udpConn.WriteToUDP(pkt, dstAddr)
	if err != nil {
		log.Error(err)
		return
	}
	if n != len(pkt) {
		log.Errorf("Not all of the data is sent. Total length: %d. Sent: %d.", len(pkt), n)
		return
	}
}
