package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a AcronicUser
type SmartContract struct {
	contractapi.Contract
}

// AcronicUser describes basic details of what makes up a simple user
type AcronicUser struct {
	ID             string `json:"ID"`
	Access_lvl     int    `json:"accessLevel"` // Access Level
}

// InitLedger adds a base set of users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []AcronicUser{
		{ID: "Dennis", Access_lvl: 0},
		{ID: "Mark", Access_lvl: 88},
		{ID: "Sanjeev", Access_lvl: 88},
	}

	for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(user.ID, userJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// RegisterUser issues a new user to the world state with given details.
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, id string, accessLevel int) error {
	exists, err := s.AcronicUserExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", id)
	}

	if accessLevel < 0 || accessLevel > 100 {
		return fmt.Errorf("Invalid access level %v. Range is 0..100!", accessLevel)
	}

	user := AcronicUser{
		ID:             id,
		Access_lvl:     accessLevel, // Access level
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, userJSON)
}

// GetUser returns the user stored in the world state with given id.
func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*AcronicUser, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", id)
	}

	var user AcronicUser
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user in the world state with provided parameters.
func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, accessLevel int) error {
	exists, err := s.AcronicUserExists(ctx, id)
	if err != nil {
		return err
	}

	if accessLevel < 0 || accessLevel > 100 {
		return fmt.Errorf("Invalid access level %v. Range is 0..100!", accessLevel)
	}

	// Create a new user if not exist
	if !exists {
		return s.RegisterUser(ctx, id, accessLevel)
	}

	// overwriting original user with new user
	user := AcronicUser{
		ID:             id,
		Access_lvl:     accessLevel,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, userJSON)
}

// DeleteUser deletes an given user from the world state.
func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AcronicUserExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AcronicUserExists returns true when user with given ID exists in world state
func (s *SmartContract) AcronicUserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return userJSON != nil, nil
}

// GetAllUsers returns all users found in world state
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*AcronicUser, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all users in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*AcronicUser
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user AcronicUser
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
