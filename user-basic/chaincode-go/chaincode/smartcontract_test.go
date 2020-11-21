package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	userBasic := chaincode.SmartContract{}
	err := userBasic.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = userBasic.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to the ledger. failed inserting key")
}

func TestRegisterUser(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	userBasic := chaincode.SmartContract{}
	err := userBasic.RegisterUser(transactionContext, "", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = userBasic.RegisterUser(transactionContext, "Dennis", 0)
	require.EqualError(t, err, "the user Dennis already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve user"))
	err = userBasic.RegisterUser(transactionContext, "Dennis", 0)
	require.EqualError(t, err, "failed to read from the ledger: unable to retrieve user")
}

func TestGetUser(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedUser := &chaincode.BasicUser{ID: "Dennis"}
	bytes, err := json.Marshal(expectedUser)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	userBasic := chaincode.SmartContract{}
	user, err := userBasic.GetUser(transactionContext, "")
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve user"))
	_, err = userBasic.GetUser(transactionContext, "")
	require.EqualError(t, err, "failed to read from the ledger: unable to retrieve user")

	chaincodeStub.GetStateReturns(nil, nil)
	user, err = userBasic.GetUser(transactionContext, "Dennis")
	require.EqualError(t, err, "the user Dennis does not exist")
	require.Nil(t, user)
}

/*
func TestUpdateUser(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedUser := &chaincode.BasicUser{ID: "Dennis"}
	bytes, err := json.Marshal(expectedUser)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	userBasic := chaincode.SmartContract{}
	err = userBasic.UpdateUser(transactionContext, "", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = userBasic.UpdateUser(transactionContext, "Dennis", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve user"))
	err = userBasic.UpdateUser(transactionContext, "Dennis", 0)
	require.EqualError(t, err, "failed to read from the ledger: unable to retrieve user")
}
*/

func TestDeleteUser(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	user := &chaincode.BasicUser{ID: "Dennis"}
	bytes, err := json.Marshal(user)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	userBasic := chaincode.SmartContract{}
	err = userBasic.DeleteUser(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = userBasic.DeleteUser(transactionContext, "Dennis")
	require.EqualError(t, err, "the user Dennis does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve user"))
	err = userBasic.DeleteUser(transactionContext, "")
	require.EqualError(t, err, "failed to read from the ledger: unable to retrieve user")
}

func TestGetAllUsers(t *testing.T) {
	user := &chaincode.BasicUser{ID: "Dennis"}
	bytes, err := json.Marshal(user)
	require.NoError(t, err)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	userBasic := &chaincode.SmartContract{}
	users, err := userBasic.GetAllUsers(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*chaincode.BasicUser{user}, users)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next user"))
	users, err = userBasic.GetAllUsers(transactionContext)
	require.EqualError(t, err, "failed retrieving next user")
	require.Nil(t, users)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all users"))
	users, err = userBasic.GetAllUsers(transactionContext)
	require.EqualError(t, err, "failed retrieving all users")
	require.Nil(t, users)
}
