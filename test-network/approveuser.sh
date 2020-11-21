#!/bin/bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

if [[ $# -ne 1 ]]; then
    echo "Need 1 parameters: ./approveuser.sh Authority"
    echo "- Authority   : Either org1 or org2"
    echo "  This script will approve user from either 'org1/org2' to update the ledger."
    exit 2
fi

if [[ "$1" = "org1" ]]
then
    echo "Approve user from 'org1' to update the ledger"
    export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/

    fabric-ca-client identity modify  org1admin --attrs 'role=approval:ecert' --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem

    fabric-ca-client enroll -u https://org1admin:org1adminpw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem    

elif [[ "$1" = "org2" ]]
then
    echo "Approve user from 'org2' to update the ledger"
    export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org2.example.com/

    fabric-ca-client identity modify  org2admin --attrs 'role=approval:ecert' --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem

    fabric-ca-client enroll -u https://org2admin:org2adminpw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem    
else
    echo "Authority can only be org1 or org2."
    exit 2
fi
