package non3gppue

import (
	"fmt"
	"math/big"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/non3gppue/context"
	"my5G-RANTester/internal/control_test_engine/non3gppue/ike/handler"
	"my5G-RANTester/internal/control_test_engine/non3gppue/ike/message"
	"my5G-RANTester/internal/control_test_engine/non3gppue/nas/service"
	"my5G-RANTester/internal/control_test_engine/non3gppue/nas/trigger"
	"my5G-RANTester/internal/monitoring"
	"my5G-RANTester/lib/nas/security"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func RegistrationUe(conf config.Config, id uint8, wg *sync.WaitGroup) {
	n3iwfIP := "10.100.200.14"

	// wg := sync.WaitGroup{}

	// new UE instance.
	ue := &context.UEContext{}

	// new UE context
	ue.NewRanUeContext(
		conf.Ue.Msin,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		"c9e8763286b5b9ffbdf56e1297d0887b", // This is for the Op field that problably is not used
		conf.Ue.Amf,
		conf.Ue.Sqn,
		conf.Ue.Hplmn.Mcc,
		conf.Ue.Hplmn.Mnc,
		conf.Ue.Dnn,
		int32(conf.Ue.Snssai.Sst),
		conf.Ue.Snssai.Sd,
		id)

	// Initiate UDP socket for N3IWF
	n3iwfUDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:500", n3iwfIP))
	if err != nil {
		log.Fatal(err)
	}
	udpConnection := setupUDPSocket(ue)

	// Starting IKE Protocol
	// IKE_SA_INIT
	ikeMessage := message.BuildIKEHeader(123123, 0, message.IKE_SA_INIT, message.InitiatorBitCheck, 0)

	// Security Association
	proposal := message.BuildProposal(1, message.TypeIKE, nil)
	var attributeType uint16 = message.AttributeTypeKeyLength
	var keyLength uint16 = 256
	encryptTransform := message.BuildTransform(message.TypeEncryptionAlgorithm, message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	message.AppendTransformToProposal(proposal, encryptTransform)
	integrityTransform := message.BuildTransform(message.TypeIntegrityAlgorithm, message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	message.AppendTransformToProposal(proposal, integrityTransform)
	pseudorandomFunctionTransform := message.BuildTransform(message.TypePseudorandomFunction, message.PRF_HMAC_SHA1, nil, nil, nil)
	message.AppendTransformToProposal(proposal, pseudorandomFunctionTransform)
	diffiehellmanTransform := message.BuildTransform(message.TypeDiffieHellmanGroup, message.DH_2048_BIT_MODP, nil, nil, nil)
	message.AppendTransformToProposal(proposal, diffiehellmanTransform)
	securityAssociation := message.BuildSecurityAssociation([]*message.Proposal{proposal})
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, securityAssociation)

	// Key exchange data
	generator := new(big.Int).SetUint64(handler.Group14Generator)
	factor, ok := new(big.Int).SetString(handler.Group14PrimeString, 16)
	if !ok {
		log.Fatal("Generate key exchange datd failed")
	}
	secert := handler.GenerateRandomNumber()
	localPublicKeyExchangeValue := new(big.Int).Exp(generator, secert, factor).Bytes()
	prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
	localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
	keyExchangeData := message.BUildKeyExchange(message.DH_2048_BIT_MODP, localPublicKeyExchangeValue)
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, keyExchangeData)

	// Nonce
	localNonce := handler.GenerateRandomNumber().Bytes()
	nonce := message.BuildNonce(localNonce)
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, nonce)

	// Send to N3IWF
	ikeMessageData, err := message.Encode(ikeMessage)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		log.Fatal(err)
	}

	log.Info("THIS IS THE END!!")
	os.Exit(0)

	// registration procedure started.
	// trigger.InitRegistration(ue)

	// wg.Wait()

	// control the signals
	sigUe := make(chan os.Signal, 1)
	signal.Notify(sigUe, os.Interrupt)

	// Block until a signal is received.
	<-sigUe
	ue.Terminate()
	wg.Done()
	// os.Exit(0)

}

func setupUDPSocket(ctx *context.UEContext) *net.UDPConn {
	bindAddr := fmt.Sprintf("%s:500", "192.168.2.209") // Static IP
	udpAddr, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		log.Fatal("Resolve UDP address failed")
	}
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Listen UDP socket failed: %+v", err)
	}
	return udpListener
}

func RegistrationUeMonitor(conf config.Config,
	id uint8, monitor *monitoring.Monitor, wg *sync.WaitGroup, start time.Time) {

	// new UE instance.
	ue := &context.UEContext{}

	// new UE context
	ue.NewRanUeContext(
		conf.Ue.Msin,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		"c9e8763286b5b9ffbdf56e1297d0887b",
		conf.Ue.Amf,
		conf.Ue.Sqn,
		conf.Ue.Hplmn.Mcc,
		conf.Ue.Hplmn.Mnc,
		conf.Ue.Dnn,
		int32(conf.Ue.Snssai.Sst),
		conf.Ue.Snssai.Sd,
		id)

	// starting communication with GNB and listen.
	err := service.InitConn(ue)
	if err != nil {
		log.Fatal("Error in", err)
	} else {
		log.Info("[UE] UNIX/NAS service is running")
		// wg.Add(1)
	}

	// registration procedure started.
	trigger.InitRegistration(ue)

	for {

		// UE is register in network
		if ue.GetStateMM() == 0x03 {
			elapsed := time.Since(start)
			monitor.LtRegisterLocal = elapsed.Milliseconds()
			log.Warn("[TESTER][UE] UE LATENCY IN REGISTRATION: ", monitor.LtRegisterLocal, " ms")
			break
		}

		// timeout is 10 000 ms
		if time.Since(start).Milliseconds() >= 10000 {
			log.Warn("[TESTER][UE] TIME EXPIRED IN UE REGISTRATION 10 000 ms")
			break
		}
	}

	wg.Done()
	// ue.Terminate()
	// os.Exit(0)
}
