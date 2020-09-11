# my5G-RANTester

## Running in docker container
**Build**

``sudo docker build -t my5g-rantester .``

**Running**

``sudo docker run --rm -e AMF_IP=10.100.200.101 -e RAN_IP=10.100.200.50 -e UPF_IP=10.100.200.223 --network some_network --ip 10.100.200.50 my5g-rantester AttachGnb``
