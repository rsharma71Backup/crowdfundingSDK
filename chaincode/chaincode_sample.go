/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var userIndexStr = "_userindex"
var campaignIndexStr = "_campaignindex"
var transactionIndexStr = "_transactionindex"

type User struct {
	Name           string `json:"name"` //the fieldtags of user are needed to store in the ledger
	Id             int    `json:"id"`
	Phone          int    `json:"phone"`
	Email          string `json:"email"`
	User_Type      string `json:"usertype"`
	Password       string `json:"password"`
	Reviews_Count  int    `json:"reviewcount"`
	Average_Rating int    `json:"averagerating"`
	Pan_No         string `json:"panno"`
	Aadhar_id      int    `json:"aadharid"`
}
type Campaign struct {
	Campaign_Id           int    `json:"campaign_id"`
	Campaign_Title        string `json:"campaign_title"`
	Campaign_Description  string `json:"campaign_description"`
	Campaign_Story        string `json:"campaign_story"`
	Campaign_Goal_amount  int    `json:"campaign_goal_amount"`
	Campaign_Start_amount int    `json:"campaign_start_amount"`
	Campaign_Start_Date   string `json:"campaign_start_date"`
	Campaign_End_Date     string `json:"campaign_end_date"`
	Project_Start_Date    string `json:"project_start_date"`
	Project_End_Date      string `json:"project_end_date"`
	Rewards               string `json:"rewards"`
}
type Transaction struct {
	Transaction_Id    string   `json:"transaction_id"`
	Amount_Transfered int      `json:"amount_transfered"`
	User              User     `json:"user"`
	Campaign          Campaign `json:"campaign"`
}

/* type Rewards struct{
Campaign_Id Campaign `json:"campaign"`
RewardTitle string `json:"rewardtitle"`
RewardDescription string `json:"rewarddescription"`
}
/*type Login struct{
Emailid string `json:"emailid"`
UserPassword string `json:"userpassword"`
}*/

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreateLoanApplication")
	return nil, nil
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//_, args := stub.GetFunctionAndParameters()
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(userIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)

	} else if function == "User_register" {
		return t.User_register(stub, args)

	} else if function == "Delete" {
		return t.Delete(stub, args)

	} else if function == "create_campaign" {
		return t.create_campaign(stub, args)

	} else if function == "transactions" {
		return t.transactions(stub, args)

	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "readuser" { //read a variable
		return t.readuser(stub, args)
	} else if function == "readcampaign" {
		return t.readcampaign(stub, args)
	} else if function == "readtxdetails" {
		return t.readtxdetails(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair

func (t *SimpleChaincode) readuser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error
	//var campaign_title,jsonResp string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}
func (t *SimpleChaincode) readcampaign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var campaign_title, jsonResp string
	var err error
	//var campaign_title,jsonResp string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	campaign_title = args[0]
	valAsbytes, err := stub.GetState(campaign_title) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + campaign_title + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}
func (t *SimpleChaincode) readtxdetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var transaction_id, jsonResp string
	var err error
	//var campaign_title,jsonResp string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	transaction_id = args[0]
	valAsbytes, err := stub.GetState(transaction_id) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + transaction_id + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}
func (t *SimpleChaincode) User_register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "lol", "1", "323323", "r@r.com"
	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	name := args[0]
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Failed to get id as cannot convert it to int")
	}
	phone, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Failed to get phone as cannot convert it to int")
	}
	email := args[3]
	usertype := args[4]
	password := args[5]
	panno := args[6]
	aadharid, err := strconv.Atoi(args[7])
	if err != nil {
		return nil, errors.New("Failed to get aadaharid as cannot convert it to int")
	}

	//check if marble already exists
	userAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get user name")
	}
	res := User{}
	json.Unmarshal(userAsBytes, &res)
	if res.Name == name {
		fmt.Println("This name arleady exists: " + name)
		fmt.Println(res)
		fmt.Println(userAsBytes)
		return nil, errors.New("This user arleady exists") //all stop a marble by this name exists
	}

	//build the marble json string manually
	str := `{"name": "` + name + `", "id": "` + strconv.Itoa(id) + `", "phone": "` + strconv.Itoa(phone) + `", "email": "` + email + `", "usertype": "` + usertype + `", "password": "` + password + `", "panno": "` + panno + `", "aadharid": "` + strconv.Itoa(aadharid) + `"}`
	err = stub.PutState(name, []byte(str)) //store marble with id as key
	fmt.Println(name)
	fmt.Println([]byte(str))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	//get the marble index
	userAsBytes, err1 := stub.GetState(userIndexStr)
	if err1 != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var userIndex []string
	json.Unmarshal(userAsBytes, &userIndex) //un stringify it aka JSON.parse()

	//append
	userIndex = append(userIndex, name) //add marble name to index list
	fmt.Println("! user index: ", userIndex)
	jsonAsBytes, _ := json.Marshal(userIndex)
	err = stub.PutState(userIndexStr, jsonAsBytes) //store name of marble

	fmt.Println("- end user_register")
	return nil, nil
}
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	err := stub.DelState(name) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the marble index
	userAsBytes, err := stub.GetState(userIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get user index")
	}
	var userIndex []string
	json.Unmarshal(userAsBytes, &userIndex) //un stringify it aka JSON.parse()

	//remove user from index
	for i, val := range userIndex {
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name { //find the correct marble
			fmt.Println("found user")
			userIndex = append(userIndex[:i], userIndex[i+1:]...) //remove it
			for x := range userIndex {                            //debug prints...
				fmt.Println(string(x) + " - " + userIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(userIndex) //save new index
	err = stub.PutState(userIndexStr, jsonAsBytes)
	return nil, nil
}

func (t *SimpleChaincode) create_campaign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "lol", "1", "323323", "r@r.com"
	if len(args) != 11 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[9]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[10]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	campaign_id, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Failed to get campaign_id as cannot convert it to int")
	}
	campaign_title := args[1]
	campaign_description := args[2]
	campaign_story := args[3]
	campaign_goal_amount, err := strconv.Atoi(args[4])
	if err != nil {
		return nil, errors.New("Failed to get id as cannot convert it to int")
	}
	campaign_start_amount, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("Failed to get phone as cannot convert it to int")
	}

	campaign_start_date := args[6]
	campaign_end_date := args[7]
	project_start_date := args[8]
	project_end_date := args[9]
	rewards := args[10]

	//check if marble already exists
	campaignAsBytes, err := stub.GetState(campaign_title)
	if err != nil {
		return nil, errors.New("Failed to get user name")
	}
	//reward:=Rewards{}
	//json.Unmarshal(rewardAsBytes, &reward)
	res := Campaign{}
	json.Unmarshal(campaignAsBytes, &res)
	if res.Campaign_Title == campaign_title {
		fmt.Println("This campaign arleady exists: " + campaign_title)
		fmt.Println(res)
		fmt.Println(campaignAsBytes)
		return nil, errors.New("This campaign arleady exists") //all stop a marble by this name exists
	}

	//build the marble json string manually
	str := `{"campaign_id": "` + strconv.Itoa(campaign_id) + `", "campaign_title": "` + campaign_title + `", "campaign_description": ` + campaign_description + `, "campaign_story": "` + campaign_story + `, "campaign_goal_amount": "` + strconv.Itoa(campaign_goal_amount) + `, "campaign_start_amount": "` + strconv.Itoa(campaign_start_amount) + `, "campaign_start_date": "` + campaign_start_date + `, "campaign_end_date": "` + campaign_end_date + `, "project_start_date": "` + project_start_date + `, "project_end_date": "` + project_end_date + `, "rewards": "` + rewards + `"}`
	err = stub.PutState(campaign_title, []byte(str)) //store marble with id as key
	fmt.Println(campaign_title)
	fmt.Println([]byte(str))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	campaignAsBytes, err1 := stub.GetState(campaignIndexStr)
	if err1 != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var camapignIndex []string

	json.Unmarshal(campaignAsBytes, &camapignIndex) //un stringify it aka JSON.parse()

	//append
	camapignIndex = append(camapignIndex, campaign_title) //add marble name to index list
	fmt.Println("! user index: ", camapignIndex)
	jsonAsBytes, _ := json.Marshal(camapignIndex)
	err = stub.PutState(campaignIndexStr, jsonAsBytes) //store name of marble

	fmt.Println("- end create campaign")
	return nil, nil
}
func (t *SimpleChaincode) transactions(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
	var user_id int
	//input sanitation
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	transaction_id := args[0]
	amount_transfered, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Failed to get campaign_id as cannot convert it to int")
	}
	name := args[2]
	fmt.Println(name)
	userAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get user name")
	}

	fmt.Println("userasbytes", userAsBytes)
	var res User
	json.Unmarshal(userAsBytes, &res)

	fmt.Printf("%+v", res)
	fmt.Println(userAsBytes)
	user_id = res.Id
	fmt.Println("user_id", user_id)
	campaign_title := args[3]
	fmt.Println(campaign_title)
	campaignAsBytes, err := stub.GetState(campaign_title)
	if err != nil {
		return nil, errors.New("Failed to get user name")
	}

	fmt.Println("campaignasbytes", campaignAsBytes)
	res1 := Campaign{}
	json.Unmarshal(campaignAsBytes, &res1)
	campaign_id := res1.Campaign_Id
	fmt.Println("campaign_id", campaign_id)

	//build the marble json string manually
	str := `{"transaction_id": "` + transaction_id + `", "amount_transfered": "` + strconv.Itoa(amount_transfered) + `", "user_id": ` + strconv.Itoa(user_id) + `, "campaign_id": "` + strconv.Itoa(campaign_id) + `"}`
	err = stub.PutState(transaction_id, []byte(str)) //store marble with id as key
	fmt.Println(transaction_id)
	fmt.Println([]byte(str))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	tansactionAsBytes, err1 := stub.GetState(transactionIndexStr)
	if err1 != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var transactionIndex []string

	json.Unmarshal(tansactionAsBytes, &transactionIndex) //un stringify it aka JSON.parse()

	//append
	transactionIndex = append(transactionIndex, transaction_id) //add marble name to index list
	fmt.Println("! user index: ", transactionIndex)
	jsonAsBytes, _ := json.Marshal(transactionIndex)
	err = stub.PutState(transactionIndexStr, jsonAsBytes) //store name of marble

	fmt.Println("- end create campaign")
	return nil, nil
}
