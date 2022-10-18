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
	"reflect"
	"time"
)

func HandlerAuthenticationReject(ue *context.UEContext, message *nas.Message) {

	log.Info("[UE][NAS] Authentication of UE ", ue.GetUeId(), " failed")

	ue.SetStateMM_DEREGISTERED()
}

func HandlerAuthenticationRequest(ue *context.UEContext, message *nas.Message) {
	var authenticationResponse []byte

	// check the mandatory fields
	if reflect.ValueOf(message.AuthenticationRequest.ExtendedProtocolDiscriminator).IsZero() {
		log.Fatal("[UE][NAS] Error in Authentication Request, Extended Protocol is missing")
	}

	if message.AuthenticationRequest.ExtendedProtocolDiscriminator.GetExtendedProtocolDiscriminator() != 126 {
		log.Fatal("[UE][NAS] Error in Authentication Request, Extended Protocol not the expected value")
	}

	if message.AuthenticationRequest.SpareHalfOctetAndSecurityHeaderType.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Authentication Request, Spare Half Octet not the expected value")
	}

	if message.AuthenticationRequest.SpareHalfOctetAndSecurityHeaderType.GetSecurityHeaderType() != 0 {
		log.Fatal("[UE][NAS] Error in Authentication Request, Security Header Type not the expected value")
	}

	if reflect.ValueOf(message.AuthenticationRequest.AuthenticationRequestMessageIdentity).IsZero() {
		log.Fatal("[UE][NAS] Error in Authentication Request, Message Type is missing")
	}

	if message.AuthenticationRequest.AuthenticationRequestMessageIdentity.GetMessageType() != 86 {
		log.Fatal("[UE][NAS] Error in Authentication Request, Message Type not the expected value")
	}

	if message.AuthenticationRequest.SpareHalfOctetAndNgksi.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Authentication Request, Spare Half Octet not the expected value")
	}

	if message.AuthenticationRequest.SpareHalfOctetAndNgksi.GetNasKeySetIdentifiler() == 7 {
		log.Fatal("[UE][NAS] Error in Authentication Request, ngKSI not the expected value")
	}

	if reflect.ValueOf(message.AuthenticationRequest.ABBA).IsZero() {
		log.Fatal("[UE][NAS] Error in Authentication Request, ABBA is missing")
	}

	if message.AuthenticationRequest.GetABBAContents() == nil {
		log.Fatal("[UE][NAS] Error in Authentication Request, ABBA Content is missing")
	}

	// getting RAND and AUTN from the message.
	rand := message.AuthenticationRequest.GetRANDValue()
	autn := message.AuthenticationRequest.GetAUTN()

	// getting 5G NAS security identifier.
	ngksi := message.AuthenticationRequest.GetNasKeySetIdentifiler()
	ue.SetNgKsi(ngksi)

	if ue.GetTesting() == "test-authentication-reject" {
		// produce wrong messages in the process of authentication
		// getting resStar
		paramAutn := ue.DeriveRESstarAndSetKeyWrongly(ue.UeSecurity.AuthenticationSubs, rand[:], ue.UeSecurity.Snn, autn[:])
		log.Info("[UE][NAS] Send authentication response")
		log.Info("[UE][NAS] 5G NAS security identifier: ", ue.GetNgKsi())
		authenticationResponse = mm_5gs.AuthenticationResponse(paramAutn, "")
	} else {
		// getting resStar
		paramAutn, check := ue.DeriveRESstarAndSetKey(ue.UeSecurity.AuthenticationSubs, rand[:], ue.UeSecurity.Snn, autn[:])
		// handled correctly the process of authentication
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
			log.Info("[UE][NAS] 5G NAS security identifier: ", ue.GetNgKsi())
			authenticationResponse = mm_5gs.AuthenticationResponse(paramAutn, "")

			// change state of UE for registered-initiated
			ue.SetStateMM_REGISTERED_INITIATED()
		}
	}

	// sending to GNB
	sender.SendToGnb(ue, authenticationResponse)
}

func HandlerSecurityModeCommand(ue *context.UEContext, message *nas.Message) {

	// check the mandatory fields
	if reflect.ValueOf(message.SecurityModeCommand.ExtendedProtocolDiscriminator).IsZero() {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Extended Protocol is missing")
	}

	if message.SecurityModeCommand.ExtendedProtocolDiscriminator.GetExtendedProtocolDiscriminator() != 126 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Extended Protocol not the expected value")
	}

	if message.SecurityModeCommand.SpareHalfOctetAndSecurityHeaderType.GetSecurityHeaderType() != 0 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Security Header Type not the expected value")
	}

	if message.SecurityModeCommand.SpareHalfOctetAndSecurityHeaderType.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Spare Half Octet not the expected value")
	}

	if reflect.ValueOf(message.SecurityModeCommand.SecurityModeCommandMessageIdentity).IsZero() {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Message Type is missing")
	}

	if message.SecurityModeCommand.SecurityModeCommandMessageIdentity.GetMessageType() != 93 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Message Type not the expected value")
	}

	if reflect.ValueOf(message.SecurityModeCommand.SelectedNASSecurityAlgorithms).IsZero() {
		log.Fatal("[UE][NAS] Error in Security Mode Command, NAS Security Algorithms is missing")
	}

	if message.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm() != 0 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, NAS Security Algorithms not the expected value")
	}

	if message.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm() != 2 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, NAS Security Algorithms not the expected value")
	}

	if message.SecurityModeCommand.SpareHalfOctetAndNgksi.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Spare Half Octet is missing")
	}

	if message.SecurityModeCommand.SpareHalfOctetAndNgksi.GetNasKeySetIdentifiler() == 7 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, ngKSI not the expected value")
	}

	if reflect.ValueOf(message.SecurityModeCommand.ReplayedUESecurityCapabilities).IsZero() {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Replayed UE Security Capabilities is missing")
	}

	if message.SecurityModeCommand.ReplayedUESecurityCapabilities.GetEA0_5G() != 1 || message.SecurityModeCommand.ReplayedUESecurityCapabilities.GetIA2_128_5G() != 1 || message.SecurityModeCommand.ReplayedUESecurityCapabilities.GetIA1_128_5G() != 0 || message.SecurityModeCommand.ReplayedUESecurityCapabilities.GetIA3_128_5G() != 0 {
		log.Fatal("[UE][NAS] Error in Security Mode Command, Replayed UE Security Capabilities not the expected value")
	}

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

	// send invalid flows.
	if ue.GetTesting() == "test-invalid-flows" {

		// getting ul nas transport and pduSession establishment request.
		ulNasTransport, err := mm_5gs.UlNasTransport(ue, nasMessage.ULNASTransportRequestTypeInitialRequest, "test-invalid-flows")
		if err != nil {
			log.Fatal("[UE][NAS] Error sending ul nas transport and pdu session establishment request: ", err)
		}

		// sending to GNB
		sender.SendToGnb(ue, ulNasTransport)

	} else {

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
}

func HandlerRegistrationAccept(ue *context.UEContext, message *nas.Message) {

	// check the mandatory fields
	if reflect.ValueOf(message.RegistrationAccept.ExtendedProtocolDiscriminator).IsZero() {
		log.Fatal("[UE][NAS] Error in Registration Accept, Extended Protocol is missing")
	}

	if message.RegistrationAccept.ExtendedProtocolDiscriminator.GetExtendedProtocolDiscriminator() != 126 {
		log.Fatal("[UE][NAS] Error in Registration Accept, Extended Protocol not the expected value")
	}

	if message.RegistrationAccept.SpareHalfOctetAndSecurityHeaderType.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Registration Accept, Spare Half not the expected value")
	}

	if message.RegistrationAccept.SpareHalfOctetAndSecurityHeaderType.GetSecurityHeaderType() != 0 {
		log.Fatal("[UE][NAS] Error in Registration Accept, Security Header not the expected value")
	}

	if reflect.ValueOf(message.RegistrationAccept.RegistrationAcceptMessageIdentity).IsZero() {
		log.Fatal("[UE][NAS] Error in Registration Accept, Message Type is missing")
	}

	if message.RegistrationAccept.RegistrationAcceptMessageIdentity.GetMessageType() != 66 {
		log.Fatal("[UE][NAS] Error in Registration Accept, Message Type not the expected value")
	}

	if reflect.ValueOf(message.RegistrationAccept.RegistrationResult5GS).IsZero() {
		log.Fatal("[UE][NAS] Error in Registration Accept, Registration Result 5GS is missing")
	}

	if message.RegistrationAccept.RegistrationResult5GS.GetRegistrationResultValue5GS() != 1 {
		log.Fatal("[UE][NAS] Error in Registration Accept, Registration Result 5GS not the expected value")
	}

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
	ulNasTransport, err := mm_5gs.UlNasTransport(ue, nasMessage.ULNASTransportRequestTypeInitialRequest, "")
	if err != nil {
		log.Fatal("[UE][NAS] Error sending ul nas transport and pdu session establishment request: ", err)
	}

	// change the sate of ue(SM).
	ue.SetStateSM_PDU_SESSION_PENDING()

	// sending to GNB
	sender.SendToGnb(ue, ulNasTransport)

	// testing Ul NAS transport replicate
	if ue.GetTesting() == "test-duplicate-messages" {
		time.Sleep(2 * time.Millisecond)
		sender.SendToGnb(ue, ulNasTransport)
		time.Sleep(1 * time.Second)
		sender.SendToGnb(ue, ulNasTransport)
	}
}

func HandlerDlNasTransportPduaccept(ue *context.UEContext, message *nas.Message) {

	// check the mandatory fields
	if reflect.ValueOf(message.DLNASTransport.ExtendedProtocolDiscriminator).IsZero() {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Extended Protocol is missing")
	}

	if message.DLNASTransport.ExtendedProtocolDiscriminator.GetExtendedProtocolDiscriminator() != 126 {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Extended Protocol not expected value")
	}

	if message.DLNASTransport.SpareHalfOctetAndSecurityHeaderType.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Spare Half not expected value")
	}

	if message.DLNASTransport.SpareHalfOctetAndSecurityHeaderType.GetSecurityHeaderType() != 0 {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Security Header not expected value")
	}

	if message.DLNASTransport.DLNASTRANSPORTMessageIdentity.GetMessageType() != 104 {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Message Type is missing or not expected value")
	}

	if reflect.ValueOf(message.DLNASTransport.SpareHalfOctetAndPayloadContainerType).IsZero() {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Payload Container Type is missing")
	}

	if message.DLNASTransport.SpareHalfOctetAndPayloadContainerType.GetPayloadContainerType() != 1 {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Payload Container Type not expected value")
	}

	if reflect.ValueOf(message.DLNASTransport.PayloadContainer).IsZero() || message.DLNASTransport.PayloadContainer.GetPayloadContainerContents() == nil {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, Payload Container is missing")
	}

	if reflect.ValueOf(message.DLNASTransport.PduSessionID2Value).IsZero() {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, PDU Session ID is missing")
	}

	if message.DLNASTransport.PduSessionID2Value.GetIei() != 18 {
		log.Fatal("[UE][NAS] Error in DL NAS Transport, PDU Session ID not expected value")
	}

	// handle errors in DL NAS TRANSPORT
	if message.DLNASTransport.Cause5GMM != nil {
		if message.DLNASTransport.Cause5GMM.GetCauseValue() == 91 {
			log.Fatal("[UE][NAS] Receive Error in Selection of SMF in DL NAS TRANSPORT")
		}
	}

	//getting PDU Session establishment accept.
	payloadContainer := nas_control.GetNasPduFromPduAccept(message)
	if payloadContainer.GsmHeader.GetMessageType() == nas.MsgTypePDUSessionEstablishmentAccept {

		log.Info("[UE][NAS] Receiving PDU Session Establishment Accept")

		// check the mandatory fields
		if reflect.ValueOf(payloadContainer.PDUSessionEstablishmentAccept.ExtendedProtocolDiscriminator).IsZero() {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, Extended Protocol Discriminator is missing")
		}

		if payloadContainer.PDUSessionEstablishmentAccept.GetExtendedProtocolDiscriminator() != 46 {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, Extended Protocol Discriminator not expected value")
		}

		if reflect.ValueOf(payloadContainer.PDUSessionEstablishmentAccept.PDUSessionID).IsZero() {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, PDU Session ID is missing or not expected value")
		}

		if reflect.ValueOf(payloadContainer.PDUSessionEstablishmentAccept.PTI).IsZero() {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, PTI is missing")
		}

		if payloadContainer.PDUSessionEstablishmentAccept.PTI.GetPTI() != 1 {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, PTI not the expected value")
		}

		if payloadContainer.PDUSessionEstablishmentAccept.PDUSESSIONESTABLISHMENTACCEPTMessageIdentity.GetMessageType() != 194 {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, Message Type is missing or not expected value")
		}

		if reflect.ValueOf(payloadContainer.PDUSessionEstablishmentAccept.SelectedSSCModeAndSelectedPDUSessionType).IsZero() {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, SSC Mode or PDU Session Type is missing")
		}

		if payloadContainer.PDUSessionEstablishmentAccept.SelectedSSCModeAndSelectedPDUSessionType.GetPDUSessionType() != 1 {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, PDU Session Type not the expected value")
		}

		if reflect.ValueOf(payloadContainer.PDUSessionEstablishmentAccept.AuthorizedQosRules).IsZero() {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, Authorized QoS Rules is missing")
		}

		if reflect.ValueOf(payloadContainer.PDUSessionEstablishmentAccept.SessionAMBR).IsZero() {
			log.Fatal("[UE][NAS] Error in PDU Session Establishment Accept, Session AMBR is missing")
		}

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
	} else if payloadContainer.GsmHeader.GetMessageType() == nas.MsgTypePDUSessionEstablishmentReject {
		// handler PDU Establishment Reject
		log.Fatal("[UE][NAS] Receive PDU Establishment Reject")
	}
}

func HandlerIdentityRequest(ue *context.UEContext, message *nas.Message) {

	var identity5gs string

	// check the mandatory fields
	if reflect.ValueOf(message.IdentityRequest.ExtendedProtocolDiscriminator).IsZero() {
		log.Fatal("[UE][NAS] Error in Identity Request, Extended Protocol is missing")
	}

	if message.IdentityRequest.ExtendedProtocolDiscriminator.GetExtendedProtocolDiscriminator() != 126 {
		log.Fatal("[UE][NAS] Error in Identity Request, Extended Protocol not the expected value")
	}

	if message.IdentityRequest.SpareHalfOctetAndSecurityHeaderType.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Identity Request, Spare Half Octet not the expected value")
	}

	if message.IdentityRequest.SpareHalfOctetAndSecurityHeaderType.GetSecurityHeaderType() != 0 {
		log.Fatal("[UE][NAS] Error in Identity Request, Security Header Type not the expected value")
	}

	if reflect.ValueOf(message.IdentityRequest.IdentityRequestMessageIdentity).IsZero() {
		log.Fatal("[UE][NAS] Error in Identity Request, Message Type is missing")
	}

	if message.IdentityRequest.IdentityRequestMessageIdentity.GetMessageType() != 91 {
		log.Fatal("[UE][NAS] Error in Identity Request, Message Type not the expected value")
	}

	if reflect.ValueOf(message.IdentityRequest.SpareHalfOctetAndIdentityType).IsZero() {
		log.Fatal("[UE][NAS] Error in Identity Request, Spare Half Octet And Identity Type is missing")
	}

	if ue.Test == "test-5g-guti" {

		if message.IdentityRequest.SpareHalfOctetAndIdentityType.GetTypeOfIdentity() != 1 {
			log.Fatal("[UE][NAS] Error in Identity Request, Type of Identity not the expected value")
		}

	}

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

	// check the mandatory fields
	if reflect.ValueOf(message.ConfigurationUpdateCommand.ExtendedProtocolDiscriminator).IsZero() {
		log.Fatal("[UE][NAS] Error in Configuration Update Command, Extended Protocol Discriminator is missing")
	}

	if message.ConfigurationUpdateCommand.ExtendedProtocolDiscriminator.GetExtendedProtocolDiscriminator() != 126 {
		log.Fatal("[UE][NAS] Error in Configuration Update Command, Extended Protocol Discriminator not the expected value")
	}

	if message.ConfigurationUpdateCommand.SpareHalfOctetAndSecurityHeaderType.GetSpareHalfOctet() != 0 {
		log.Fatal("[UE][NAS] Error in Configuration Update Command, Spare Half not the expected value")
	}

	if message.ConfigurationUpdateCommand.SpareHalfOctetAndSecurityHeaderType.GetSecurityHeaderType() != 0 {
		log.Fatal("[UE][NAS] Error in Configuration Update Command, Security Header not the expected value")
	}

	if reflect.ValueOf(message.ConfigurationUpdateCommand.ConfigurationUpdateCommandMessageIdentity).IsZero() {
		log.Fatal("[UE][NAS] Error in Configuration Update Command, Message type not the expected value")
	}

	if message.ConfigurationUpdateCommand.ConfigurationUpdateCommandMessageIdentity.GetMessageType() != 84 {
		log.Fatal("[UE][NAS] Error in Configuration Update Command, Message Type not the expected value")
	}

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
