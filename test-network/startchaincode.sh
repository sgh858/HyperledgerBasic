
source scriptUtils.sh

#Start chaincode
cd ~/acronics/test-network/
./network.sh up
./network.sh createChannel
./network.sh deployCC

#Interacting with the network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

#Run the following command to initialize the ledger with simple users:
echo "**********************************************"
echo "*       Chaincode BasicUser started          *"
echo "*     Initialized ledger with basic users    *"
echo "**********************************************"
echo ""
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basicuser --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

echo "Chaincode initialized. Please enter your command..."
echo ""
echo "You're free to use the following commands:"

println "${C_GREEN}./getallusers.sh${C_RESET} : No parameter, display all user in the ledger"
println "${C_GREEN}./getusers.sh Name${C_RESET} : 1 parameter, display user with Name in the ledger"
println "${C_GREEN}./deleteuser.sh Name${C_RESET} : 1 parameter, delete user with Name in the ledger"
println "${C_GREEN}./registeruser.sh Name accessLevel${C_RESET} : 2 parameters, create a new user with given params in the ledger"
println "${C_GREEN}./updateuser.sh Name accessLevel${C_RESET} : 2 parameters, update or create a new user with given params  in the ledger"
echo ""
