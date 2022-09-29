package errorhandler

import "github.com/gridscale/gsclient-go/v3"

// SuppressHTTPErrorCodes suppresses the error, if the error
// is in the list errorCodes.
func SuppressHTTPErrorCodes(err error, errorCodes ...int) error {
	if requestError, ok := err.(gsclient.RequestError); ok {
		if containsInt(errorCodes, requestError.StatusCode) {
			err = nil
		}
	}
	return err
}

// containsInt check if an int array contains a specific int.
func containsInt(arr []int, target int) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}
