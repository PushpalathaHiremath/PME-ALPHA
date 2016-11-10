package pme
//
// import (
// 	"fmt"
// 	"errors"
// 	"github.com/hyperledger/fabric/core/chaincode/shim"
//   "strings"
// )
//
// /*
//   input   :
//   output  :
//   purpose :
// */
// func CreateNXPTable(stub shim.ChaincodeStubInterface) (error) {
// 	myLogger.Debug("Creating wts table NAME-PHONE")
// 	err := stub.CreateTable("NAME-PHONE", []*shim.ColumnDefinition{
// 		&shim.ColumnDefinition{Name: "Comparison_Type", Type: shim.ColumnDefinition_STRING, Key: true},
// 		&shim.ColumnDefinition{Name: "Index", Type: shim.ColumnDefinition_INT32, Key: true},
// 		&shim.ColumnDefinition{Name: "Ph-ED0", Type: shim.ColumnDefinition_INT32, Key: false},
// 		&shim.ColumnDefinition{Name: "Ph-ED1", Type: shim.ColumnDefinition_INT32, Key: false},
// 		&shim.ColumnDefinition{Name: "Ph-ED2", Type: shim.ColumnDefinition_INT32, Key: false},
// 		&shim.ColumnDefinition{Name: "Ph-ED3", Type: shim.ColumnDefinition_INT32, Key: false},
// 	})
// 	if err != nil {
// 		return errors.New("Failed creating bucket")
// 	}
// 	return nil
// }
//
// /*
//   input   :
//   output  :
//   purpose :
// */
// func CreateNameXactTable(stub shim.ChaincodeStubInterface) (error) {
// 	myLogger.Debug("Creating wts table NAME")
// 	err := stub.CreateTable("NAME", []*shim.ColumnDefinition{
// 		&shim.ColumnDefinition{Name: "Comparison_Type", Type: shim.ColumnDefinition_STRING, Key: true},
// 		&shim.ColumnDefinition{Name: "Index", Type: shim.ColumnDefinition_INT32, Key: true},
// 		&shim.ColumnDefinition{Name: "weight", Type: shim.ColumnDefinition_INT32, Key: false},
// 	})
// 	if err != nil {
// 		return errors.New("Failed creating bucket")
// 	}
// 	return nil
// }
//
// /*
//   input   :
//   output  :
//   purpose :
// */
// func loadNameXact(stub shim.ChaincodeStubInterface) (error) {
// 	var data map[int]int
// 	data[0] = 500
// 	data[1] = 400
// 	data[2] = 300
// 	data[3] = 200
//
// 	for k, v := range data{
// 		ok, err := stub.InsertRow("NAME", shim.Row{
// 			Columns: []*shim.Column{
// 				&shim.Column{Value: &shim.Column_String_{String_: "NAME-XACT"}},
// 				&shim.Column{Value: &shim.Column_String_{String_: k}},
// 				&shim.Column{Value: &shim.Column_String_{String_: v}},
// 			},
// 		})
// 		if !ok && err == nil {
// 			return errors.New("Error in adding data to bucket")
// 		}
// 	}
// 	return nil
// }
//
// /*
//   input   :
//   output  :
//   purpose :
// */
// func loadNamePhone(stub shim.ChaincodeStubInterface) (error) {
// 	var data map[int]string
// 	data[0] = "123|123|123|123"
// 	data[1] = "123|123|123|123"
// 	data[2] = "123|123|123|123"
// 	data[3] = "123|123|123|123"
//
// 	for k, v := range data{
// 		values := strings.Split(v, "|")
// 		ok, err := stub.InsertRow("NAME", shim.Row{
// 			Columns: []*shim.Column{
// 				&shim.Column{Value: &shim.Column_String_{String_: "NXP-2DM"}},
// 				&shim.Column{Value: &shim.Column_String_{String_: k}},
// 				&shim.Column{Value: &shim.Column_String_{String_: values[0]}},
// 				&shim.Column{Value: &shim.Column_String_{String_: values[1]}},
// 				&shim.Column{Value: &shim.Column_String_{String_: values[2]}},
// 				&shim.Column{Value: &shim.Column_String_{String_: values[3]}},
// 			},
// 		})
// 		if !ok && err == nil {
// 			return errors.New("Error in adding data to bucket")
// 		}
// 	}
// 	return nil
// }
//
//
// /*
//   input   :
//   output  :
//   purpose :
// */
// func GetNameXactScore(stub shim.ChaincodeStubInterface, index int) (int, error) {
// 	BucketHash := generateBucketHash(data)
//
// 	var columns []shim.Column
// 	col1 := shim.Column{Value: &shim.Column_String_{String_: "NAME-XACT"}}
// 	col2 := shim.Column{Value: &shim.Column_String_{String_: index}}
// 	columns = append(columns, col1)
//   columns = append(columns, col2)
//
// 	row, err := stub.GetRow("NAME", columns)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed retriving data for ", index, err)
// 	}
//
//   return matchingRecords, nil
// }
