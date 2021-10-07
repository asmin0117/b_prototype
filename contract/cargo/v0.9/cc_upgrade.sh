!/bin/bash


CHNAME=mychannel
CCNAME=cargo
VER=$1


#chaincode insall
docker exec cli peer chaincode install -n $CCNAME -v $VER -p github.com/cargo/v0.9
#chaincode instatiate
docker exec cli peer chaincode upgrade -n $CCNAME -v $VER -C $CHNAME -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'
sleep 3
#chaincode test

docker exec cli peer chaincode invoke -n $CCNAME -C $CHNAME -c '{"Args":["pay", "TC0"]}'
sleep 3


docker exec cli peer chaincode query -n $CCNAME -C $CHNAME -c '{"Args":["queryResTransportation","TC0"]}'
docker exec cli peer chaincode query -n $CCNAME -C $CHNAME -c '{"Args":["queryReqTransportation","TC1"]}'

echo '-------------------------------------END-------------------------------------'
