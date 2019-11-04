package aws

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// Sts wrapper
type Sts struct {
	Client stsiface.STSAPI
}

// Sts wrapper constructor
func NewSts(stsapi stsiface.STSAPI) *Sts {
	return &Sts{Client: stsapi}
}

// Fetach API caller identity
func (s *Sts) GetCallerId() (*sts.GetCallerIdentityOutput, error) {
	result, err := s.Client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return result, nil
}
