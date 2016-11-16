package pme

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/ibm/ciav"
  "strings"
	"strconv"
)

/*
  input   : comparison string
  output  : encoded hash string
  purpose : create encoding hash using SHA256
*/
func generateBucketHash(message string) string {
	key := []byte("secret")
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/*
  input   : Bucket criteria
  output  :
  purpose : creates the bucket table
*/
func CreateBucketHashTable(stub shim.ChaincodeStubInterface, bucketCriteria string) (error) {
	myLogger.Debug("Creating bucket : ", bucketCriteria)
	err := stub.CreateTable(bucketCriteria, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "dummy", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "BucketHash", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "CustomerIds", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating bucket")
	}
	return nil
}

/*
  input   : comparison string, CustomerId, comparison string
  output  :
  purpose : add the customer id against corresponding hash
*/
func AddToBkt(stub shim.ChaincodeStubInterface, bucketCriteria string, CustomerId string,data string) (error) {
	BucketHash := generateBucketHash(data)

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: BucketHash}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	hashrow, err := stub.GetRow(bucketCriteria, columns)
	if err != nil {
		return fmt.Errorf("Failed retriving details  for [%s]: [%s]", BucketHash, err)
	}
	isRowExists := (hashrow.Columns != nil)

	var ok bool
	if isRowExists {
		if !strings.Contains(hashrow.Columns[2].GetString_(), CustomerId){
			CustomerIds := hashrow.Columns[2].GetString_() + "|" + CustomerId
			ok, err = stub.ReplaceRow(bucketCriteria, shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: dummyValue}},
					&shim.Column{Value: &shim.Column_String_{String_: BucketHash}},
					&shim.Column{Value: &shim.Column_String_{String_: CustomerIds}},
				},
			})
			if !ok && err == nil {
				return errors.New("ERROR : error in adding data to bucket")
			}
			return nil
		}
	} else {
  		ok, err := stub.InsertRow(bucketCriteria, shim.Row{
  			Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: dummyValue}},
  				&shim.Column{Value: &shim.Column_String_{String_: BucketHash}},
  				&shim.Column{Value: &shim.Column_String_{String_: CustomerId}},
  			},
  		})
  		if !ok && err == nil {
  			return errors.New("Error in adding data to bucket")
  		}
	}
	return nil
}

/*
  input   : bucket criteria, comparison string
  output  : customers who fall into same bucket
  purpose : get the customers corresponding to given comparison string
*/
func SearchInBkt(stub shim.ChaincodeStubInterface, bucketCriteria string, data string) ([]string, error) {
	BucketHash := generateBucketHash(data)

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: BucketHash}}
	columns = append(columns, col1)
  columns = append(columns, col2)

	hashrow, err := stub.GetRow(bucketCriteria, columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retriving data for ", BucketHash, err)
	}
  var matchingRecords []string
	isRowExists := (hashrow.Columns != nil)
	if isRowExists {
    matchingRecords = strings.Split(hashrow.Columns[2].GetString_(), "|")
		return matchingRecords, nil
	}
  fmt.Println("No similar records found")
  return matchingRecords, nil
}

/*
  input   : bucket criteria, CustomerId, comparison string
  output  :
  purpose : Clear given customer from all the buckets
*/
func updateBucket(stub shim.ChaincodeStubInterface, bucketCriteria string, CustomerId string, data string) (error) {
	myLogger.Debugf("Updating Bucket ", bucketCriteria, "for CustomerId ", CustomerId)
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: dummyValue}}
	columns = append(columns, col1)

	rows, err := ciav.GetAllRows(stub, bucketCriteria, columns)
	if err != nil {
		return fmt.Errorf("Failed retriving data", err)
	}
	var idBuffer bytes.Buffer
	for i := range rows {
		myLogger.Debugf("Bucket rows : ")
		myLogger.Debugf("Hash - ", rows[i].Columns[1].GetString_())
		myLogger.Debugf("Customer Ids - ", rows[i].Columns[2].GetString_())
		if strings.Contains(rows[i].Columns[2].GetString_(), CustomerId) {
			ids := strings.Split(rows[i].Columns[2].GetString_(), "|")
			for _,val := range ids{
				if val != CustomerId {
					if idBuffer.String() != ""{
						idBuffer.WriteString("|")
					}
					idBuffer.WriteString(val)
				}
			}

			ok, err := stub.ReplaceRow(bucketCriteria, shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: dummyValue}},
					&shim.Column{Value: &shim.Column_String_{String_: rows[i].Columns[1].GetString_()}},
					&shim.Column{Value: &shim.Column_String_{String_: idBuffer.String()}},
				},
			})
			if !ok && err == nil {
				myLogger.Debugf("Error in updating bucket",err)
				return errors.New("Error in updating bucket")
			}
		}
	}
	myLogger.Debugf("Updating(Add) bucket for CustomerId : ", CustomerId)
	AddToBkt(stub, bucketCriteria, CustomerId, data)
	return nil
}

func UpdateBuckets(stub shim.ChaincodeStubInterface, CustomerId string, data string, op_type string) () {
	for _ , bucket := range BUCKET_CRITERIAS {
		val := strings.Split(bucket, "+")
		bucket_key := ""
		for _,v := range val{
				attr_def := strings.Split(CRITICAL_DATA[v], "|")
				idx,_ := strconv.Atoi(attr_def[0])
				myLogger.Debugf("data : ", data)
				myLogger.Debugf("idx : ", idx)
				tval := StandardizeL2(getValAt(data, "~", idx), attr_def[1])
				if bucket_key != ""{
					bucket_key = bucket_key + "~"
				}
				bucket_key = bucket_key + tval
		}
		myLogger.Debugf("Bucket data ")
		myLogger.Debugf("Type : ", op_type)
		myLogger.Debugf("bucket : ", bucket)
		myLogger.Debugf("bucket_key : ", bucket_key)
		if op_type == "add"{
			AddToBkt(stub, bucket, CustomerId, bucket_key)
		}else if op_type == "update" {
			updateBucket(stub, bucket, CustomerId, bucket_key)
		}
	}
}
