package message

import (
	"git.cs.nctu.edu.tw/calee/sctp"

	"my5G-RANTester/lib/n3iwf/context"
	"my5G-RANTester/lib/ngap/ngapType"

	log "github.com/sirupsen/logrus"
)

func SendToAmf(amf *context.N3IWFAMF, pkt []byte) {
	if amf == nil {
		log.Errorf("[N3IWF] AMF Context is nil ")
	} else {
		if n, err := amf.SCTPConn.Write(pkt); err != nil {
			log.Errorf("Write to SCTP socket failed: %+v", err)
		} else {
			log.Tracef("Wrote %d bytes", n)
		}
	}
}

func SendNGSetupRequest(conn *sctp.SCTPConn) {
	log.Infoln("[N3IWF] Send NG Setup Request")

	sctpAddr := conn.RemoteAddr().String()

	if available, _ := context.N3IWFSelf().AMFReInitAvailableListLoad(sctpAddr); !available {
		log.Warnf("[N3IWF] Please Wait at least for the indicated time before reinitiating toward same AMF[%s]", sctpAddr)
		return
	}
	pkt, err := BuildNGSetupRequest()
	if err != nil {
		log.Errorf("Build NGSetup Request failed: %+v\n", err)
		return
	}

	if n, err := conn.Write(pkt); err != nil {
		log.Errorf("Write to SCTP socket failed: %+v", err)
	} else {
		log.Tracef("Wrote %d bytes", n)
	}
}

// partOfNGInterface: if reset type is "reset all", set it to nil TS 38.413 9.2.6.11
func SendNGReset(
	amf *context.N3IWFAMF,
	cause ngapType.Cause,
	partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList) {

	log.Infoln("[N3IWF] Send NG Reset")

	pkt, err := BuildNGReset(cause, partOfNGInterface)
	if err != nil {
		log.Errorf("Build NGReset failed : %s", err.Error())
		return
	}

	SendToAmf(amf, pkt)
}

func SendNGResetAcknowledge(
	amf *context.N3IWFAMF,
	partOfNGInterface *ngapType.UEAssociatedLogicalNGConnectionList,
	diagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send NG Reset Acknowledge")

	if partOfNGInterface != nil && len(partOfNGInterface.List) == 0 {
		log.Error("length of partOfNGInterface is 0")
		return
	}

	pkt, err := BuildNGResetAcknowledge(partOfNGInterface, diagnostics)
	if err != nil {
		log.Errorf("Build NGReset Acknowledge failed : %s", err.Error())
		return
	}

	SendToAmf(amf, pkt)
}

func SendInitialContextSetupResponse(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	responseList *ngapType.PDUSessionResourceSetupListCxtRes,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send Initial Context Setup Response")

	if responseList != nil && len(responseList.List) > context.MaxNumOfPDUSessions {
		log.Errorln("Pdu List out of range")
		return
	}

	if failedList != nil && len(failedList.List) > context.MaxNumOfPDUSessions {
		log.Errorln("Pdu List out of range")
		return
	}

	pkt, err := BuildInitialContextSetupResponse(ue, responseList, failedList, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build Initial Context Setup Response failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendInitialContextSetupFailure(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	cause ngapType.Cause,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send Initial Context Setup Failure")

	if failedList != nil && len(failedList.List) > context.MaxNumOfPDUSessions {
		log.Errorln("Pdu List out of range")
		return
	}

	pkt, err := BuildInitialContextSetupFailure(ue, cause, failedList, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build Initial Context Setup Failure failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendUEContextModificationResponse(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send UE Context Modification Response")

	pkt, err := BuildUEContextModificationResponse(ue, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build UE Context Modification Response failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendUEContextModificationFailure(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	cause ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send UE Context Modification Failure")

	pkt, err := BuildUEContextModificationFailure(ue, cause, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build UE Context Modification Failure failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendUEContextReleaseComplete(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send UE Context Release Complete")

	pkt, err := BuildUEContextReleaseComplete(ue, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build UE Context Release Complete failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendUEContextReleaseRequest(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe, cause ngapType.Cause) {

	log.Infoln("[N3IWF] Send UE Context Release Request")

	pkt, err := BuildUEContextReleaseRequest(ue, cause)
	if err != nil {
		log.Errorf("Build UE Context Release Request failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendInitialUEMessage(amf *context.N3IWFAMF,
	ue *context.N3IWFUe, nasPdu []byte) {
	log.Infoln("[N3IWF] Send Initial UE Message")
	// Attach To AMF

	pkt, err := BuildInitialUEMessage(ue, nasPdu, nil)
	if err != nil {
		log.Errorf("Build Initial UE Message failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
	// ue.AttachAMF()
}

func SendUplinkNASTransport(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	nasPdu []byte) {

	log.Infoln("[N3IWF] Send Uplink NAS Transport")

	if len(nasPdu) == 0 {
		log.Errorln("NAS Pdu is nil")
		return
	}

	pkt, err := BuildUplinkNASTransport(ue, nasPdu)
	if err != nil {
		log.Errorf("Build Uplink NAS Transport failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendNASNonDeliveryIndication(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	nasPdu []byte,
	cause ngapType.Cause) {
	log.Infoln("[N3IWF] Send NAS NonDelivery Indication")

	if len(nasPdu) == 0 {
		log.Errorln("NAS Pdu is nil")
		return
	}

	pkt, err := BuildNASNonDeliveryIndication(ue, nasPdu, cause)
	if err != nil {
		log.Errorf("Build Uplink NAS Transport failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendRerouteNASRequest() {
	log.Infoln("[N3IWF] Send Reroute NAS Request")
}

func SendPDUSessionResourceSetupResponse(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	responseList *ngapType.PDUSessionResourceSetupListSURes,
	failedListSURes *ngapType.PDUSessionResourceFailedToSetupListSURes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send PDU Session Resource Setup Response")

	if ue == nil {
		log.Error("UE context is nil, this information is mandatory.")
		return
	}

	pkt, err := BuildPDUSessionResourceSetupResponse(ue, responseList, failedListSURes, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build PDU Session Resource Setup Response failed : %+v", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendPDUSessionResourceModifyResponse(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	responseList *ngapType.PDUSessionResourceModifyListModRes,
	failedList *ngapType.PDUSessionResourceFailedToModifyListModRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send PDU Session Resource Modify Response")

	if ue == nil && criticalityDiagnostics == nil {
		log.Error("UE context is nil, this information is mandatory")
		return
	}

	pkt, err := BuildPDUSessionResourceModifyResponse(ue, responseList, failedList, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build PDU Session Resource Modify Response failed : %+v", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendPDUSessionResourceModifyIndication(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	modifyList []ngapType.PDUSessionResourceModifyItemModInd) {

	log.Infoln("[N3IWF] Send PDU Session Resource Modify Indication")

	if ue == nil {
		log.Error("UE context is nil, this information is mandatory")
		return
	}
	if modifyList == nil {
		log.Errorln("PDU Session Resource Modify Indication List is nil. This message shall contain at least one Item")
		return
	}

	pkt, err := BuildPDUSessionResourceModifyIndication(ue, modifyList)
	if err != nil {
		log.Errorf("Build PDU Session Resource Modify Indication failed : %+v", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendPDUSessionResourceNotify(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	notiList *ngapType.PDUSessionResourceNotifyList,
	relList *ngapType.PDUSessionResourceReleasedListNot) {

	log.Infoln("[N3IWF] Send PDU Session Resource Notify")

	if ue == nil {
		log.Error("UE context is nil, this information is mandatory")
		return
	}

	pkt, err := BuildPDUSessionResourceNotify(ue, notiList, relList)
	if err != nil {
		log.Errorf("Build PDUSession Resource Notify failed : %+v", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendPDUSessionResourceReleaseResponse(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	relList ngapType.PDUSessionResourceReleasedListRelRes,
	diagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send PDU Session Resource Release Response")

	if ue == nil {
		log.Error("UE context is nil, this information is mandatory")
		return
	}
	if len(relList.List) < 1 {
		log.Errorln("PDUSessionResourceReleasedListRelRes is nil. This message shall contain at least one Item")
		return
	}

	pkt, err := BuildPDUSessionResourceReleaseResponse(ue, relList, diagnostics)
	if err != nil {
		log.Errorf("Build PDU Session Resource Release Response failed : %+v", err)
		return
	}

	SendToAmf(amf, pkt)

}

func SendErrorIndication(
	amf *context.N3IWFAMF,
	amfUENGAPID *int64,
	ranUENGAPID *int64,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send Error Indication")

	if (cause == nil) && (criticalityDiagnostics == nil) {
		log.Errorln("Both cause and criticality is nil. This message shall contain at least one of them.")
		return
	}

	pkt, err := BuildErrorIndication(amfUENGAPID, ranUENGAPID, cause, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build Error Indication failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendErrorIndicationWithSctpConn(
	sctpConn *sctp.SCTPConn,
	amfUENGAPID *int64,
	ranUENGAPID *int64,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send Error Indication")

	if (cause == nil) && (criticalityDiagnostics == nil) {
		log.Errorln("Both cause and criticality is nil. This message shall contain at least one of them.")
		return
	}

	pkt, err := BuildErrorIndication(amfUENGAPID, ranUENGAPID, cause, criticalityDiagnostics)
	if err != nil {
		log.Errorf("Build Error Indication failed : %+v\n", err)
		return
	}

	if n, err := sctpConn.Write(pkt); err != nil {
		log.Errorf("Write to SCTP socket failed: %+v", err)
	} else {
		log.Tracef("Wrote %d bytes", n)
	}
}

func SendUERadioCapabilityInfoIndication() {
	log.Infoln("[N3IWF] Send UE Radio Capability Info Indication")
}

func SendUERadioCapabilityCheckResponse(
	amf *context.N3IWFAMF,
	ue *context.N3IWFUe,
	diagnostics *ngapType.CriticalityDiagnostics) {
	log.Infoln("[N3IWF] Send UE Radio Capability Check Response")

	pkt, err := BuildUERadioCapabilityCheckResponse(ue, diagnostics)
	if err != nil {

		log.Errorf("Build UERadio Capability Check Response failed : %+v\n", err)
		return
	}
	SendToAmf(amf, pkt)
}

func SendAMFConfigurationUpdateAcknowledge(
	amf *context.N3IWFAMF,
	setupList *ngapType.AMFTNLAssociationSetupList,
	failList *ngapType.TNLAssociationList,
	diagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send AMF Configuration Update Acknowledge")

	pkt, err := BuildAMFConfigurationUpdateAcknowledge(setupList, failList, diagnostics)
	if err != nil {
		log.Errorf("Build AMF Configuration Update Acknowledge failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendAMFConfigurationUpdateFailure(
	amf *context.N3IWFAMF,
	ngCause ngapType.Cause,
	time *ngapType.TimeToWait,
	diagnostics *ngapType.CriticalityDiagnostics) {

	log.Infoln("[N3IWF] Send AMF Configuration Update Failure")
	pkt, err := BuildAMFConfigurationUpdateFailure(ngCause, time, diagnostics)
	if err != nil {
		log.Errorf("Build AMF Configuration Update Failure failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendRANConfigurationUpdate(amf *context.N3IWFAMF) {

	log.Infoln("[N3IWF] Send RAN Configuration Update")

	if available, _ := context.N3IWFSelf().AMFReInitAvailableListLoad(amf.SCTPAddr); !available {
		log.Warnf(
			"[N3IWF] Please Wait at least for the indicated time before reinitiating toward same AMF[%s]", amf.SCTPAddr)
		return
	}

	pkt, err := BuildRANConfigurationUpdate()
	if err != nil {
		log.Errorf("Build AMF Configuration Update Failure failed : %+v\n", err)
		return
	}

	SendToAmf(amf, pkt)
}

func SendUplinkRANConfigurationTransfer() {
	log.Infoln("[N3IWF] Send Uplink RAN Configuration Transfer")
}

func SendUplinkRANStatusTransfer() {
	log.Infoln("[N3IWF] Send Uplink RAN Status Transfer")
}

func SendLocationReportingFailureIndication() {
	log.Infoln("[N3IWF] Send Location Reporting Failure Indication")
}

func SendLocationReport() {
	log.Infoln("[N3IWF] Send Location Report")
}

func SendRRCInactiveTransitionReport() {
	log.Infoln("[N3IWF] Send RRC Inactive Transition Report")
}
