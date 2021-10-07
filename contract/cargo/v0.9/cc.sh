#!/bin/bash

set -ev

CHNAME=mychannel
CCNAME=cargo
VER=0.9


#chaincode install
docker exec cli peer chaincode install -n $CCNAME -v $VER -p github.com/cargo/v0.9
#chaincode instatiate
docker exec cli peer chaincode instantiate -n $CCNAME -v $VER -C $CHNAME -c '{"Args":[]}' -P 'OR ("Org1MSP.member","Org2MSP.member","Org3MSP.member")'
sleep 3
#chaincode test


docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["regPERSON", "sumin", "a/c111", "5000"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["regPERSON", "dr.choi", "a/c222", "5000"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["reqTransportation","TC0", "a/c111","5","Seo0000ul","Busan","1000","10"]}'
sleep 3
docker exec cli peer chaincode query -n $CCNAME -C $CHNAME -c '{"Args":["queryReqTransportation", "TC0"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["regTransportation", "TC1", "a/c222",  "5", "Seoul", "Busan", "1000", "10"]}'
sleep 3
docker exec cli peer chaincode query -n $CCNAME -C $CHNAME -c '{"Args":["queryRegTransportation", "TC1"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["respond", "TC0", "D", "a/c222"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["confirmContract", "TC0", "D", "a/c222"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["load", "TC0"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["depart", "TC0"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["arrive", "TC0"]}'
sleep 3
docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["pay", "TC0"]}'
sleep 3


docker exec cli peer chaincode query -n $CCNAME -C $CHNAME -c '{"Args":["history","TC0"]}'
docker exec cli peer chaincode query -n $CCNAME -C $CHNAME -c '{"Args":["history","TC1"]}'

docker exec cli peer chaincode query -n cargo -C mychannel -c '{"Args":["history","a/c111"]}'

echo '-------------------------------------END-------------------------------------'
