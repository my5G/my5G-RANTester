package ue_context_management

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/pdu_session_management"
	"my5G-RANTester/lib/ngap"
	"my5G-RANTester/lib/ngap/ngapType"
	"net"
)

func InitialContextSetupResponse(ue *context.GNBUe, allocationPdu bool, ipv4 string) ([]byte, error) {

	var message ngapType.NGAPPDU

	// check hostname.
	nameIp, err := net.LookupHost(ipv4)
	if err != nil {
		return nil, err
	}

	if allocationPdu {
		message = BuildInitialContextSetupResponseWithPduSession(ue.GetAmfUeId(), ue.GetRanUeId(),
			nameIp[0], ue.GetPduSessionId(), ue.GetTeidDownlink(), ue.GetQosId())
	} else {
		message = BuildInitialContextSetupResponseForRegistrationTest(ue.GetAmfUeId(), ue.GetRanUeId())
	}

	return ngap.Encoder(message)
}

func BuildInitialContextSetupResponseForRegistrationTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	return
}

func BuildInitialContextSetupResponseWithPduSession(amfUeNgapID, ranUeNgapID int64, ipv4 string, pduId int64, teid uint32, qosId int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// PDU session resources response
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceSetupListCxtRes
	ie.Value.PDUSessionResourceSetupListCxtRes = new(ngapType.PDUSessionResourceSetupListCxtRes)

	PDUSessionResourceSetupListCxtRes := ie.Value.PDUSessionResourceSetupListCxtRes

	PDUSessionResourceSetupItemCxtRes := ngapType.PDUSessionResourceSetupItemCxtRes{}

	// PDU Session ID
	PDUSessionResourceSetupItemCxtRes.PDUSessionID.Value = pduId

	// get response transfer
	PDUSessionResourceSetupItemCxtRes.PDUSessionResourceSetupResponseTransfer = pdu_session_management.GetPDUSessionResourceSetupResponseTransfer(ipv4, teid, qosId)

	PDUSessionResourceSetupListCxtRes.List = append(PDUSessionResourceSetupListCxtRes.List, PDUSessionResourceSetupItemCxtRes)

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	return
}
