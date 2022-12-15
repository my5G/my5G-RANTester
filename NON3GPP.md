# my5G-RANTester (non3gpp feature)
The non3gpp feature of the my5G-RANTester enables the emulation of a non3gpp UE with the N3IWF function of the core.

## Experimental environment
For this development, it was used two virtual machines using Ubuntu 20.04.5 LTS with kernel 5.4.0-135-generic. Both VMs were deployed in the same host using the same network. The first VM was used for the deployment of the UE and the second was used for the 5G Core.
Two scenarios was tested for implementation validation. First we deployed the 5G Core directly on the VM OS. Secondly, we used the docker-compose deployment of the 5G core. The core utilized for this experiments was the [free5gc](https://github.com/free5gc/free5gc).

## Deployment
### 5G Core
Run the following commands at the 5G core VM:
```bash
git clone https://github.com/free5gc/free5gc-compose
cd free5gc-compose
docker-compose up -d

# Enable forward between docker network and VM network
DOCKER_IF_NAME=br-free5gc
IF_NAME=enp0s3

sysctl -w net.ipv4.ip_forward=1
iptables -A FORWARD -i $IF_NAME -o $DOCKER_IF_NAME -j ACCEPT
iptables -A FORWARD -i $DOCKER_IF_NAME -o $IF_NAME -j ACCEPT

# Follow 5g core logs
docker-compose logs -f
```

### Non-3GPP UE
Run the following commands at the UE VM

```bash
git clone -b non3gpp https://github.com/my5G/my5G-RANTester
cd my5G-RANTester
cd dev

# Edit IP of the 5G Core host
./include_non3gpp_ue.sh

# Set up IPSEC tunnel
sudo ./set_ipsec.sh

# There statics IPs that needs to be changed in the files
# ./my5G-RANTester/internal/control_test_engine/non3gppue/non3gppue.go
cd ../cmd
go build app.go
./app non3gpp-ue
```

## Current Status
In this current status of the implementation, we can complete the IKE negotiation but we still cannot start the communication through the IPSEC tunnel. It is still under investigation.

