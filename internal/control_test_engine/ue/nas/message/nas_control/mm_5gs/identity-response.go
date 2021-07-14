package mm_5gs

import (
	"bytes"
	"fmt"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
)

func IdentityResponse(identity5gs string, ue *context.UEContext) (nasPdu []byte) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeIdentityResponse)

	identityResponse := nasMessage.NewIdentityResponse(0)
	identityResponse.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	identityResponse.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	identityResponse.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0x00)
	identityResponse.IdentityResponseMessageIdentity.SetMessageType(nas.MsgTypeIdentityResponse)
	identityResponse.MobileIdentity5GS = ue.GetMobileIdentity(identity5gs)

	m.GmmMessage.IdentityResponse = identityResponse

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	nasPdu = data.Bytes()
	return
}
