#Interacting with the network

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

echo "**********************************************"
echo "*    Below is all users currently in the     *"
echo "*      hyperledger fabric blockchain         *"
echo "**********************************************"
echo ""
#Query the ledger from your CLI
peer chaincode query -C mychannel -n basicuser -c '{"Args":["GetAllUsers"]}'
