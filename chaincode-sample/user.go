package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type User struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
	Activity string `json:"activity"`
	CreatedAt string `json:"created_at"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, userID string, id string, activity string) error {
	userExists, err := s.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}
	if userExists {
		return fmt.Errorf("user with userID %s already exists", userID)
	}
	createdAt := time.Now().Format(time.RFC3339)
	user := User {
		UserID: userID,
		ID: id,
		Activity: activity,
		CreatedAt: createdAt,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}
	return ctx.GetStub().PutState(userID, userJSON)
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, userID string) (bool, error) {
	userJSON, err := ctx. GetStub().GetState(userID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return userJSON != nil, nil
}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("user %s does not exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}
	return &user, nil
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, userID string, newActivity string) error {
    user, err := s.GetUser(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %v", err)
    }

    user.Activity = newActivity
    user.CreatedAt = time.Now().Format(time.RFC3339)

    userJSON, err := json.Marshal(user)
    if err != nil {
        return fmt.Errorf("failed to marshal updated user: %v", err)
    }

    return ctx.GetStub().PutState(userID, userJSON)
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, userID string) error {
	userExists, err := s.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}
	if !userExists {
        return fmt.Errorf("user with userID %s does not exist", userID)
    }
	return ctx.GetStub().DelState(userID)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating user management chaincode: %v\n", err)
		return
	}
	if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting user management chaincode: %v\n", err)
    }
}