package gridscale

// ReleaseFeatureIncompatibilityError struct
type ReleaseFeatureIncompatibilityError struct {
	Detail string
}

// ReleaseFeatureIncompatibilityError receiver returning the error detail
func (rfie *ReleaseFeatureIncompatibilityError) Error() string {
	return rfie.Detail
}
