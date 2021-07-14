package nasMessage

import (
	"bytes"
	"encoding/binary"
	"my5G-RANTester/lib/nas/nasType"
)

type IdentityResponse struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.IdentityResponseMessageIdentity
	nasType.MobileIdentity5GS
}

func NewIdentityResponse(iei uint8) (identityResponse *IdentityResponse) {
	identityResponse = &IdentityResponse{}
	return identityResponse
}

func (a *IdentityResponse) EncodeIdentityResponse(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet)
	binary.Write(buffer, binary.BigEndian, &a.SpareHalfOctetAndSecurityHeaderType.Octet)
	binary.Write(buffer, binary.BigEndian, &a.IdentityResponseMessageIdentity.Octet)
	binary.Write(buffer, binary.BigEndian, a.MobileIdentity5GS.GetLen())
	binary.Write(buffer, binary.BigEndian, &a.MobileIdentity5GS.Buffer)
}

func (a *IdentityResponse) DecodeIdentityResponse(byteArray *[]byte) {
	buffer := bytes.NewBuffer(*byteArray)
	binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet)
	binary.Read(buffer, binary.BigEndian, &a.SpareHalfOctetAndSecurityHeaderType.Octet)
	binary.Read(buffer, binary.BigEndian, &a.IdentityResponseMessageIdentity.Octet)
	binary.Read(buffer, binary.BigEndian, &a.MobileIdentity5GS.Len)
	a.MobileIdentity5GS.SetLen(a.MobileIdentity5GS.GetLen())
	binary.Read(buffer, binary.BigEndian, &a.MobileIdentity5GS.Buffer)
	for buffer.Len() > 0 {
		var ieiN uint8
		var tmpIeiN uint8
		binary.Read(buffer, binary.BigEndian, &ieiN)
		// fmt.Println(ieiN)
		if ieiN >= 0x80 {
			tmpIeiN = (ieiN & 0xf0) >> 4
		} else {
			tmpIeiN = ieiN
		}
		// fmt.Println("type", tmpIeiN)
		switch tmpIeiN {
		default:
		}
	}
}
