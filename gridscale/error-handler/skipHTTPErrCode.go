package error_handler

import "github.com/gridscale/gsclient-go/v3"

//RemoveErrorContainsHTTPCodes returns nil, if the error of HTTP error
//has status code that is in the given list of http status codes
func RemoveErrorContainsHTTPCodes(err error, errorCodes ...int) error {
	if requestError, ok := err.(gsclient.RequestError); ok {
		if containsInt(errorCodes, requestError.StatusCode) {
			err = nil
		}
	}
	return err
}

//containsInt check if an int array contains a specific int.
func containsInt(arr []int, target int) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}
