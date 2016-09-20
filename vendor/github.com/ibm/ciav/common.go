/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package ciav

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

var Superadmin map[string]string
var Manager map[string]string
var RelationalManager map[string]string
var Helpdesk map[string]string

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
func GetAllRows(stub *shim.ChaincodeStub, tableName string, columns []shim.Column) ([]shim.Row, error) {
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
				// myLogger.Debugf("Fetching row : [%s]", temprow.Columns[0].GetString_())
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
func GetCustomerID(stub *shim.ChaincodeStub, panId string) ([]string, error) {
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

func GetCallerRole(stub *shim.ChaincodeStub)(string){
	callerRole, _ := stub.ReadCertAttribute("role")
	return string(callerRole)
}

func GetVisibilityForCurrentUser(stub *shim.ChaincodeStub)(map[string]string){
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

func CanModifyIdentificationTable(stub *shim.ChaincodeStub)(bool){
	visibility := GetVisibilityForCurrentUser(stub)
	// "IdentityNumber": "W",
	// "PoiType":        "W",
	// "PoiDoc":         "W",
	if visibility["IdentityNumber"]=="W" && visibility["PoiType"]=="W" && visibility["PoiDoc"]=="W"{
		return true
	}
	return false
}

func CanModifyAddressTable(stub *shim.ChaincodeStub)(bool){
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

func CanModifyCustomerTable(stub *shim.ChaincodeStub)(bool){
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

func CanModifyKYCTable(stub *shim.ChaincodeStub)(bool){
	visibility := GetVisibilityForCurrentUser(stub)
	// "KycStatus":      "R",
	// "KycRiskLevel":   "N",
	// "LastUpdated":    "R",
	// callerRole := GetCallerRole(stub)

	if visibility["KycStatus"]=="W" && visibility["KycRiskLevel"]=="W" && visibility["LastUpdated"]=="W" {
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
