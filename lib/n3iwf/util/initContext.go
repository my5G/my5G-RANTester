package util

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net"
	"strings"

	"git.cs.nctu.edu.tw/calee/sctp"

	"my5G-RANTester/lib/n3iwf/context"
	"my5G-RANTester/lib/n3iwf/factory"
	"my5G-RANTester/lib/path_util"

	log "github.com/sirupsen/logrus"
)

func InitN3IWFContext() bool {
	var ok bool

	if factory.N3iwfConfig.Configuration == nil {
		log.Error("No N3IWF configuration found")
		return false
	}

	n3iwfContext := context.N3IWFSelf()

	// N3IWF NF information
	n3iwfContext.NFInfo = factory.N3iwfConfig.Configuration.N3IWFInfo
	if ok = formatSupportedTAList(&n3iwfContext.NFInfo); !ok {
		return false
	}

	// AMF SCTP addresses
	if len(factory.N3iwfConfig.Configuration.AMFSCTPAddresses) == 0 {
		log.Error("No AMF specified")
		return false
	} else {
		for _, amfAddress := range factory.N3iwfConfig.Configuration.AMFSCTPAddresses {
			amfSCTPAddr := new(sctp.SCTPAddr)
			// IP addresses
			for _, ipAddrStr := range amfAddress.IPAddresses {
				if ipAddr, err := net.ResolveIPAddr("ip", ipAddrStr); err != nil {
					log.Errorf("Resolve AMF IP address failed: %+v", err)
					return false
				} else {
					amfSCTPAddr.IPAddrs = append(amfSCTPAddr.IPAddrs, *ipAddr)
				}
			}
			// Port
			if amfAddress.Port == 0 {
				amfSCTPAddr.Port = 38412
			} else {
				amfSCTPAddr.Port = amfAddress.Port
			}
			// Append to context
			n3iwfContext.AMFSCTPAddresses = append(n3iwfContext.AMFSCTPAddresses, amfSCTPAddr)
		}
	}

	// IKE bind address
	if factory.N3iwfConfig.Configuration.IKEBindAddr == "" {
		log.Error("IKE bind address is empty")
		return false
	} else {
		n3iwfContext.IKEBindAddress = factory.N3iwfConfig.Configuration.IKEBindAddr
	}

	// IPSec gateway address
	if factory.N3iwfConfig.Configuration.IPSecGatewayAddr == "" {
		log.Error("IPSec interface address is empty")
		return false
	} else {
		n3iwfContext.IPSecGatewayAddress = factory.N3iwfConfig.Configuration.IPSecGatewayAddr
	}

	// GTP bind address
	if factory.N3iwfConfig.Configuration.GTPBindAddr == "" {
		log.Error("GTP bind address is empty")
		return false
	} else {
		n3iwfContext.GTPBindAddress = factory.N3iwfConfig.Configuration.GTPBindAddr
	}

	// TCP port
	if factory.N3iwfConfig.Configuration.TCPPort == 0 {
		log.Error("TCP port is not defined")
		return false
	} else {
		n3iwfContext.TCPPort = factory.N3iwfConfig.Configuration.TCPPort
	}

	// FQDN
	if factory.N3iwfConfig.Configuration.FQDN == "" {
		log.Error("FQDN is empty")
		return false
	} else {
		n3iwfContext.FQDN = factory.N3iwfConfig.Configuration.FQDN
	}

	// Private key
	{
		var keyPath string

		if factory.N3iwfConfig.Configuration.PrivateKey == "" {
			log.Warn("No private key file path specified, load default key file...")
			keyPath = path_util.Gofree5gcPath("free5gc/support/TLS/n3iwf.key")
		} else {
			keyPath = factory.N3iwfConfig.Configuration.PrivateKey
		}

		content, err := ioutil.ReadFile(keyPath)
		if err != nil {
			log.Errorf("Cannot read private key data from file: %+v", err)
			return false
		}
		block, _ := pem.Decode(content)
		if block == nil {
			log.Error("Parse pem failed")
			return false
		}
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Warnf("Parse PKCS8 private key failed: %+v", err)
			log.Info("Parse using PKCS1...")

			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				log.Errorf("Parse PKCS1 pricate key failed: %+v", err)
				return false
			}
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			log.Error("Private key is not an rsa private key")
			return false
		}

		n3iwfContext.N3IWFPrivateKey = rsaKey
	}

	// Certificate authority
	{
		var keyPath string

		if factory.N3iwfConfig.Configuration.CertificateAuthority == "" {
			log.Warn("No certificate authority file path specified, load default CA certificate...")
			keyPath = path_util.Gofree5gcPath("free5gc/support/TLS/n3iwf.pem")
		} else {
			keyPath = factory.N3iwfConfig.Configuration.CertificateAuthority
		}

		// Read .pem
		content, err := ioutil.ReadFile(keyPath)
		if err != nil {
			log.Errorf("Cannot read certificate authority data from file: %+v", err)
			return false
		}
		// Decode pem
		block, _ := pem.Decode(content)
		if block == nil {
			log.Error("Parse pem failed")
			return false
		}
		// Parse DER-encoded x509 certificate
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Errorf("Parse certificate authority failed: %+v", err)
			return false
		}
		// Get sha1 hash of subject public key info
		sha1Hash := sha1.New()
		if _, err := sha1Hash.Write(cert.RawSubjectPublicKeyInfo); err != nil {
			log.Errorf("Hash function writing failed: %+v", err)
			return false
		}

		n3iwfContext.CertificateAuthority = sha1Hash.Sum(nil)
	}

	// Certificate
	{
		var keyPath string

		if factory.N3iwfConfig.Configuration.Certificate == "" {
			log.Warn("No certificate file path specified, load default certificate...")
			keyPath = path_util.Gofree5gcPath("free5gc/support/TLS/n3iwf.pem")
		} else {
			keyPath = factory.N3iwfConfig.Configuration.Certificate
		}

		// Read .pem
		content, err := ioutil.ReadFile(keyPath)
		if err != nil {
			log.Errorf("Cannot read certificate data from file: %+v", err)
			return false
		}
		// Decode pem
		block, _ := pem.Decode(content)
		if block == nil {
			log.Error("Parse pem failed")
			return false
		}

		n3iwfContext.N3IWFCertificate = block.Bytes
	}

	// UE IP address range
	if factory.N3iwfConfig.Configuration.UEIPAddressRange == "" {
		log.Error("UE IP address range is empty")
		return false
	} else {
		_, ueIPRange, err := net.ParseCIDR(factory.N3iwfConfig.Configuration.UEIPAddressRange)
		if err != nil {
			log.Errorf("Parse CIDR failed: %+v", err)
			return false
		}
		n3iwfContext.Subnet = ueIPRange
	}

	if factory.N3iwfConfig.Configuration.InterfaceMark == 0 {
		log.Warn("IPSec interface mark is not defined, set to default value 7")
		n3iwfContext.Mark = 7
	} else {
		n3iwfContext.Mark = factory.N3iwfConfig.Configuration.InterfaceMark
	}

	return true
}

func formatSupportedTAList(info *context.N3IWFNFInfo) bool {
	for taListIndex := range info.SupportedTAList {

		supportedTAItem := &info.SupportedTAList[taListIndex]

		// Checking TAC
		if supportedTAItem.TAC == "" {
			log.Error("TAC is mandatory.")
			return false
		}
		if len(supportedTAItem.TAC) < 6 {
			log.Trace("Detect configuration TAC length < 6")
			supportedTAItem.TAC = strings.Repeat("0", 6-len(supportedTAItem.TAC)) + supportedTAItem.TAC
			log.Tracef("Changed to %s", supportedTAItem.TAC)
		} else if len(supportedTAItem.TAC) > 6 {
			log.Error("Detect configuration TAC length > 6")
			return false
		}

		// Checking SST and SD
		for plmnListIndex := range supportedTAItem.BroadcastPLMNList {

			broadcastPLMNItem := &supportedTAItem.BroadcastPLMNList[plmnListIndex]

			for sliceListIndex := range broadcastPLMNItem.TAISliceSupportList {

				sliceSupportItem := &broadcastPLMNItem.TAISliceSupportList[sliceListIndex]

				// SST
				if sliceSupportItem.SNSSAI.SST == "" {
					log.Error("SST is mandatory.")
				}
				if len(sliceSupportItem.SNSSAI.SST) < 2 {
					log.Trace("Detect configuration SST length < 2")
					sliceSupportItem.SNSSAI.SST = "0" + sliceSupportItem.SNSSAI.SST
					log.Tracef("Change to %s", sliceSupportItem.SNSSAI.SST)
				} else if len(sliceSupportItem.SNSSAI.SST) > 2 {
					log.Error("Detect configuration SST length > 2")
					return false
				}

				// SD
				if sliceSupportItem.SNSSAI.SD != "" {
					if len(sliceSupportItem.SNSSAI.SD) < 6 {
						log.Trace("Detect configuration SD length < 6")
						sliceSupportItem.SNSSAI.SD = strings.Repeat("0", 6-len(sliceSupportItem.SNSSAI.SD)) + sliceSupportItem.SNSSAI.SD
						log.Tracef("Change to %s", sliceSupportItem.SNSSAI.SD)
					} else if len(sliceSupportItem.SNSSAI.SD) > 6 {
						log.Error("Detect configuration SD length > 6")
						return false
					}
				}

			}
		}

	}

	return true
}
