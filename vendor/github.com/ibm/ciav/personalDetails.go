/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package ciav

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type PersonalDetails struct {
	CustomerId   string
	FirstName    string
	LastName     string
	Sex          string
	EmailId      string
	Dob          string
	PhoneNumber  string
	Occupation   string
	AnnualIncome string
	IncomeSource string
	Source       string
}

/*
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
																				Customer
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
*/

/*
	Create customer table
*/
func CreateCustomerTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// myLogger.Debug("Creating Customer table...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	// Create Customer table
	err := stub.CreateTable("Customer", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "customer_id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "first_name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "last_name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "sex", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "email_id", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "dob", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "phone_number", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "occupation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "annual_income", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "income_source", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "source", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed to create Customer table.")
	}
	// myLogger.Debug("Customer table created Successfully... !!! ")
	return nil, nil
}

/*
	Add customer record
*/
func AddCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debug("Adding Customer record ...")

	if len(args) != 11 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}
	customerId := args[0]
	firstName := args[1]
	lastName := args[2]
	sex := args[3]
	emailId := args[4]
	dob := args[5]
	phoneNumber := args[6]
	occupation := args[7]
	annualIncome := args[8]
	incomeSource := args[9]
	source := args[10]

	// Add customer entry to Customer table
	// myLogger.Debugf("Adding Customer [%s] ", firstName)

	ok, err := stub.InsertRow("Customer", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			&shim.Column{Value: &shim.Column_String_{String_: firstName}},
			&shim.Column{Value: &shim.Column_String_{String_: lastName}},
			&shim.Column{Value: &shim.Column_String_{String_: sex}},
			&shim.Column{Value: &shim.Column_String_{String_: emailId}},
			&shim.Column{Value: &shim.Column_String_{String_: dob}},
			&shim.Column{Value: &shim.Column_String_{String_: phoneNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: occupation}},
			&shim.Column{Value: &shim.Column_String_{String_: annualIncome}},
			&shim.Column{Value: &shim.Column_String_{String_: incomeSource}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding customer record.")
	}
	myLogger.Debug("Congratulations !!! Successfully added PersonalDetails", firstName)
	return nil, err
}

/*
	Update customer record
*/
func UpdateCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// myLogger.Debug("Updating Customer record ...")

	if len(args) != 11 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}

	customerId := args[0]
	firstName := args[1]
	lastName := args[2]
	sex := args[3]
	emailId := args[4]
	dob := args[5]
	phoneNumber := args[6]
	occupation := args[7]
	annualIncome := args[8]
	incomeSource := args[9]
	source := args[10]

	// Update customer entry to Customer table
	// myLogger.Debugf("Updating Customer [%s] ", firstName)

	ok, err := stub.ReplaceRow("Customer", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: customerId}},
			&shim.Column{Value: &shim.Column_String_{String_: firstName}},
			&shim.Column{Value: &shim.Column_String_{String_: lastName}},
			&shim.Column{Value: &shim.Column_String_{String_: sex}},
			&shim.Column{Value: &shim.Column_String_{String_: emailId}},
			&shim.Column{Value: &shim.Column_String_{String_: dob}},
			&shim.Column{Value: &shim.Column_String_{String_: phoneNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: occupation}},
			&shim.Column{Value: &shim.Column_String_{String_: annualIncome}},
			&shim.Column{Value: &shim.Column_String_{String_: incomeSource}},
			&shim.Column{Value: &shim.Column_String_{String_: source}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in updated customer record.")
	}
	// myLogger.Debug("Congratulations !!! Successfully updated [%s] ", firstName)
	return nil, err
}

/*
 Get customer record
*/
func GetCustomer(stub shim.ChaincodeStubInterface, customerId string) (string, error) {
	var err error
	myLogger.Debugf("Get personal details record for customer : [%s]", customerId)

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: customerId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Customer", columns)
	if err != nil {
		// myLogger.Debugf("Failed retriving Customer details [%s]: [%s]", string(customerId), err)
		return "", fmt.Errorf("Failed retriving Customer details [%s]: [%s]", string(customerId), err)
	}

       	 myLogger.Debugf("ERROR : Not found : ", row)
	 if row == nil {
	 	myLogger.Debugf("ERROR : Not found : ", customerId)
		return "", fmt.Errorf("Failed retriving Customer details [%s]: [%s]", string(customerId), err)
	 }
	jsonResp := "{\"customerId\":\"" + row.Columns[0].GetString_() + "\"" +
		",\"firstName\":\"" + row.Columns[1].GetString_() + "\"" +
		",\"lastName\":\"" + row.Columns[2].GetString_() + "\"" +
		",\"sex\":\"" + row.Columns[3].GetString_() + "\"" +
		",\"emailId\":\"" + row.Columns[4].GetString_() + "\"" +
		",\"dob\":\"" + row.Columns[5].GetString_() + "\"" +
		",\"phoneNumber\":\"" + row.Columns[6].GetString_() + "\"" +
		",\"occupation\":\"" + row.Columns[7].GetString_() + "\"" +
		",\"annualIncome\":\"" + row.Columns[8].GetString_() + "\"" +
		",\"incomeSource\":\"" + row.Columns[9].GetString_() + "\"" +
		",\"source\":\"" + row.Columns[10].GetString_() + "\"}"
	return jsonResp, nil
}
