package gsclient

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type isContinue func() (bool, error)

//isValidUUID validates the uuid
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

//retryWithTimeout reruns a function within a period of time
func retryWithTimeout(targetFunc isContinue, timeout, delay time.Duration) error {
	timer := time.After(timeout)
	var err error
	var continueRetrying bool
	for {
		select {
		case <-timer:
			if err != nil {
				return err
			}
			return errors.New("timeout reached")
		default:
			time.Sleep(delay) //delay between retries
			continueRetrying, err = targetFunc()
			if !continueRetrying {
				return err
			}
		}
	}
}

//retryWithLimitedNumOfRetries reruns a function within a number of retries
func retryWithLimitedNumOfRetries(targetFunc isContinue, numOfRetries int, delay time.Duration) error {
	retryNo := 0
	var err error
	var continueRetrying bool
	for retryNo <= numOfRetries {
		time.Sleep(delay) //delay between retries
		continueRetrying, err = targetFunc()
		if !continueRetrying {
			return err
		}
		retryNo++
	}
	if err != nil {
		return err
	}
	return errors.New("timeout reached")
}
