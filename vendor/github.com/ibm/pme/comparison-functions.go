/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package pme

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/ibm/ciav"
	"strconv"
	"strings"
)

var UPDATED_ATTRIBUTES map[string]string

/*
  input   :
  output  :
  purpose : creates the Comparison table. This table is used to store the conparison string created for each customer
            using the critical data
*/
func CreateComparisonStringTable(stub shim.ChaincodeStubInterface) error {
	myLogger.Debug("Creating Comparison string table...")
	err := stub.CreateTable("Comparison", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "dummy", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "CustomerId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ComparisonStr", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating Comparison string table")
	}
	return nil
}

/*
  input   : CustomerId , comparison string
  output  :
  purpose : Whenever a new customer is added, an entry is made for that customer in Comparison table
*/
func AddComparisonString(stub shim.ChaincodeStubInterface, CustomerId string, ComparisonStr string) error {
	myLogger.Debug("Adding Comparison string for customer : ", CustomerId, " is : ", ComparisonStr)
	ok, err := stub.InsertRow("Comparison", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: dummyValue}},
			&shim.Column{Value: &shim.Column_String_{String_: CustomerId}},
			&shim.Column{Value: &shim.Column_String_{String_: ComparisonStr}},
		},
	})
	if !ok && err == nil {
		myLogger.Debugf("ERROR : Error in inserting Comparison record.")
		return errors.New("Error in inserting Comparison record.")
	}
	myLogger.Debug("Comparison string for customer : ", CustomerId, " is : ", ComparisonStr, " added.")
	return err
}

/*
  input   : CustomerId, updated attributes
  output  : updated comparison string
  purpose : Whenever a new customer is updated, comparison string is also updated accordingly
*/
func UpdateComparisonString(stub shim.ChaincodeStubInterface, CustomerId string, updated_attrs map[string]string) (string, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: CustomerId}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("Comparison", columns)
	if err != nil {
		return "", fmt.Errorf("Failed retriving Comparison string details  for ID", err)
	}

	myLogger.Debugf("Customer Id : ", CustomerId)
	myLogger.Debugf("Before CS : ", row.Columns[2].GetString_())
	updatedCS, _ := getUpdatedComparisonString(row.Columns[2].GetString_(), updated_attrs)
	myLogger.Debugf("After CS : ", updatedCS)
	ok, err := stub.ReplaceRow("Comparison", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: dummyValue}},
			&shim.Column{Value: &shim.Column_String_{String_: CustomerId}},
			&shim.Column{Value: &shim.Column_String_{String_: updatedCS}},
		},
	})

	if !ok && err == nil {
		return "", errors.New("ERROR : while updating Comparison string record.")
	}
	return updatedCS, nil
}

/*
  input   : CustomerId
  output  : comparison string
  purpose : This function returns the comparison string for given customerId
*/
func GetComparisonString(stub shim.ChaincodeStubInterface, CustomerId string) (string, error) {
	var err error
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: CustomerId}}
	columns = append(columns, col1)
	columns = append(columns, col2)
	row, err := stub.GetRow("Comparison", columns)
	if err != nil {
		return "", fmt.Errorf("Failed retriving Comparison string for customer", CustomerId, err)
	}
	return row.Columns[2].GetString_(), nil
}

/*
  input   : comparison string
  output  : verifies if the critical data already exists
  purpose : We want to prevent adding duplicate data.
*/
func ComparisonStringExists(stub shim.ChaincodeStubInterface, data string) (bool, error) {
	myLogger.Debugf("Verifying CD for duplicate", data)
	var err error
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	columns = append(columns, col1)
	rows, err := ciav.GetAllRows(stub, "Comparison", columns)
	if err != nil {
		return false, fmt.Errorf("Failed retriving Comparison string for customer" , err)
	}
	for i := range rows {
		myLogger.Debugf("CS : ", rows[i].Columns[2].GetString_())
		if rows[i].Columns[2].GetString_() == data {
			return true, nil
		}
	}
	return false, nil
}

/*
  input   : search criteria and search keys
  output  : json string containing list of matching customers
  purpose : Get all the matching customers in
*/
func SearchRecords(stub shim.ChaincodeStubInterface, searchCriteria string, searchkeys string) string {

	bucket_key := ""
	searchCriteria = strings.Replace(searchCriteria, "|", "+", -1)
	searchkeys = strings.Replace(searchkeys, "|", "+", -1)
	sCriterias := strings.Split(searchCriteria, "+")
	sKeys := strings.Split(searchkeys, "+")
	for k,v := range sCriterias{
			attr_def := strings.Split(CRITICAL_DATA[v], "|")
			tval := StandardizeData(sKeys[k], attr_def[1])
			tval = StandardizeL2(tval, attr_def[1])
			if bucket_key != ""{
				bucket_key = bucket_key + "~"
			}
			bucket_key = bucket_key + tval
	}
	myLogger.Debugf("searchCriteria :",searchCriteria)
	myLogger.Debugf("bucket_key :",bucket_key)
	matchingCustomers, _ := SearchInBkt(stub, searchCriteria, bucket_key)

	var jsonResp string
	if !strings.Contains(searchCriteria, "PhoneNumber") {
		jsonResp = "["
		for i := range matchingCustomers {
			if i != 0 {
				jsonResp = jsonResp + ","
			}
			jsonResp = jsonResp + ciav.GetCustomerData(stub, matchingCustomers[i])
		}
		jsonResp = jsonResp + "]"
		return jsonResp
	}

	jsonResp = "["
	for i := range matchingCustomers {
			if i != 0 {
				jsonResp = jsonResp + ","
			}
			customerStr,riskLevel := ciav.GetCustomerRecord(stub, matchingCustomers[i])
			allowedActions := ciav.GetPermissionMetadata(stub, riskLevel)

			var cust ciav.Customer
			json.Unmarshal([]byte(customerStr), &cust)
			myLogger.Debugf("customerStr : ", customerStr)
			myLogger.Debugf("sKeys[1] : ", sKeys[1])
			myLogger.Debugf("cust.PersonalDetails.PhoneNumber : ", cust.PersonalDetails.PhoneNumber)
			phoneA := StandardizePhoneNumber(sKeys[1], 5)
			phoneB := StandardizePhoneNumber(cust.PersonalDetails.PhoneNumber, 5)

			jsonResp = jsonResp + "{\"data\":" + customerStr + "," +
				"\"visibility\":" + ciav.GetVisibility(ciav.GetCallerRole(stub)) + "," +
				"\"allowedActions\":" + allowedActions +  "," +
				"\"phED\":" + strconv.Itoa(getEDScore(phoneA, phoneB, len(phoneA), len(phoneB))) +"}"

			// jsonResp = jsonResp + "{ data : "
			// jsonResp = jsonResp + customerStr + ","
			// jsonResp = jsonResp + "PH-ED:" + strconv.Itoa(getEDScore(phoneA, phoneB, len(phoneA), len(phoneB)))
			// jsonResp = jsonResp + "}"
		}
		jsonResp = jsonResp + "]"
	return jsonResp
}

func min(x int, y int, z int) int {
	if x < y && x < z {
		return x
	}
	if y < x && y < z {
		return y
	}
	return z
}

func getEDScore(str1 string, str2 string, m int, n int) int {
	// If first string is empty, the only option is to
	// insert all characters of second string into first
	if m == 0 {
		return n
	}
	// If second string is empty, the only option is to
	// remove all characters of first string
	if n == 0 {
		return m
	}
	// If last characters of two strings are same, nothing
	// much to do. Ignore last characters and get count for
	// remaining strings.
	if str1[m-1] == str2[n-1] {
		return getEDScore(str1, str2, m-1, n-1)
	}
	// If last characters are not same, consider all three(Insert,Remove, Replace)
	// operations on last character of first string, recursively
	// compute minimum cost for all three operations and take
	// minimum of three values.
	return 1 + min(getEDScore(str1, str2, m, n-1), getEDScore(str1, str2, m-1, n), getEDScore(str1, str2, m-1, n-1))
}

/*
  input   : customer
  output  : list of critical fields updated
  purpose : compare the existing record with the tobe updated records and find out the list of critical data
            which are getting updated
*/
func GetCriticalDataModified(stub shim.ChaincodeStubInterface, CustB ciav.Customer) (map[string]string, error) {
	UPDATED_ATTRIBUTES = make(map[string]string)
	CustomerAJSON := ciav.GetCustomerData(stub, CustB.PersonalDetails.CustomerId)
	var CustA ciav.Customer
	err := json.Unmarshal([]byte(CustomerAJSON), &CustA)
	if err != nil {
		fmt.Errorf("ERROR : Converting customer json to object", err)
		return nil, err
	}
	CustomerA := ciav.ConvertStructToMap(CustA)
	CustomerB := ciav.ConvertStructToMap(CustB)
	for k, v := range CRITICAL_DATA {
		myLogger.Debugf("Attribute  : ",k)
		myLogger.Debugf("CustomerA : ",CustomerA[k])
		myLogger.Debugf("CustomerB : ",CustomerB[k])

		val := strings.Split(v, "|")
		attrValB := StandardizeData(CustomerB[k], val[1])

		if CustomerA[k] != attrValB {
			myLogger.Debugf("Not Equal . . .")
			UPDATED_ATTRIBUTES[k] = attrVal
		}else{
			myLogger.Debugf("Equal . . .")
		}
	}
	return UPDATED_ATTRIBUTES, nil
}

/*
  input   : comparision string , updated attributes
  output  : updated comparision string
  purpose : updates comparision string based on the changes done on customer record
*/
func getUpdatedComparisonString(csToBeUpdated string, updated_attrs map[string]string) (string, error) {
	recToBeUpdated := strings.Split(csToBeUpdated, "~")
	myLogger.Debugf("Started updating CD : ", csToBeUpdated)
	for attr, val := range updated_attrs {
		myLogger.Debugf("Attr : ", attr)
		myLogger.Debugf("Val : ", val)
		tmp := strings.Split(CRITICAL_DATA[attr], "|")
		tAttr, err := strconv.Atoi(tmp[0])
		if err != nil {
			fmt.Errorf("ERROR : Incorrect index value for attribute", err)
			return "", err
		}
		recToBeUpdated[tAttr] = val
	}
	updatedCD := ""
	for _, val := range recToBeUpdated {
		if updatedCD != "" {
			updatedCD = updatedCD + "~"
		}
		updatedCD = updatedCD + val
	}
	return updatedCD, nil
}
