LOCAL_UE_IP=10.0.2.5
N3IWF_IP=10.100.200.14

ip route add 10.100.200.0/24 via 10.0.2.4

ip link add ipsec-default type vti local $LOCAL_UE_IP remote $N3IWF_IP key 5
ip link set ipsec-default up

