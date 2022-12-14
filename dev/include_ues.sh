#!/bin/bash

HOST=192.168.2.162
PORT=5000

BASE_URL="http://$HOST:$PORT/api/subscriber"
K=8baf473f2f8fd09487cccbd7097c6862
OPC=8e27b6af0e692e750f32667a3b14605d

i=0

imsi=208930000000000
plmnid=20893

while [ $i -lt $1 ] ;
do
    imsi=$((imsi+1))
    curl "$BASE_URL/imsi-$imsi/$plmnid" \
        -H 'Accept: application/json' \
        -H 'Accept-Language: en-US,en;q=0.9,pt;q=0.8' \
        -H 'Content-Type: application/json;charset=UTF-8' \
        -H 'Token: admin' \
        -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36' \
        -H 'X-Requested-With: XMLHttpRequest' \
        --data-raw '{"plmnID":"'$plmnid'","ueId":"imsi-'$imsi'","AuthenticationSubscription":{"authenticationManagementField":"8000","authenticationMethod":"5G_AKA","milenage":{"op":{"encryptionAlgorithm":0,"encryptionKey":0,"opValue":""}},"opc":{"encryptionAlgorithm":0,"encryptionKey":0,"opcValue":"'$OPC'"},"permanentKey":{"encryptionAlgorithm":0,"encryptionKey":0,"permanentKeyValue":"'$K'"},"sequenceNumber":"16f3b3f70fc2"},"AccessAndMobilitySubscriptionData":{"gpsis":["msisdn-0900000000"],"nssai":{"defaultSingleNssais":[{"sst":1,"sd":"010203","isDefault":true},{"sst":1,"sd":"112233","isDefault":true}],"singleNssais":[]},"subscribedUeAmbr":{"downlink":"2 Gbps","uplink":"1 Gbps"}},"SessionManagementSubscriptionData":[{"singleNssai":{"sst":1,"sd":"010203"},"dnnConfigurations":{"internet":{"sscModes":{"defaultSscMode":"SSC_MODE_1","allowedSscModes":["SSC_MODE_2","SSC_MODE_3"]},"pduSessionTypes":{"defaultSessionType":"IPV4","allowedSessionTypes":["IPV4"]},"sessionAmbr":{"uplink":"200 Mbps","downlink":"100 Mbps"},"5gQosProfile":{"5qi":9,"arp":{"priorityLevel":8},"priorityLevel":8}},"internet2":{"sscModes":{"defaultSscMode":"SSC_MODE_1","allowedSscModes":["SSC_MODE_2","SSC_MODE_3"]},"pduSessionTypes":{"defaultSessionType":"IPV4","allowedSessionTypes":["IPV4"]},"sessionAmbr":{"uplink":"200 Mbps","downlink":"100 Mbps"},"5gQosProfile":{"5qi":9,"arp":{"priorityLevel":8},"priorityLevel":8}}}},{"singleNssai":{"sst":1,"sd":"112233"},"dnnConfigurations":{"internet":{"sscModes":{"defaultSscMode":"SSC_MODE_1","allowedSscModes":["SSC_MODE_2","SSC_MODE_3"]},"pduSessionTypes":{"defaultSessionType":"IPV4","allowedSessionTypes":["IPV4"]},"sessionAmbr":{"uplink":"200 Mbps","downlink":"100 Mbps"},"5gQosProfile":{"5qi":9,"arp":{"priorityLevel":8},"priorityLevel":8}},"internet2":{"sscModes":{"defaultSscMode":"SSC_MODE_1","allowedSscModes":["SSC_MODE_2","SSC_MODE_3"]},"pduSessionTypes":{"defaultSessionType":"IPV4","allowedSessionTypes":["IPV4"]},"sessionAmbr":{"uplink":"200 Mbps","downlink":"100 Mbps"},"5gQosProfile":{"5qi":9,"arp":{"priorityLevel":8},"priorityLevel":8}}}}],"SmfSelectionSubscriptionData":{"subscribedSnssaiInfos":{"01010203":{"dnnInfos":[{"dnn":"internet"},{"dnn":"internet2"}]},"01112233":{"dnnInfos":[{"dnn":"internet"},{"dnn":"internet2"}]}}},"AmPolicyData":{"subscCats":["free5gc"]},"SmPolicyData":{"smPolicySnssaiData":{"01010203":{"snssai":{"sst":1,"sd":"010203"},"smPolicyDnnData":{"internet":{"dnn":"internet"},"internet2":{"dnn":"internet2"}}},"01112233":{"snssai":{"sst":1,"sd":"112233"},"smPolicyDnnData":{"internet":{"dnn":"internet"},"internet2":{"dnn":"internet2"}}}}},"FlowRules":[]}' \
        --compressed \
        --insecure

    i=$((i+1))
done