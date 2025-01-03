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

const unsupportedK8SRelease = "1.30"

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
func (r *Release) CheckIfFeatureIsKnown(f *Feature) error {
	if r.LessThan(&f.Version) {
		return &ReleaseFeatureIncompatibilityError{
			Detail: fmt.Sprintf("Feature '%s' is part of release %s but requested for release %s.", f.Description, f.String(), r.String()),
		}
	}
	return nil
}

// CheckIfK8SReleaseIsSupported checks if the Kubernetes release is supported by this gridscale terraform provider.
func (r *Release) CheckIfK8SReleaseIsSupported() error {
	if r.GreaterThanOrEqual(goVersion.Must(goVersion.NewVersion(unsupportedK8SRelease))) {
		return &ReleaseFeatureIncompatibilityError{
			Detail: fmt.Sprintf("this gridscale terraform provider version v1 supports only Kubernetes release < %s, for Kubernetes release %s please use gridscale terraform provider version v2", unsupportedK8SRelease, r.String()),
		}
	}
	return nil
}
