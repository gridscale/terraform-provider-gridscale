package gridscale

import (
	"fmt"

	goVersion "github.com/hashicorp/go-version"
)

// Release struct
type Release struct {
	goVersion.Version
}

// Release struct
type ReleaseSpan struct {
	Start Release
	End   *Release
}

// Feature struct
type Feature struct {
	Description  string
	ReleaseSpans []ReleaseSpan
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

// NewReleaseSpan creates a new ReleaseSpan instance.
func NewReleaseSpan(startRepresentation string, endRepresentation *string) (*ReleaseSpan, error) {
	startRelease, err := NewRelease(startRepresentation)
	if err != nil {
		return nil, err
	}
	if endRepresentation == nil {
		return &ReleaseSpan{Start: *startRelease, End: nil}, nil
	}
	endRelease, err := NewRelease(*endRepresentation)
	if err != nil {
		return nil, err
	}
	if startRelease.CheckIfAhead(endRelease) {
		panic("Wrongly defined release span where start is ahead of end.")
	}
	return &ReleaseSpan{Start: *startRelease, End: endRelease}, nil
}

// CheckIfFeatureRequestIsApplicable checks by a Release receiver if feature request represented by a passed Feature instance is applicable.
func (r *Release) CheckIfFeatureRequestIsApplicable(fc *Feature) *ReleaseFeatureIncompatibilityError {
	var featureRequestIsApplicable Truth
	for i := 0; featureRequestIsApplicable == false && i < len(fc.ReleaseSpans); i++ {
		featureRequestIsApplicable = fc.ReleaseSpans[i].CheckIfReleaseIsPartOf(r)
	}
	if featureRequestIsApplicable == true {
		return nil
	}
	return &ReleaseFeatureIncompatibilityError{
		Detail: fmt.Sprintf("Feature '%s' is requested for but not available at release %s", fc.Description, r.String()),
	}
}

// CheckIfAhead checks by a Release receiver if it is ahead a passed Release instance.
func (r *Release) CheckIfAhead(releaseToCompare *Release) Truth {
	return Truth(r.GreaterThan(&releaseToCompare.Version))
}

// CheckIfBehind checks by a Release receiver if it is behind a passed Release instance.
func (r *Release) CheckIfBehind(releaseToCompare *Release) Truth {
	return Truth(r.GreaterThan(&releaseToCompare.Version))
}

// CheckIfEqual checks by a Release receiver if it equals a passed Release instance.
func (r *Release) CheckIfEqual(releaseToCompare *Release) Truth {
	return Truth(r.Equal(&releaseToCompare.Version))
}

// CheckIfReleaseIsPartOf checks by a ReleaseSpan receiver if a passed Release instance is known.
func (rs *ReleaseSpan) CheckIfReleaseIsPartOf(r *Release) Truth {
	equalKnot := (r.CheckIfEqual(&rs.Start) || r.CheckIfEqual(rs.End))
	inBetweenKnots := (r.CheckIfAhead(&rs.Start) && (&rs.End == nil || r.CheckIfBehind(rs.End)))
	return (equalKnot || inBetweenKnots)
}
