package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract provides functions for managing a BasicUser
type SmartContract struct {
	contractapi.Contract
}

// BasicUser describes basic details of what makes up a simple user
type BasicUser struct {
	ID             string `json:"ID"`
	Access_lvl     int    `json:"accessLevel"` // Access Level
}

// InitLedger adds a base set of users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []BasicUser{
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
			return fmt.Errorf("failed to put to the ledger. %v", err)
		}
	}

	return nil
}

// RegisterUser issues a new user to the the ledger with given details.
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, id string, accessLevel int) error {
	exists, err := s.BasicUserExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", id)
	}

	if accessLevel < 0 || accessLevel > 100 {
		return fmt.Errorf("Invalid access level %v. Range is 0..100!", accessLevel)
	}

	user := BasicUser{
		ID:             id,
		Access_lvl:     accessLevel, // Access level
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, userJSON)
}

// GetUser returns the user stored in the the ledger with given id.
func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*BasicUser, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from the ledger: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", id)
	}

	var user BasicUser
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user in the the ledger with provided parameters.
func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, accessLevel int) error {
	exists, err := s.BasicUserExists(ctx, id)
	if err != nil {
		return err
	}

	approved, err := s.GetApproval(ctx)
	if err != nil {
		return err
	}

	if !approved {
		return fmt.Errorf("You don't have authority to use this function!") 
	}

	if accessLevel < 0 || accessLevel > 100 {
		return fmt.Errorf("Invalid access level %v. Range is 0..100!", accessLevel)
	}

	// Create a new user if not exist
	if !exists {
		return s.RegisterUser(ctx, id, accessLevel)
	}

	// overwriting original user with new user
	user := BasicUser{
		ID:             id,
		Access_lvl:     accessLevel,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, userJSON)
}

// DeleteUser deletes an given user from the the ledger.
func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.BasicUserExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// BasicUserExists returns true when user with given ID exists in the ledger
func (s *SmartContract) BasicUserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from the ledger: %v", err)
	}

	return userJSON != nil, nil
}

// GetApproval returns true when chaincode invoking user got the registered approval and role
func (s *SmartContract) GetApproval(ctx contractapi.TransactionContextInterface) (bool, error) {
	attrrole, ok, err := cid.GetAttributeValue(ctx.GetStub(), "role")
	if err != nil {
		fmt.Println("There was an error trying to retrieve the attribute.")
		return false, err 
	}
	if !ok {
   		fmt.Println("The client identity does not possess the attribute role.") 
		return false, err
	}
	if attrrole == "approval" {
		return true, nil
	} else {
		return false, nil
	}
}

// GetAllUsers returns all users found in the ledger
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*BasicUser, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all users in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*BasicUser
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user BasicUser
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
