package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control/mm_5gs"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/sender"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
	"time"
)

func HandlerAuthenticationReject(ue *context.UEContext, message *nas.Message) {

	log.Info("[UE][NAS] Authentication of UE ", ue.GetUeId(), " failed")

	ue.SetStateMM_DEREGISTERED()
}

func HandlerAuthenticationRequest(ue *context.UEContext, message *nas.Message) {
	var authenticationResponse []byte

	// getting RAND and AUTN from the message.
	rand := message.AuthenticationRequest.GetRANDValue()
	autn := message.AuthenticationRequest.GetAUTN()

	// getting resStar
	paramAutn, check := ue.DeriveRESstarAndSetKey(ue.UeSecurity.AuthenticationSubs, rand[:], ue.UeSecurity.Snn, autn[:])

	switch check {

	case "MAC failure":
		log.Info("[UE][NAS][MAC] Authenticity of the authentication request message: FAILED")
		log.Info("[UE][NAS] Send authentication failure with MAC failure")
		authenticationResponse = mm_5gs.AuthenticationFailure("MAC failure", "", paramAutn)
		// not change the state of UE.

	case "SQN failure":
		log.Info("[UE][NAS][MAC] Authenticity of the authentication request message: OK")
		log.Info("[UE][NAS][SQN] SQN of the authentication request message: INVALID")
		log.Info("[UE][NAS] Send authentication failure with Synch failure")
		authenticationResponse = mm_5gs.AuthenticationFailure("SQN failure", "", paramAutn)
		// not change the state of UE.

	case "successful":
		// getting NAS Authentication Response.
		log.Info("[UE][NAS][MAC] Authenticity of the authentication request message: OK")
		log.Info("[UE][NAS][SQN] SQN of the authentication request message: VALID")
		log.Info("[UE][NAS] Send authentication response")
		authenticationResponse = mm_5gs.AuthenticationResponse(paramAutn, "")

		// change state of UE for registered-initiated
		ue.SetStateMM_REGISTERED_INITIATED()
	}

	// sending to GNB
	sender.SendToGnb(ue, authenticationResponse)
}

func HandlerSecurityModeCommand(ue *context.UEContext, message *nas.Message) {

	switch message.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm() {
	case 0:
		log.Info("[UE][NAS] Type of ciphering algorithm is 5G-EA0")
	case 1:
		log.Info("[UE][NAS] Type of ciphering algorithm is 128-5G-EA1")
	case 2:
		log.Info("[UE][NAS] Type of ciphering algorithm is 128-5G-EA2")
	}

	switch message.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm() {
	case 0:
		log.Info("[UE][NAS] Type of integrity protection algorithm is 5G-IA0")
	case 1:
		log.Info("[UE][NAS] Type of integrity protection algorithm is 128-5G-IA1")
	case 2:
		log.Info("[UE][NAS] Type of integrity protection algorithm is 128-5G-IA2")
	}

	// checking BIT RINMR that triggered registration request in security mode complete.
	rinmr := message.SecurityModeCommand.Additional5GSecurityInformation.GetRINMR()

	// getting NAS Security Mode Complete.
	securityModeComplete, err := mm_5gs.SecurityModeComplete(ue, rinmr)
	if err != nil {
		log.Fatal("[UE][NAS] Error sending Security Mode Complete: ", err)
	}

	// sending to GNB
	sender.SendToGnb(ue, securityModeComplete)
}

func HandlerRegistrationAccept(ue *context.UEContext, message *nas.Message) {

	// change the state of ue for registered
	ue.SetStateMM_REGISTERED()

	// saved 5g GUTI and others information.
	ue.SetAmfRegionId(message.RegistrationAccept.GetAMFRegionID())
	ue.SetAmfPointer(message.RegistrationAccept.GetAMFPointer())
	ue.SetAmfSetId(message.RegistrationAccept.GetAMFSetID())
	ue.Set5gGuti(message.RegistrationAccept.GetTMSI5G())

	// check allowed slices.
	log.Info("[UE][NAS] Allowed NSSAI", message.RegistrationAccept.AllowedNSSAI.GetSNSSAIValue())

	log.Info("[UE][NAS] UE 5G GUTI: ", ue.Get5gGuti())

	// getting NAS registration complete.
	registrationComplete, err := mm_5gs.RegistrationComplete(ue)
	if err != nil {
		log.Fatal("[UE][NAS] Error sending Registration Complete: ", err)
	}

	// sending to GNB
	sender.SendToGnb(ue, registrationComplete)

	// waiting receive Configuration Update Command.
	time.Sleep(20 * time.Millisecond)

	// getting ul nas transport and pduSession establishment request.
	ulNasTransport, err := mm_5gs.UlNasTransport(ue, nasMessage.ULNASTransportRequestTypeInitialRequest)
	if err != nil {
		log.Fatal("[UE][NAS] Error sending ul nas transport and pdu session establishment request: ", err)
	}

	// change the sate of ue(SM).
	ue.SetStateSM_PDU_SESSION_PENDING()

	// sending to GNB
	sender.SendToGnb(ue, ulNasTransport)
}

func HandlerDlNasTransportPduaccept(ue *context.UEContext, message *nas.Message) {

	//getting PDU Session establishment accept.
	payloadContainer := nas_control.GetNasPduFromPduAccept(message)
	if payloadContainer.GsmHeader.GetMessageType() == nas.MsgTypePDUSessionEstablishmentAccept {
		log.Info("[UE][NAS] Receiving PDU Session Establishment Accept")

		// update PDU Session information.

		// change the state of ue(SM)(PDU Session Active).
		ue.SetStateSM_PDU_SESSION_ACTIVE()

		// get QoS Rules
		QosRule := payloadContainer.PDUSessionEstablishmentAccept.AuthorizedQosRules.GetQosRule()

		// get PDU session IP.
		UeIp := payloadContainer.PDUSessionEstablishmentAccept.GetPDUAddressInformation()

		// get DNN
		dnn := payloadContainer.PDUSessionEstablishmentAccept.DNN.GetDNN()

		// get SNSSAI
		sst := payloadContainer.PDUSessionEstablishmentAccept.SNSSAI.GetSST()
		sd := payloadContainer.PDUSessionEstablishmentAccept.SNSSAI.GetSD()

		// set UE ip
		ue.SetIp(UeIp)

		log.Info("[UE][NAS] PDU session QoS RULES: ", QosRule)
		log.Info("[UE][NAS] PDU session DNN: ", string(dnn))
		log.Info("[UE][NAS] PDU session NSSAI -- sst: ", sst, " sd: ",
			fmt.Sprintf("%x%x%x", sd[0], sd[1], sd[2]))
		log.Info("[UE][NAS] PDU address received: ", ue.GetIp())
	}
}

func HandlerIdentityRequest(ue *context.UEContext, message *nas.Message) {

	var identity5gs string

	type5gs := message.IdentityRequest.GetTypeOfIdentity()

	switch type5gs {

	case 1:
		log.Info("[UE][NAS] Requested SUCI 5GS type")
		identity5gs = "suci"

	case 2:
		log.Info("[UE][NAS] Requested 5G-GUTI 5GS type")
		identity5gs = "guti"

	case 3:
		log.Info("[UE][NAS] Requested IMEI 5GS type")
		identity5gs = "imei"

	case 4:
	case 5:
		log.Info("[UE][NAS] Requested IMEISV 5GS type")
		identity5gs = "imeisv"

	}

	// trigger identity response.
	identityResponse := mm_5gs.IdentityResponse(identity5gs, ue)

	// send to GNB.
	sender.SendToGnb(ue, identityResponse)
}

func HandlerConfigurationUpdateCommand(ue *context.UEContext, message *nas.Message) {

	networkName := message.ConfigurationUpdateCommand.FullNameForNetwork.GetTextString()
	log.Info("[UE][NAS] Network Name: ", string(networkName))

	// time zone
	timeZone := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetTimeZone()
	log.Info("[UE][NAS] Time Zone: ", timeZone)

	//time
	/*
		year := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetYear()
		day := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetDay()
		mounth := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetMonth()
		hour := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetHour()
		minute := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetMinute()
		second := message.ConfigurationUpdateCommand.UniversalTimeAndLocalTimeZone.GetSecond()
		log.Info("[UE][NAS] Time: ", mounth, year, day, hour, ":", minute, ":", second)
	*/

	// return configuration update complete
	//message.ConfigurationUpdateCommand.ConfigurationUpdateIndication.GetACK()
}
