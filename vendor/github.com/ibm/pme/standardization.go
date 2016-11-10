/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the IBM India Pvt Ltd, Version 1.0 (the "License");
*/

package pme

import (
	"fmt"
	"github.com/phonetics"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"github.com/ibm/ciav"
)

/*
  input   : customer
  output  : comparison string
  purpose : standardizes the critical data elements of given customer and creates the comparison string
*/
func Standardize(Cust ciav.Customer) string {
	data := ""
	custData := ciav.ConvertStructToMap(Cust)

	for i := range COMPARISON_ATTR_DEF {
		val := strings.Split(CRITICAL_DATA[COMPARISON_ATTR_DEF[i]], "|")
		myLogger.Debugf("Attributes : ", COMPARISON_ATTR_DEF[i], val[1])
		attrVal := StandardizeData(custData[COMPARISON_ATTR_DEF[i]], val[1])
		if data != "" {
			data = data + "~"
		}
		data = data + attrVal
	}

	// for k, v := range CRITICAL_DATA {
	// 	val := strings.Split(v, "|")
	// 	myLogger.Debugf("Attributes : ", k, val[1])
	// 	attrVal := StandardizeData(custData[k], val[1])
	// 	if data != "" {
	// 		data = data + "~"
	// 	}
	// 	data = data + attrVal
	// }
	myLogger.Debugf("Critical data : ", data)
	return data
}

/*
  input   : customer data, field value and type
  output  : standardized field value
  purpose : standardizes given field value based on the type mentioned
*/
func StandardizeData(attrVal string, attrType string) string {
	switch strings.ToLower(attrType) {
	case "name":
		return standardizeName(attrVal)
	case "phone:home":
		return StandardizePhoneNumber(attrVal, 6)
	case "phone:ofc":
		return StandardizePhoneNumber(attrVal, 6)
	case "phone:mob":
		return StandardizePhoneNumber(attrVal, 5)
	case "id":
		return StandardizeID(attrVal)
	default:
		myLogger.Errorf("ERROR: unrecognized attribute type in comparision element configuration. Please modify and try again.")
	}
	return ""
}

/*
  input   : name
  output  : standardized name
  purpose : trims out spaces & special charecters, removes anonymous values and replaces nicknames,
						get the sound marker
*/
func standardizeName(name string) string {
	myLogger.Debugf("Standardizing : ", name)
	// Trim spaces and special charecters
	str := trim(name)
	// Ignore if anonymous values found
	if str != "" {
		if isAnonymous(str) {
			str = ""
		}
	}
	// replace the nick names with std name
	if str != "" {
		str = getStdName(str)
	}
	// get sound markers
	// if str != "" {
	// 	str = getSoundMarker(str)
	// }
	return str
}

/*
  input   : PhoneNumber, trim lenghth (based on country and type)
  output  : standardized PhoneNumber
  purpose : trims out the spaces & special charecters, trims to the lenghth mentioned and sort the number
*/
func StandardizePhoneNumber(PhoneNumber string, trimLen int) string {
	// Trim spaces and special charecters and
	// TODO: truncate to last 7, 8 or 10 digits based on country
	PhoneNumber = trimPhoneNumber(PhoneNumber, trimLen)
	// PhoneNumber, _ = sortStr(PhoneNumber)
	return PhoneNumber
}

/*
  input   : id
  output  : standardized id
  purpose : trims out spaces & special charecters and sorts the id
*/
func StandardizeID(id string) string {
	// Trim spaces and special charecters
	id = trimID(id)
	id, _ = sortStr(id)
	return id
}

/*
  input   : string
  output  : boolean true/false
  purpose : checks if the given string is a anonymous value
*/
func isAnonymous(oName string) bool {
	// var err error
	// if len(ANONYMOUS) == 0 {
	//     ANONYMOUS,err = readFile("./anonymous.txt")
	//     if err != nil{
	//         myLogger.Debugf("Error reading anonymous dictionary.",err)
	//     }
	// }
	// for i := 0; i < len(ANONYMOUS); i++ {
	for _, anonymousName := range ANONYMOUS {
		if anonymousName == oName {
			return true
		}
	}
	return false
}

/*
  input   : string
  output  : Standardized name
  purpose : returns the standard name for the given nickname
*/
func getStdName(oName string) string {
	if len(NICKNAMES) == 0 {
		// NICKNAMES = make(map[string]string)
		// nickNames,err := readFile("./nicknames.txt")
		// if err != nil{
		//     myLogger.Debugf("Error reading nicknames dictionary.",err)
		// }
		// for i := 0; i < len(nickNames); i++ {
		//   if nickNames[i] != "" {
		//     prop := strings.Split(nickNames[i], "=")
		//     NICKNAMES[trimSpace(prop[0])]=trimSpace(prop[1])
		//   }
		// }
	}

	for stdName, nickName := range NICKNAMES {
		if strings.Contains(nickName, oName) {
			tmpNames := strings.Split(nickName, "|")
			// for i := 0; i < len(tmpNames); i++ {
			for _, tmpName := range tmpNames {
				if tmpName == oName {
					myLogger.Debugf("Nickname found for ", oName, " is ", stdName)
					return stdName
				}
			}
		}
	}
	myLogger.Debugf("No nickname for :", oName)
	return oName
}

/*
  input   : FirstName, LastName
  output  : formatted name string
  purpose : formats the name
*/
func formatName(FirstName string, LastName string) string {
	return LastName + ":" + FirstName
}

/*
  input   : name
  output  : sound marker string
  purpose : get the sound marker for the given name using Domnikov's phonetics library
*/
func getSoundMarker(name string) string {
	soundex := phonetics.EncodeSoundex(name)
	return soundex
}

/*
  input   : name
  output  : trimmed name
  purpose : trims the spaces, numbers and special charecters found in given name
*/
func trim(name string) string {
	// Method 01
	// r := strings.NewReplacer("<", "", ">", "","*", "","/", ""," ", "","\\", "","|", "","-", "","!", "")
	// name = r.Replace(name)

	// Method 02
	for _, runeValue := range name {
		if runeValue < 'A' || runeValue > 'z' {
			name = strings.Replace(name, string(runeValue), "", -1)
		}
	}
	return strings.ToUpper(name)
}

/*
  input   : array, string
  output  : boolean true/false
  purpose : verifies if the given string is present in the given array
*/
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/*
  input   : fileName
  output  : array containing file's content
  purpose : read the given file and add the content line by line to and array
*/
func readFile(fileName string) ([]string, error) {
	myLogger.Debugf("Open file: ", fileName)
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(contents), string('\n'))
	var values []string
	for _, line := range lines {
		values = append(values, trimSpace(line))
	}
	return values, nil
}

/*
  input   : string
  output  : trimmed string
  purpose : trims spaces in the given string
*/
func trimSpace(str string) string {
	return strings.TrimSpace(str)
}

/*
  input   : string
  output  : sorted string
  purpose : converts the string to int, sort them and convert back to string
*/
func sortStr(str string) (string, error) {
	var intArr []int
	intArr = make([]int, 0, len(str))

	for i := 0; i < len(str); i++ {
		s, err := strconv.Atoi(string(str[i]))
		if err != nil {
			return "", err
		}
		intArr = append(intArr, s)
	}

	sort.Ints(intArr)
	sorteStr := ""
	for i := 0; i < len(intArr); i++ {
		sorteStr = sorteStr + strconv.Itoa(intArr[i])
	}
	return sorteStr, nil
}

/*
  input   : PhoneNumber
  output  : trimmed PhoneNumber
  purpose : trims spaces, special characters from phone number
*/
func trimPhoneNumber(PhoneNumber string, trimLen int) string {
	for _, runeValue := range PhoneNumber {
		if runeValue < '0' || runeValue > '9' {
			PhoneNumber = strings.Replace(PhoneNumber, string(runeValue), "", -1)
		}
	}
	if trimLen > len(PhoneNumber) {
		fmt.Errorf("ERROR: Incorrect phone number")
		return ""
	}
	return string(PhoneNumber[len(PhoneNumber)-trimLen : len(PhoneNumber)])
}

/*
  input   : id
  output  : trimmed id
  purpose : trims spaces and special charecters from given id
*/
func trimID(id string) string {
	for _, runeValue := range id {
		if runeValue < '0' || runeValue > '9' || runeValue < 'A' || runeValue > 'z' {
			id = strings.Replace(id, string(runeValue), "", -1)
		}
	}
	return strings.ToUpper(id)
}

/*
  input   : attribute definition
  output  : attribute value, attribute type
  purpose : to get the value and type from attribute definitions given in configuration file
*/
func getAttr(attr string) (string, string) {
	tmpArr := strings.Split(attr, "|")
	if attr != "" && len(tmpArr) < 2 {
		myLogger.Errorf("ERROR: Attribute definitions are not proper. Please modify and try again.")
	}
	atype := ""
	tArr := strings.Split(tmpArr[1], ":")
	if len(tArr) == 1 {
		atype = tArr[0]
	} else {
		for i, val := range tArr {
			if i != 0 {
				atype = atype + ":"
			}
			atype = atype + val
		}
	}
	return tmpArr[0], atype
}

func StandardizeL2(attr_val string, attr_type string) (string) {
	switch strings.ToLower(attr_type) {
	case "name":
		return standardizeL2Name(attr_val)
	case "phone:home":
		return standardizeL2Ph(attr_val)
	case "phone:ofc":
		return standardizeL2Ph(attr_val)
	case "phone:mob":
		return standardizeL2Ph(attr_val)
	default:
		myLogger.Errorf("ERROR: unrecognized attribute type in comparision element configuration. Please modify and try again.")
	}
	return ""
}

func standardizeL2Name(name string)(string){
	if name != "" {
		return getSoundMarker(name)
	}
	return ""
}

func standardizeL2Ph(PhoneNumber string)(string){
	PhoneNumber, _ = sortStr(PhoneNumber)
	return PhoneNumber
}
