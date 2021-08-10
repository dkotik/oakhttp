package secrets

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func FromAWS(config *aws.Config, rotation time.Duration) (SecretsRepository, error) {

	go func(ctx context.Context) {

	}(context.Background())
	return nil, nil
}

type awsSecretsMap struct {
	// client   nil
	values map[string]interface{}
	mutex  *sync.RWMutex
	// rotation time.Duration
}

func (a *awsSecretsMap) Rotate(ctx context.Context) error {

	return nil
}
