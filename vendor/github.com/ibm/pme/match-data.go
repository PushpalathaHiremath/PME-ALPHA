package pme

import (
  "strconv"
  "github.com/op/go-logging"
  "github.com/ibm/ciav"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  "strings"
)

var myLogger = logging.MustGetLogger("pme-logs")
var dummyValue = "99999"
var BUCKET_CRITERIA1 = "Name+Phone"
var ANONYMOUS []string
var NICKNAMES map[string]string
var COMPARISON_ATTR []string
var CRITICAL_DATA map[string]string
var BUCKET_CRITERIAS []string
var COMPARISON_ATTR_DEF []string

/*
  input   :
  output  :
  purpose : Get critical data from configuration file
*/
func InitMatching ()(){
	CRITICAL_DATA = make(map[string]string)
	for idx, comparisonAttr := range COMPARISON_ATTR {
		if comparisonAttr != ""{
			attr, atype := getAttr(comparisonAttr)
			idxStr := strconv.Itoa(idx)
      COMPARISON_ATTR_DEF = append(COMPARISON_ATTR_DEF, attr)
			CRITICAL_DATA[attr] = idxStr + "|" + atype
		}
	}
}

func DuplicateExists(stub shim.ChaincodeStubInterface, Cust ciav.Customer)(string, bool){
  data := Standardize(Cust)
  isExists, _ := ComparisonStringExists(stub, data)
  return data, isExists
}

/*
  input   : Customer
  output  :
  purpose : create comparison string using the given customer's critical data,
            store it in comparison table and add the customer to the corresponding bucket
*/
func CollectMatchData(stub shim.ChaincodeStubInterface, Cust ciav.Customer, data string)(){
	myLogger.Debugf("Customer Id : ", Cust.PersonalDetails.CustomerId)
	myLogger.Debugf("CD : ", data)
	AddComparisonString(stub, Cust.PersonalDetails.CustomerId, data)
  UpdateBuckets(stub, Cust.PersonalDetails.CustomerId, data, "add")

  // for _ , bucket := range BUCKET_CRITERIAS {
  //   val := strings.Split(bucket, "+")
  //   bucket_key := ""
  //   for _,v := range val{
  //       attr_def := strings.Split(CRITICAL_DATA[v], "|")
  //       idx,_ := strconv.Atoi(attr_def[0])
  //       tval := StandardizeL2(getValAt(data, "~", idx), attr_def[1])
  //       if bucket_key != ""{
  //         bucket_key = bucket_key + "~"
  //       }
  //       bucket_key = bucket_key + tval
  //   }
  //     AddToBkt(stub, bucket, Cust.PersonalDetails.CustomerId, bucket_key)
  // }
}

func getValAt(str string, separator string, idx int)(string){
  val := strings.Split(str, separator)
  return val[idx]
}
