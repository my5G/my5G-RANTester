package service

import (
	"errors"
	"io"
	"time"

	"git.cs.nctu.edu.tw/calee/sctp"

	"my5G-RANTester/lib/n3iwf/context"
	"my5G-RANTester/lib/n3iwf/ngap"
	"my5G-RANTester/lib/n3iwf/ngap/message"

	log "github.com/sirupsen/logrus"
)

const NGAP_PPID_BigEndian = 0x3c000000

// Run start the N3IWF SCTP process.
func Run() error {
	// n3iwf context
	n3iwfSelf := context.N3IWFSelf()
	// load amf SCTP address slice
	amfSCTPAddresses := n3iwfSelf.AMFSCTPAddresses

	localAddr := new(sctp.SCTPAddr)

	for _, remoteAddr := range amfSCTPAddresses {
		errChan := make(chan error)
		go listenAndServe(localAddr, remoteAddr, errChan)
		if err, ok := <-errChan; ok {
			log.Errorln(err)
			return errors.New("NGAP service run failed")
		}
	}

	return nil
}

func listenAndServe(localAddr, remoteAddr *sctp.SCTPAddr, errChan chan<- error) {
	var conn *sctp.SCTPConn
	var err error

	// Connect the session
	for i := 0; i < 3; i++ {
		conn, err = sctp.DialSCTP("sctp", localAddr, remoteAddr)
		if err != nil {
			log.Errorf("[SCTP] DialSCTP(): %+v", err)
		} else {
			break
		}

		if i != 2 {
			log.Info("Retry to connect AMF after 1 second...")
			time.Sleep(1 * time.Second)
		} else {
			log.Debugf("[SCTP] AMF SCTP address: %+v", remoteAddr.String())
			errChan <- errors.New("Failed to connect to AMF.")
			return
		}
	}

	// Set default sender SCTP infomation sinfo_ppid = NGAP_PPID = 60
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		log.Errorf("[SCTP] GetDefaultSentParam(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			log.Errorf("conn close error in GetDefaultSentParam(): %+v", errConn)
		}
		errChan <- errors.New("Get socket infomation failed.")
		return
	}
	info.PPID = NGAP_PPID_BigEndian
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		log.Errorf("[SCTP] SetDefaultSentParam(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			log.Errorf("conn close error in SetDefaultSentParam(): %+v", errConn)
		}
		errChan <- errors.New("Set socket parameter failed.")
		return
	}

	// Subscribe receiver SCTP information
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		log.Errorf("[SCTP] SubscribeEvents(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			log.Errorf("conn close error in SubscribeEvents(): %+v", errConn)
		}
		errChan <- errors.New("Subscribe SCTP event failed.")
		return
	}

	// Send NG setup request
	go message.SendNGSetupRequest(conn)

	close(errChan)

	data := make([]byte, 65535)

	for {
		n, info, _, err := conn.SCTPRead(data)

		if err != nil {
			log.Debugf("[SCTP] AMF SCTP address: %+v", conn.RemoteAddr().String())
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Warn("[SCTP] Close connection.")
				errConn := conn.Close()
				if errConn != nil {
					log.Errorf("conn close error: %+v", errConn)
				}
				return
			}
			log.Errorf("[SCTP] Read from SCTP connection failed: %+v", err)
		} else {
			log.Tracef("[SCTP] Successfully read %d bytes.", n)

			if info == nil || info.PPID != NGAP_PPID_BigEndian {
				log.Warn("Received SCTP PPID != 60")
				continue
			}

			forwardData := make([]byte, n)
			copy(forwardData, data[:n])

			go ngap.Dispatch(conn, forwardData)
		}
	}
}
