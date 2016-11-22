/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/ibm/ciav"
	"github.com/ibm/pme"
	"github.com/op/go-logging"
	"strconv"
	"io/ioutil"
	"strings"
)

var myLogger = logging.MustGetLogger("customer_details")
var dummyValue = "99999"
var BKT_CRITERIA_DEFINITION = "FirstName+LastName|FirstName+PhoneNumber|LastName+PhoneNumber|FirstName+DOB1|LastName+DOB1|FirstName+DOB2|LastName+DOB2"

type ServicesChaincode struct {
}

func readFile(fileName string)([]string , error){
	myLogger.Debugf("Open file: ", fileName)
	contents, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(contents), string('\n'))
		var values []string
		for _, line := range lines {
			if line != "" {
					values = append(values, strings.TrimSpace(line))
			}
		}
		return values, nil
}

/*
   Deploy KYC data model
*/
func (t *ServicesChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	stub.PutState("bucket-criteria", []byte(BKT_CRITERIA_DEFINITION))
	stub.PutState("anonymous", []byte("XXX|ZZZ|GIRL1|GIRL2|GIRL3|XXXXXXXXXXXX|XXXXXX|XXXXXXX|XXXXXXXXXX|BABY|BB10|BBOY|BG10|BOY1|BOY2|BOY3|INVALID|FATHER|GIRL|MALE|DUPLICATE|TEST|TWIN|UNKNOWN|XXXXXXXXX|XXXXXXXX|XXXX|ZZZZ|MOTHER|SPOUSE|XXXXX|BABYBOY|DEPENDENT|BABYGIRL|TRAUMA|BGIRL|XXXXXXXXXXX|NOFIRSTNAME|NEWBORN|NOLASTNAME"))
	stub.PutState("nickNames", []byte("WILLIAM=BILL+ADELAIDE=ALEY|ELA|ELKE|LAIDEY|LAIDY+BENJAMIN=JAMIE|BIN|BENN|JAMEY+MADELINE=MADGE|MADIE+JOHNSON=JOHNSUN|JONSON|JONSUN+JENKINSON=JANKINSON|JAINKINSUN|JENKINSUN|JANKINSUN"))
	stub.PutState("comparison-attributes", []byte("FirstName|name+LastName|name+PhoneNumber|phone:home+DOB|date"))

	pme.BUCKET_CRITERIAS = strings.Split(BKT_CRITERIA_DEFINITION, "|")

	ciav.GetVisibility(ciav.GetCallerRole(stub))
	ciav.CreateIdentificationTable(stub, args)
	ciav.CreateCustomerTable(stub, args)
	ciav.CreateKycTable(stub, args)
	ciav.CreateAddressTable(stub, args)
	for _,bucket := range pme.BUCKET_CRITERIAS {
			pme.CreateBucketHashTable(stub, bucket)
	}
	pme.CreateComparisonStringTable(stub)
	return nil, nil
}

/*
  Add Customer record
*/
func (t *ServicesChaincode) addCIAV(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	myLogger.Debugf("Adding Customer record started...")
	myLogger.Debugf(args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments for addCIAV. Expecting 1")
	}

	myLogger.Debugf("Adding Customer record ...")
	var Cust ciav.Customer
	err := json.Unmarshal([]byte(string(args[0])), &Cust)
	if err != nil {
		fmt.Println("Error is :", err)
	}

	data, isCDExists := pme.DuplicateExists(stub, Cust)
	if isCDExists {
		return nil, errors.New("ERROR : Add failed. Duplicate customer - Data already exists.")
	}

	for i := range Cust.Identification {
		ciav.AddIdentification(stub, []string{Cust.Identification[i].CustomerId, Cust.Identification[i].IdentityNumber, Cust.Identification[i].PoiType, Cust.Identification[i].PoiDoc,
			Cust.Identification[i].Source})
	}

	ciav.AddCustomer(stub, []string{Cust.PersonalDetails.CustomerId, Cust.PersonalDetails.FirstName, Cust.PersonalDetails.LastName,
		Cust.PersonalDetails.Sex, Cust.PersonalDetails.EmailId, Cust.PersonalDetails.Dob, Cust.PersonalDetails.PhoneNumber, Cust.PersonalDetails.Occupation,
		Cust.PersonalDetails.AnnualIncome, Cust.PersonalDetails.IncomeSource, Cust.PersonalDetails.Source})

	ciav.AddKYC(stub, []string{Cust.Kyc.CustomerId, Cust.Kyc.KycStatus, Cust.Kyc.LastUpdated, Cust.Kyc.Source, Cust.Kyc.KycRiskLevel})

	for i := range Cust.Address {
		ciav.AddAddress(stub, []string{Cust.Address[i].CustomerId, Cust.Address[i].AddressId, Cust.Address[i].AddressType,
			Cust.Address[i].DoorNumber, Cust.Address[i].Street, Cust.Address[i].Locality, Cust.Address[i].City, Cust.Address[i].State,
			Cust.Address[i].Pincode, Cust.Address[i].PoaType, Cust.Address[i].PoaDoc, Cust.Address[i].Source})
	}
	// match data using PME
	pme.CollectMatchData(stub, Cust, data)

	myLogger.Debugf("Add . . .")
	return []byte("Add Successful...."), nil
}

/*
 Update customer record
*/
func (t *ServicesChaincode) updateCIAV(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments for updateCIAV. Expecting 1")
	}

	var Cust ciav.Customer
	err := json.Unmarshal([]byte(string(args[0])), &Cust)
	if err != nil {
		fmt.Println("Error is :", err)
	}

	isCDModified := false
	updatedCD, _ := pme.GetCriticalDataModified(stub, Cust)
	if len(updatedCD) != 0{
		isCDModified = true
	}
	if ciav.CanModifyIdentificationTable(stub){
		for i := range Cust.Identification {
			ciav.UpdateIdentification(stub, []string{Cust.Identification[i].CustomerId, Cust.Identification[i].IdentityNumber, Cust.Identification[i].PoiType, Cust.Identification[i].PoiDoc,
				Cust.Identification[i].Source})
		}
	}
	if ciav.CanModifyCustomerTable(stub){
		ciav.UpdateCustomer(stub, []string{Cust.PersonalDetails.CustomerId, Cust.PersonalDetails.FirstName, Cust.PersonalDetails.LastName,
			Cust.PersonalDetails.Sex, Cust.PersonalDetails.EmailId, Cust.PersonalDetails.Dob, Cust.PersonalDetails.PhoneNumber, Cust.PersonalDetails.Occupation,
			Cust.PersonalDetails.AnnualIncome, Cust.PersonalDetails.IncomeSource, Cust.PersonalDetails.Source})
	}
	if ciav.CanModifyKYCTable(stub){
		ciav.UpdateKYC(stub, []string{Cust.Kyc.CustomerId, Cust.Kyc.KycStatus, Cust.Kyc.LastUpdated, Cust.Kyc.Source, Cust.Kyc.KycRiskLevel})
	}

	if ciav.CanModifyAddressTable(stub){
		for i := range Cust.Address {
			ciav.UpdateAddress(stub, []string{Cust.Address[i].CustomerId, Cust.Address[i].AddressId, Cust.Address[i].AddressType,
				Cust.Address[i].DoorNumber, Cust.Address[i].Street, Cust.Address[i].Locality, Cust.Address[i].City, Cust.Address[i].State,
				Cust.Address[i].Pincode, Cust.Address[i].PoaType, Cust.Address[i].PoaDoc, Cust.Address[i].Source})
		}
	}

	if isCDModified {
		updatedCS,_ := pme.UpdateComparisonString(stub, Cust.PersonalDetails.CustomerId, updatedCD)
		// pme.UpdateBuckets(stub, ciav.BUCKET_CRITERIA1, Cust.PersonalDetails.CustomerId, updatedCS)
		pme.UpdateBuckets(stub, Cust.PersonalDetails.CustomerId, updatedCS, "update")
	}
	return nil, nil
}

/*
   Invoke : addCIAV and updateCIAV
*/
func (t *ServicesChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(pme.COMPARISON_ATTR) == 0{
		err := getConfiguration(stub);
		if err != nil{
			return nil, fmt.Errorf("Failed retriving configurations [%s]", err)
		}
	}

	if function == "addCIAV" {
		// add customer
		return t.addCIAV(stub, args)
	} else {
		// update customer
		return t.updateCIAV(stub, args)
	}
	return nil, errors.New("Received unknown function invocation")
}

func getConfiguration(stub shim.ChaincodeStubInterface)(error){
			bucketCriteriaStr, bc_err := stub.GetState("bucket-criteria")
		  pme.BUCKET_CRITERIAS = strings.Split(string(bucketCriteriaStr), "|")

			anonymousStr, a_err := stub.GetState("anonymous")
		  pme.ANONYMOUS = strings.Split(string(anonymousStr), "|")

			// comparisinAttrsStr, ca_err := stub.GetState("comparison-attributes")
			// pme.COMPARISON_ATTR = strings.Split(string(comparisinAttrsStr), "+")

			pme.COMPARISON_ATTR = append(pme.COMPARISON_ATTR, "FirstName|name")
			pme.COMPARISON_ATTR = append(pme.COMPARISON_ATTR, "LastName|name")
			pme.COMPARISON_ATTR = append(pme.COMPARISON_ATTR, "PhoneNumber|phone:home")
			pme.COMPARISON_ATTR = append(pme.COMPARISON_ATTR, "DOB1|date1")
			pme.COMPARISON_ATTR = append(pme.COMPARISON_ATTR, "DOB2|date2")

			nickNamesStr, n_err := stub.GetState("nickNames")
		  nickNames := strings.Split(string(nickNamesStr), "+")

			pme.NICKNAMES = make(map[string]string)

			for i := 0; i < len(nickNames); i++ {
				if nickNames[i] != "" {
					prop := strings.Split(nickNames[i], "=")
					pme.NICKNAMES[strings.TrimSpace(prop[0])]=strings.TrimSpace(prop[1])
				}
			}

			if bc_err != nil || a_err != nil || n_err != nil {
				return errors.New("ERROR : Fetching configurations.")
			}

			pme.InitMatching()
			return nil
}

/*
	Get Customer record by customer id or PAN number
*/
func (t *ServicesChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "getCIAV" {
		return t.getCIAV(stub, args)
	} else if function == "searchRecords" {
		return t.searchRecords(stub, args)
	}else if function == "getKYCStats" {
		return t.getKYCStats(stub)
	}
	return nil, errors.New("Received unknown function invocation")
}
func (t *ServicesChaincode) searchRecords(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments for searchRecords. Expecting 2")
	}
	jsonResp := pme.SearchRecords(stub, args[0], args[1])

	// customerRecord,_ := ciav.GetCustomerRecord(stub, args[0])
	// var Cust ciav.Customer
	// err := json.Unmarshal([]byte(customerRecord), &Cust)
	//
	// jsonResp := ciav.SearchMatches(stub, Cust)
	bytes, err := json.Marshal(jsonResp)
	if err != nil {
		return nil, errors.New("Error converting kyc record")
	}
	return bytes, nil
}

func (t *ServicesChaincode) getCIAV(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments for getCIAV. Expecting 2")
	}
	var jsonResp string
	var customerIds []string
	var err error

	if args[0] == "PAN" {
		customerIds, err = ciav.GetCustomerID(stub, args[1])
		// jsonResp = "["
		for i := range customerIds {
			if i != 0 {
				jsonResp = jsonResp + ","
			}
			jsonResp = ciav.GetCustomerData(stub, customerIds[i])
		}
		// jsonResp = jsonResp + "]"
	} else if args[0] == "CUST_ID" {
		jsonResp = ciav.GetCustomerData(stub, args[1])
	} else {
		return nil, errors.New("Invalid arguments. Please query by CUST_ID or PAN")
	}
	bytes, err := json.Marshal(jsonResp)
	if err != nil {
		return nil, errors.New("Error converting kyc record")
	}
	return bytes, nil
}

func (t *ServicesChaincode) getKYCStats(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var err error

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	columns = append(columns, col1)
	rows, err := ciav.GetAllRows(stub, "KYC", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retriving KYC details [%s]", err)
	}

	var kycBuffer bytes.Buffer
	// var compliantBuffer bytes.Buffer
	// var noncompliantBuffer bytes.Buffer
	var compliantCustomersCount int
	var noncompliantCustomersCount int
	var totalCustomers int

	for i := range rows {
		row := rows[i]
		totalCustomers++
		if row.Columns[2].GetString_() == "compliant" {
			compliantCustomersCount++
			// if compliantBuffer.String() != "" {
			// 	compliantBuffer.WriteString(",")
			// }
			// compliantBuffer.WriteString("{\"customerId\":\"" + row.Columns[1].GetString_() + "\"" +
			// 	",\"kycStatus\":\"" + row.Columns[2].GetString_() + "\"" +
			// 	",\"lastUpdated\":\"" + row.Columns[3].GetString_() + "\"" +
			// 	",\"source\":\"" + row.Columns[4].GetString_() + "\"}")
		} else if row.Columns[2].GetString_() == "non-compliant" {
			noncompliantCustomersCount++
			// 	if noncompliantBuffer.String() != "" {
			// 		noncompliantBuffer.WriteString(",")
			// 	}
			// 	noncompliantBuffer.WriteString("{\"customerId\":\"" + row.Columns[1].GetString_() + "\"" +
			// 		",\"kycStatus\":\"" + row.Columns[2].GetString_() + "\"" +
			// 		",\"lastUpdated\":\"" + row.Columns[3].GetString_() + "\"" +
			// 		",\"source\":\"" + row.Columns[4].GetString_() + "\"}")
		}
	}
	kycBuffer.WriteString("{" +
		"\"compliant\" : \"" + strconv.Itoa(compliantCustomersCount) + "\"," +
		"\"noncompliant\" : \"" + strconv.Itoa(noncompliantCustomersCount) + "\"," +
		"\"total\" : \"" + strconv.Itoa(totalCustomers) + "\"" +
		"}")

	bytes, err := json.Marshal(kycBuffer.String())
	if err != nil {
		return nil, errors.New("Error converting kyc stats")
	}
	return bytes, nil
}

func main() {
	err := shim.Start(new(ServicesChaincode))
	if err != nil {
		fmt.Printf("Error starting ServicesChaincode: %s", err)
	}
}
