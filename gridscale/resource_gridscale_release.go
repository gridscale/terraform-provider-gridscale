package gridscale

import (
	"fmt"

	goVersion "github.com/hashicorp/go-version"
)

// Release struct
type Release struct {
	goVersion.Version
}

// Feature struct
type Feature struct {
	Description string
	Release
}

// NewRelease creates a new Release instance.
func NewRelease(representation string) (*Release, error) {
	version, err := goVersion.NewVersion(representation)
	if err != nil {
		return nil, err
	}
	release := Release{*version}
	return &release, nil
}

// CheckIfFeatureIsKnown checks by a Release receiver if a passed Feature instance is known.
func (r *Release) CheckIfFeatureIsKnown(f *Feature) *ReleaseFeatureIncompatibilityError {
	if r.LessThan(&f.Version) {
		return &ReleaseFeatureIncompatibilityError{
			Detail: fmt.Sprintf("Feature '%s' is part of release %s but requested for release %s.", f.Description, f.String(), r.String()),
		}
	}
	return nil
}
