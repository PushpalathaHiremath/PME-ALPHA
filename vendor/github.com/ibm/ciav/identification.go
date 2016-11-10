/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package ciav

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

var dummyValue = "99999"

type Identification struct {
	CustomerId     string
	IdentityNumber string
	PoiType        string
	PoiDoc         string
	Source         string
}

/*
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
																				identification
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
*/

/*
 Create Identification table
*/
func CreateIdentificationTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// myLogger.Debug("Init Identification Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	// Create Identification table
	err := stub.CreateTable("Identification", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "customer_id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "identity_number", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "poi_type", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "poi_doc", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "source", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	// Create identification relation table
	err = stub.CreateTable("IDRelation", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "identity_number", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "customer_id", Type: shim.ColumnDefinition_STRING, Key: false},
	})

	if err != nil {
		return nil, errors.New("Failed creating Identification table.")
	}
	// myLogger.Debug("Identification table initialization done Successfully... !!! ")
	return nil, nil
}

/*
	add Identification record
*/
func AddIdentification(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Add Identification record ...")
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	customerId := args[0]
	identityNumber := args[1]
	poiType := args[2]
	poiDoc := args[3]
	source := args[4]

	// myLogger.Debugf("Adding identity doc : [%s] ", poiType)

	ok, err := stub.InsertRow("Identification", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			&shim.Column{Value: &shim.Column_String_{String_: identityNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: poiType}},
			&shim.Column{Value: &shim.Column_String_{String_: poiDoc}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	// update identification relation
	_, err = updateIDRelation(stub, identityNumber, customerId, "add")
	if !ok && err == nil {
		return nil, errors.New("Error in adding Identification record.")
	}

	myLogger.Debug("Congratulations !!! Successfully added,", customerId, identityNumber, poiType, poiDoc, source)
	return nil, err
}

/*
 Update Identification record
*/
func UpdateIdentification(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// myLogger.Debug("Update Identification record ...")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	// customerId := args[0]
	// identityNumber := args[1]
	// poiType := args[2]
	// poiDoc := args[3]
	// source := args[4]

	var customerId string
	var identityNumber string
	var poiType string
	var poiDoc string
	var source string

	myLogger.Debugf("Updating identity : [%s] ", poiType)

	isOk, _ := stub.VerifyAttribute("role", []byte("Helpdesk"))

	callerRole, _ := stub.ReadCertAttribute("role")
	myLogger.Debugf("caller role : [%s] ", string(callerRole))

	if isOk {
		identificationStr, _ := GetIdentification(stub, args[0])
		myLogger.Debugf("Before identificationStr : [%s] ", identificationStr)
		// identificationStr1 := identificationStr[2:len(identificationStr)-2]
		// identificationStr1 := strings.Replace(identificationStr, "[[", "[", -1)
		// identificationStr1 = strings.Replace(identificationStr1, "]]", "]", -1)
		// identificationStr1 = strings.Replace(identificationStr1, "CustomerId", "custId", -1)
		//
		// myLogger.Debugf("After identificationStr1 : [%s] ", identificationStr1)
		//
		identificationStr1 := strings.Replace(string(identificationStr), "[[", "[", -1)
		identificationStr1 = strings.Replace(identificationStr1, "]]", "]", -1)
		// identificationStr1 = strings.Replace(identificationStr1, "]", " ", -1)
		// identificationStr1 = strings.Replace(identificationStr1, "[", " ", -1)
		// identificationStr1 =  "[" +identificationStr1 + "]"

		myLogger.Debugf("After Trial identificationStr1 : [%s] ", identificationStr1)

		var identification []Identification
		err := json.Unmarshal([]byte(string(identificationStr1)), &identification)
		if err == nil {
			// return nil, errors.New("Error in getting Identification record.")
		}
		// customerId = identification[0].CustomerId
		// identityNumber = identification[0].IdentityNumber
		// 	poiType = args[2]
		// 	poiDoc = args[3]
		// 	source = args[4]

		customerId = args[0]
		identityNumber = args[1]
		poiType = args[2]
		poiDoc = args[3]
		source = args[4]

	} else {
		customerId = args[0]
		identityNumber = args[1]
		poiType = args[2]
		poiDoc = args[3]
		source = args[4]
	}
	myLogger.Debugf("customerId : [%s] ", customerId)
	myLogger.Debugf("identityNumber : [%s] ", identityNumber)
	myLogger.Debugf("poiType : [%s] ", poiType)
	myLogger.Debugf("poiDoc : [%s] ", poiDoc)
	myLogger.Debugf("source : [%s] ", source)

	ok, err := stub.ReplaceRow("Identification", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			&shim.Column{Value: &shim.Column_String_{String_: identityNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: poiType}},
			&shim.Column{Value: &shim.Column_String_{String_: poiDoc}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	_, err = updateIDRelation(stub, identityNumber, customerId, "update")
	ok, err = stub.ReplaceRow("IDRelation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: identityNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in updating Identification record.")
	}
	// myLogger.Debug("Congratulations !!! Successfully updated [%s]", res)
	return nil, err
}

/*
	Get Identification record
*/
func GetIdentification(stub shim.ChaincodeStubInterface, customerId string) (string, error) {
	var err error
	// myLogger.Debugf("Get identification record for customer : [%s]", string(customerId))

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: customerId}}
	columns = append(columns, col1)
	rows, err := GetAllRows(stub, "Identification", columns)
	if err != nil {
		// myLogger.Debugf("Failed retriving Identification details [%s]: [%s]", string(customerId), err)
		return "", fmt.Errorf("Failed retriving Identification details [%s]: [%s]", string(customerId), err)
	}

	var jsonRespBuffer bytes.Buffer
	jsonRespBuffer.WriteString("[")
	for i := range rows {
		row := rows[i]
		myLogger.Debugf("Identification rows [%s], is : [%s]", i, row)
		fmt.Println(row)
		if i != 0 {
			jsonRespBuffer.WriteString(",")
		}
		jsonRespBuffer.WriteString("{\"CustomerId\":\"" + row.Columns[0].GetString_() + "\"" +
			",\"IdentityNumber\":\"" + row.Columns[1].GetString_() + "\"" +
			",\"PoiType\":\"" + row.Columns[2].GetString_() + "\"" +
			",\"PoiDoc\":\"" + row.Columns[3].GetString_() + "\"" +
			",\"Source\":\"" + row.Columns[4].GetString_() + "\"}")
	}
	jsonRespBuffer.WriteString("]")
	myLogger.Debugf("Id Arr : [%s] ", jsonRespBuffer.String())
	return jsonRespBuffer.String(), nil
}

/*
	Update ID relation table
*/
func updateIDRelation(stub shim.ChaincodeStubInterface, identityNumber string, customerId string, functionType string) (string, error) {
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: identityNumber}}
	columns = append(columns, col1)

	idrow, err := stub.GetRow("IDRelation", columns)
	if err != nil {
		// myLogger.Debugf("Failed retriving Identification relation details for ID [%s]: [%s]", string(identityNumber), err)
		return "", fmt.Errorf("Failed retriving Identification relation details  for ID [%s]: [%s]", string(identityNumber), err)
	}

	var isRowExists bool
	isRowExists = (idrow.Columns != nil)

	var ok bool
	if isRowExists {
		if functionType == "update" {
			if strings.Contains(idrow.Columns[1].GetString_(), customerId) {
				// myLogger.Debugf("Identification relation exists. Do nothing.")
				return "", nil
			}
		}
		customerId = idrow.Columns[1].GetString_() + "|" + customerId
		ok, err = stub.ReplaceRow("IDRelation", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: identityNumber}},
				&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			},
		})
	} else {
		ok, err = stub.InsertRow("IDRelation", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: identityNumber}},
				&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			},
		})
	}
	if !ok && err == nil {
		return "", errors.New("Error in updating Identification relation record.")
	}
	return "", nil
}
