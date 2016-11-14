/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package ciav

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Address struct {
	CustomerId  string
	AddressId   string
	AddressType string
	DoorNumber  string
	Street      string
	Locality    string
	City        string
	State       string
	Pincode     string
	PoaType     string
	PoaDoc      string
	Source      string
}

/*
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
																				Address
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
*/

/*
	Create address table
*/
func CreateAddressTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//	myLogger.Debug("Creating Address Table ...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err := stub.CreateTable("Address", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "customer_id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "address_id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "address_type", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "door_tumber", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "street", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "locality", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "city", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "state", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "pincode", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "poa_type", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "poa_doc", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "source", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Address table.")
	}

	// myLogger.Debug("Address table created Successfully... !!! ")
	return nil, nil
}

/*
	Add Address record
*/
func AddAddress(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// myLogger.Debug("Adding Address record ...")

	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12")
	}
	customerId := args[0]
	addressId := args[1]
	addressType := args[2]
	doorNumber := args[3]
	street := args[4]
	locality := args[5]
	city := args[6]
	state := args[7]
	pincode := args[8]
	poaType := args[9]
	poaDoc := args[10]
	source := args[11]

	ok, err := stub.InsertRow("Address", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			&shim.Column{Value: &shim.Column_String_{String_: addressId}},
			&shim.Column{Value: &shim.Column_String_{String_: addressType}},
			&shim.Column{Value: &shim.Column_String_{String_: doorNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: street}},
			&shim.Column{Value: &shim.Column_String_{String_: locality}},
			&shim.Column{Value: &shim.Column_String_{String_: city}},
			&shim.Column{Value: &shim.Column_String_{String_: state}},
			&shim.Column{Value: &shim.Column_String_{String_: pincode}},
			&shim.Column{Value: &shim.Column_String_{String_: poaType}},
			&shim.Column{Value: &shim.Column_String_{String_: poaDoc}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding address record.")
	}
	myLogger.Debug("Congratulations !!! Successfully added address",)
	return nil, err
}

func UpdateAddress(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Updating address record ...")

	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12")
	}

	customerId := args[0]
	addressId := args[1]
	addressType := args[2]
	doorNumber := args[3]
	street := args[4]
	locality := args[5]
	city := args[6]
	state := args[7]
	pincode := args[8]
	poaType := args[9]
	poaDoc := args[10]
	source := args[11]

	ok, err := stub.ReplaceRow("Address", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			&shim.Column{Value: &shim.Column_String_{String_: addressId}},
			&shim.Column{Value: &shim.Column_String_{String_: addressType}},
			&shim.Column{Value: &shim.Column_String_{String_: doorNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: street}},
			&shim.Column{Value: &shim.Column_String_{String_: locality}},
			&shim.Column{Value: &shim.Column_String_{String_: city}},
			&shim.Column{Value: &shim.Column_String_{String_: state}},
			&shim.Column{Value: &shim.Column_String_{String_: pincode}},
			&shim.Column{Value: &shim.Column_String_{String_: poaType}},
			&shim.Column{Value: &shim.Column_String_{String_: poaDoc}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in updated customer address record.")
	}
	// myLogger.Debug("Congratulations !!! Successfully updated ",)
	return nil, err
}

/*
 Get address record
*/
func GetAddress(stub shim.ChaincodeStubInterface, customerId string) (string, error) {
	var err error
	// myLogger.Debugf("Getting address record for customer : [%s]", string(customerId))

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: customerId}}
	columns = append(columns, col1)
	rows, err := GetAllRows(stub, "Address", columns)
	if err != nil {
		// myLogger.Debugf("Failed retriving Address details [%s]: [%s]", string(customerId), err)
		return "", fmt.Errorf("Failed retriving Address details [%s]: [%s]", string(customerId), err)
	}

	var jsonRespBuffer bytes.Buffer
	jsonRespBuffer.WriteString("[")
	for i := range rows {
		row := rows[i]
		// myLogger.Debugf("Identification rows [%s], is : [%s]", i, row)
		fmt.Println(row)
		if i != 0 {
			jsonRespBuffer.WriteString(",")
		}
		
		if len(row.Columns) == 0 {
		 	myLogger.Debugf("ERROR : Not found : ", customerId)
			return "", fmt.Errorf("Failed retriving Customer details [%s]: [%s]", string(customerId), err)
		 }
		jsonRespBuffer.WriteString("{\"customerId\":\"" + row.Columns[0].GetString_() + "\"" +
			",\"addressId\":\"" + row.Columns[1].GetString_() + "\"" +
			",\"addressType\":\"" + row.Columns[2].GetString_() + "\"" +
			",\"doorNumber\":\"" + row.Columns[3].GetString_() + "\"" +
			",\"street\":\"" + row.Columns[4].GetString_() + "\"" +
			",\"locality\":\"" + row.Columns[5].GetString_() + "\"" +
			",\"city\":\"" + row.Columns[6].GetString_() + "\"" +
			",\"state\":\"" + row.Columns[7].GetString_() + "\"" +
			",\"pincode\":\"" + row.Columns[8].GetString_() + "\"" +
			",\"poaType\":\"" + row.Columns[9].GetString_() + "\"" +
			",\"poaDoc\":\"" + row.Columns[10].GetString_() + "\"" +
			",\"source\":\"" + row.Columns[11].GetString_() + "\"}")
	}
	jsonRespBuffer.WriteString("]")

	return jsonRespBuffer.String(), nil
}
