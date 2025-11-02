package test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/nskforward/iac"
	"github.com/nskforward/iac/provider/aws"
	"github.com/nskforward/iac/provider/aws/cloudformation"
)

func TestCloudFormationStack(t *testing.T) {
	provider := aws.NewProvider()

	stack := cloudformation.NewStack(provider.Config, "test-stack", 5*time.Minute)
	stack.AddTag("app", "test")

	vpc1 := cloudformation.NewVPC("sanbox-1", "10.1.0.0/16", true, true)
	vpc1.AddTag("stack", "test-stack")

	vpc2 := cloudformation.NewVPC("sanbox-2", "10.2.0.0/16", true, true)
	vpc2.AddTag("stack", "test-stack")

	stack.AddResource(vpc1)
	stack.AddResource(vpc2)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := iac.Run(ctx, stack)
	if err != nil {
		slog.Error(err.Error())
	}
}
