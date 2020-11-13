# compile and pack the new chaincode
cd ~/acronics/user-basic/chaincode-go
GO111MODULE=on go mod vendor
cd ../test-network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
peer version
peer lifecycle chaincode package basicuser.tar.gz --path ../user-basic/chaincode-go/ --lang golang --label basicuser_1.0
