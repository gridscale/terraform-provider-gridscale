package gsclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_isValidUUID(t *testing.T) {
	validationUUIDTestCases := make([]uuidTestCase, len(uuidCommonTestCases))
	copy(validationUUIDTestCases, uuidCommonTestCases)
	validationUUIDTestCases = append(validationUUIDTestCases,
		uuidTestCase{
			isFailed: true,
			testUUID: "abc-123",
		},
		uuidTestCase{
			isFailed: false,
			testUUID: "690de890-13c0-4e76-8a01-e10ba8786e53",
		},
	)
	for _, test := range validationUUIDTestCases {
		isValid := isValidUUID(test.testUUID)
		if test.isFailed {
			assert.False(t, isValid)
		} else {
			assert.True(t, isValid)
		}
	}
}
