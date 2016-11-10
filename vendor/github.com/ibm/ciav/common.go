/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package ciav

import (
	"bytes"
	"fmt"
	"github.com/op/go-logging"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

var myLogger = logging.MustGetLogger("ciav-details")
var Superadmin map[string]string
var Manager map[string]string
var RelationalManager map[string]string
var Helpdesk map[string]string

type Customer struct {
	Identification  []Identification
	PersonalDetails PersonalDetails
	Kyc             Kyc
	Address         []Address
}

func GetVisibility(callerRole string)(string) {
	Superadmin = map[string]string{
		"CustomerId":     "W",
		"IdentityNumber": "W",
		"PoiType":        "W",
		"PoiDoc":         "W",
		"Source":         "W",
		"FirstName":      "W",
		"LastName":       "W",
		"Sex":            "W",
		"EmailId":        "W",
		"Dob":            "W",
		"PhoneNumber":    "W",
		"Occupation":     "W",
		"AnnualIncome":   "W",
		"IncomeSource":   "W",
		"KycStatus":      "W",
		"KycRiskLevel":   "W",
		"LastUpdated":    "W",
		"AddressId":      "W",
		"AddressType":    "W",
		"DoorNumber":     "W",
		"Street":         "W",
		"Locality":       "W",
		"City":           "W",
		"State":          "W",
		"Pincode":        "W",
		"PoaType":        "W",
		"PoaDoc":         "W"}

	Manager = map[string]string{
		"CustomerId":     "W",
		"IdentityNumber": "W",
		"PoiType":        "W",
		"PoiDoc":         "W",
		"Source":         "W",
		"FirstName":      "W",
		"LastName":       "W",
		"Sex":            "W",
		"EmailId":        "W",
		"Dob":            "W",
		"PhoneNumber":    "W",
		"Occupation":     "W",
		"AnnualIncome":   "W",
		"IncomeSource":   "W",
		"KycStatus":      "W",
		"KycRiskLevel":   "W",
		"LastUpdated":    "W",
		"AddressId":      "W",
		"AddressType":    "W",
		"DoorNumber":     "W",
		"Street":         "W",
		"Locality":       "W",
		"City":           "W",
		"State":          "W",
		"Pincode":        "W",
		"PoaType":        "W",
		"PoaDoc":         "W"}

	RelationalManager = map[string]string{
		"CustomerId":     "W",
		"IdentityNumber": "W",
		"PoiType":        "W",
		"PoiDoc":         "W",
		"Source":         "W",
		"FirstName":      "W",
		"LastName":       "W",
		"Sex":            "W",
		"EmailId":        "W",
		"Dob":            "W",
		"PhoneNumber":    "W",
		"Occupation":     "W",
		"AnnualIncome":   "W",
		"IncomeSource":   "W",
		"KycStatus":      "W",
		"KycRiskLevel":   "N",
		"LastUpdated":    "W",
		"AddressId":      "W",
		"AddressType":    "W",
		"DoorNumber":     "W",
		"Street":         "W",
		"Locality":       "W",
		"City":           "W",
		"State":          "W",
		"Pincode":        "W",
		"PoaType":        "W",
		"PoaDoc":         "W"}
	Helpdesk = map[string]string{
		"CustomerId":     "R",
		"IdentityNumber": "W",
		"PoiType":        "W",
		"PoiDoc":         "W",
		"Source":         "W",
		"FirstName":      "R",
		"LastName":       "R",
		"Sex":            "R",
		"EmailId":        "R",
		"Dob":            "R",
		"PhoneNumber":    "R",
		"Occupation":     "R",
		"AnnualIncome":   "R",
		"IncomeSource":   "R",
		"KycStatus":      "W",
		"KycRiskLevel":   "N",
		"LastUpdated":    "W",
		"AddressId":      "W",
		"AddressType":    "W",
		"DoorNumber":     "W",
		"Street":         "W",
		"Locality":       "W",
		"City":           "W",
		"State":          "W",
		"Pincode":        "W",
		"PoaType":        "W",
		"PoaDoc":         "W"}

		visibility := Helpdesk
		if callerRole == "Superadmin" {
			visibility = Superadmin
		} else if callerRole == "RelationalManager" {
			visibility = RelationalManager
		} else if callerRole == "Manager" {
			visibility = Manager
		}

		var visibilityBuffer bytes.Buffer
		visibilityBuffer.WriteString("{")
		i := 0
		for key, value := range visibility {
			if i > 0 {
				visibilityBuffer.WriteString(",")
			}
			visibilityBuffer.WriteString("\"" + key + "\":\"" + value + "\"")
			i++
		}
		visibilityBuffer.WriteString("}")
		return visibilityBuffer.String()
}

/*
	Get all rows corresponding to the partial keys given
*/
func GetAllRows(stub shim.ChaincodeStubInterface, tableName string, columns []shim.Column) ([]shim.Row, error) {
	rowChannel, err := stub.GetRows(tableName, columns)
	if err != nil {
		// myLogger.Debugf("Failed retriving address details for : [%s]", err)
		return nil, fmt.Errorf("Failed retriving address details : [%s]", err)
	}
	var rows []shim.Row
	for {
		select {
		case temprow, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				rows = append(rows, temprow)
			}
		}
		if rowChannel == nil {
			break
		}
	}
	return rows, nil
}

/*
 Get the customer id by PAN number
*/
func GetCustomerID(stub shim.ChaincodeStubInterface, panId string) ([]string, error) {
	var err error

	// myLogger.Debugf("Get customer id for PAN : [%s]", panId)

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: panId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("IDRelation", columns)
	if err != nil {
		// myLogger.Debugf("Failed retriving Identification details for PAN [%s]: [%s]", string(panId), err)
		return nil, fmt.Errorf("Failed retriving Identification details  for PAN [%s]: [%s]", string(panId), err)
	}
	custIds := row.Columns[1].GetString_()
	custIdArray := strings.Split(custIds, "|")
	return custIdArray, nil
}

func GetCallerRole(stub shim.ChaincodeStubInterface)(string){
	callerRole, _ := stub.ReadCertAttribute("role")
	return string(callerRole)
}

func GetVisibilityForCurrentUser(stub shim.ChaincodeStubInterface)(map[string]string){
	callerRole := GetCallerRole(stub)

	visibility := Helpdesk
	if callerRole == "Superadmin" {
		visibility = Superadmin
	} else if callerRole == "RelationalManager" {
		visibility = RelationalManager
	} else if callerRole == "Manager" {
		visibility = Manager
	}
	return visibility
}

func CanModifyIdentificationTable(stub shim.ChaincodeStubInterface)(bool){
	visibility := GetVisibilityForCurrentUser(stub)
	// "IdentityNumber": "W",
	// "PoiType":        "W",
	// "PoiDoc":         "W",
	if visibility["IdentityNumber"]=="W" && visibility["PoiType"]=="W" && visibility["PoiDoc"]=="W"{
		return true
	}
	return false
}

func CanModifyAddressTable(stub shim.ChaincodeStubInterface)(bool){
	visibility := GetVisibilityForCurrentUser(stub)
	// "AddressId":      "W",
	// "AddressType":    "W",
	// "DoorNumber":     "W",
	// "Street":         "W",
	// "Locality":       "W",
	// "City":           "W",
	// "State":          "W",
	// "Pincode":        "W",
	// "PoaType":        "W",
	// "PoaDoc":         "W"}
	if visibility["AddressId"]=="W" && visibility["AddressType"]=="W" && visibility["DoorNumber"]=="W"  && visibility["Street"]=="W" && visibility["Locality"]=="W"  && visibility["City"]=="W"  && visibility["State"]=="W" && visibility["Pincode"]=="W" && visibility["PoaType"]=="W" && visibility["PoaDoc"]=="W"{
		return true
	}
	return false
}

func CanModifyCustomerTable(stub shim.ChaincodeStubInterface)(bool){
	visibility := GetVisibilityForCurrentUser(stub)
	// "FirstName":      "W",
	// "LastName":       "W",
	// "Sex":            "W",
	// "EmailId":        "W",
	// "Dob":            "W",
	// "PhoneNumber":    "W",
	// "Occupation":     "W",
	// "AnnualIncome":   "W",
	// "IncomeSource":   "W",
	if visibility["FirstName"]=="W" && visibility["LastName"]=="W" && visibility["Sex"]=="W" && visibility["EmailId"]=="W" && visibility["Dob"]=="W" && visibility["PhoneNumber"]=="W" && visibility["Occupation"]=="W" && visibility["AnnualIncome"]=="W" && visibility["IncomeSource"]=="W"{
		return true
	}
	return false
}

func CanModifyKYCTable(stub shim.ChaincodeStubInterface)(bool){
	visibility := GetVisibilityForCurrentUser(stub)
	// "KycStatus":      "R",
	// "KycRiskLevel":   "N",
	// "LastUpdated":    "R",
	// callerRole := GetCallerRole(stub)

	if visibility["KycStatus"]=="W" && visibility["LastUpdated"]=="W" {
		// if riskLevel == "3"{
		// 	return true
		// }else if riskLevel == "2"{
		// 	if callerRole == "Superadmin" || callerRole == "Manager" || callerRole == "RelationalManager"{
		// 		return true
		// 	}else{
		// 		return false
		// 	}
		// }else if riskLevel == "1"{
		// 	if callerRole == "Superadmin" || callerRole == "Manager"{
		// 		return true
		// 	}else{
		// 		return false
		// 	}
		// }
		return true
	}
	return false
}


func GetCustomerRecord(stub shim.ChaincodeStubInterface, customerId string)(string, string){
	var err error
	var identificationStr string
	var customerStr string
	var kycStr string
	var addressStr string
	var riskLevel string

	identificationStr, err = GetIdentification(stub, customerId)
	customerStr, err = GetCustomer(stub, customerId)
	kycStr, riskLevel, err = GetKYC(stub, customerId)
	addressStr, err = GetAddress(stub, customerId)

	if err != nil{
		myLogger.Debugf("Failed retriving customer details for : [%s], [%s]", customerId, err)
	}

	jsonResp := "{\"Identification\":" + identificationStr +
		",\"PersonalDetails\":" + customerStr +
		",\"KYC\":" + kycStr +
		",\"address\":" + addressStr + "}"

	return jsonResp, riskLevel
}

func GetCustomerData(stub shim.ChaincodeStubInterface, customerId string)(string){
	jsonResp, riskLevel := GetCustomerRecord(stub, customerId)
	allowedActions := GetPermissionMetadata(stub, riskLevel)

	responseStr := "{\"data\":" + jsonResp + "," +
		"\"visibility\":" + GetVisibility(GetCallerRole(stub)) + "," +
		"\"allowedActions\":" + allowedActions + "}"
	return responseStr
}

func GetPermissionMetadata(stub shim.ChaincodeStubInterface, riskLevel string)(string){
	callerRole := GetCallerRole(stub)
	allowedActions := "{\"updateKYCDocs\":\""

	if riskLevel == "3"{
		allowedActions = allowedActions + "true"
	}else if riskLevel == "2"{
		if callerRole == "Superadmin" || callerRole == "Manager" || callerRole == "RelationalManager"{
			allowedActions = allowedActions + "true"
		}else{
			allowedActions = allowedActions + "false"
		}
	}else if riskLevel == "1"{
		if callerRole == "Superadmin" || callerRole == "Manager"{
			allowedActions = allowedActions + "true"
		}else{
			allowedActions = allowedActions + "false"
		}
	}
	allowedActions = allowedActions + "\"}"
	return allowedActions
}

func ConvertStructToMap(Cust Customer)(map[string]string){
	allAttributes := map[string]string{
		"FirstName":     Cust.PersonalDetails.FirstName,
		"LastName":     Cust.PersonalDetails.LastName,
		"Sex":     Cust.PersonalDetails.Sex,
		"EmailId":     Cust.PersonalDetails.EmailId,
		"Dob":     Cust.PersonalDetails.Dob,
		"PhoneNumber":     Cust.PersonalDetails.PhoneNumber,
		"Occupation":     Cust.PersonalDetails.Occupation,
		"AnnualIncome":     Cust.PersonalDetails.AnnualIncome,
		"IncomeSource":     Cust.PersonalDetails.IncomeSource,
		"Source":     Cust.PersonalDetails.Source}
	return allAttributes
}
