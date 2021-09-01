package gridscale

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	boolInterfaceType   = "bool"
	intInterfaceType    = "int"
	floatInterfaceType  = "float"
	stringInterfaceType = "string"
)

var supportedPrimTypes = []string{boolInterfaceType, intInterfaceType, floatInterfaceType, stringInterfaceType}

//convSOStrings converts slice of interfaces to slice of strings
func convSOStrings(interfaceList []interface{}) []string {
	strList := make([]string, 0)
	for _, strIntf := range interfaceList {
		strList = append(strList, strIntf.(string))
	}
	return strList
}

//convStrToTypeInterface converts a string to a primitive type (in the form of interface{})
func convStrToTypeInterface(interfaceType, str string) (interface{}, error) {
	switch interfaceType {
	case boolInterfaceType:
		return strconv.ParseBool(str)
	case intInterfaceType:
		return strconv.Atoi(str)
	case floatInterfaceType:
		return strconv.ParseFloat(str, 64)
	case stringInterfaceType:
		return str, nil
	default:
		return nil, errors.New("type is invalid")
	}
}

//getInterfaceType gets interface type
func getInterfaceType(value interface{}) (string, error) {
	switch value.(type) {
	case bool:
		return boolInterfaceType, nil
	case int:
		return intInterfaceType, nil
	case float32:
		return floatInterfaceType, nil
	case float64:
		return floatInterfaceType, nil
	case string:
		return stringInterfaceType, nil
	default:
		return "", errors.New("type not found")
	}
}

//convInterfaceToString converts an interface of any primitive types to a  string value
func convInterfaceToString(interfaceType string, val interface{}) (string, error) {
	switch interfaceType {
	case boolInterfaceType:
		v, ok := val.(bool)
		if !ok {
			return "", fmt.Errorf("type assertion error:  value %v is not a bool value", val)
		}
		return strconv.FormatBool(v), nil
	case intInterfaceType:
		v, ok := val.(int)
		if !ok {
			return "", fmt.Errorf("type assertion error:  value %v is not an int value", val)
		}
		return strconv.FormatInt(int64(v), 10), nil
	case floatInterfaceType:
		v, ok := val.(float64)
		if !ok {
			return "", fmt.Errorf("type assertion error:  value %v is not a float64 value", val)
		}
		return strconv.FormatFloat(v, 'f', -1, 32), nil
	case stringInterfaceType:
		v, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("type assertion error:  value %v is not a string value", val)
		}
		return v, nil
	default:
		return "", errors.New("type is invalid")
	}
}
