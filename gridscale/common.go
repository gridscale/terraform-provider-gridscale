package gridscale

import (
	"github.com/hashicorp/terraform/helper/resource"
	"time"
)

//type timeOutFunc func(payload interface{}) interface{}

//convSOStrings converts slice of interfaces to slice of strings
func convSOStrings(interfaceList []interface{}) []string {
	var labels []string
	for _, labelInterface := range interfaceList {
		labels = append(labels, labelInterface.(string))
	}
	return labels
}

//func timeOutTask(seconds int, task timeOutFunc) error

type retrier struct {
	timeout time.Duration
}

func (r retrier) basisRetryAfterFail(err error) error {
	return resource.Retry(r.timeout, func() *resource.RetryError {
		return resource.RetryableError(err)
	})
}
