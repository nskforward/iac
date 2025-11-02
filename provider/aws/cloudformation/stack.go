package cloudformation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/nskforward/iac"
)

type Stack struct {
	Tagger
	cfg        aws.Config
	timeout    time.Duration
	name       string
	resources  []ResourceSerializer
	parameters []types.Parameter
}

type StackTemplate struct {
	AWSTemplateFormatVersion string                     `json:"AWSTemplateFormatVersion"`
	Description              string                     `json:"Description"`
	Resources                map[string]json.RawMessage `json:"Resources"`
	Outputs                  map[string]json.RawMessage `json:"Outputs"`
}

type StackOutput struct {
	Description string          `json:"Description"`
	Value       json.RawMessage `json:"Value"`
}

func NewStack(cfg aws.Config, name string, timeout time.Duration) *Stack {
	stack := &Stack{
		cfg:     cfg,
		timeout: timeout,
		name:    name,
	}
	stack.AddTag("Stack", name)
	return stack
}

func (stack *Stack) AddResource(resource ResourceSerializer) {
	if stack.resources == nil {
		stack.resources = make([]ResourceSerializer, 0, 16)
	}
	for _, tag := range stack.tags {
		resource.AddTag(tag.Key, tag.Value)
	}
	stack.resources = append(stack.resources, resource)
}

func (stack *Stack) AddParam(key, value string) {
	if stack.parameters == nil {
		stack.parameters = make([]types.Parameter, 0, 16)
	}
	stack.parameters = append(stack.parameters, types.Parameter{
		ParameterKey:   aws.String(key),
		ParameterValue: aws.String(value),
	})
}

func (stack *Stack) Run(ctx context.Context, args ...any) iac.Output {
	output := make(iac.Output, 1)

	template := StackTemplate{
		AWSTemplateFormatVersion: "2010-09-09",
		Description:              "CloudFormation automatic stack",
		Resources:                make(map[string]json.RawMessage),
		Outputs:                  make(map[string]json.RawMessage),
	}

	go func() {
		defer close(output)

		for _, resource := range stack.resources {
			res := resource.Serialize()
			template.Resources[res.ResourceID] = res.ResourceBody
			for key, val := range res.Outputs {
				valBytes, err := json.Marshal(val)
				if err != nil {
					output.PushError([]byte(err.Error()))
					return
				}
				template.Outputs[key] = valBytes
			}
		}

		templateBody, err := json.Marshal(template)
		if err != nil {
			output.PushError([]byte(err.Error()))
			return
		}

		input := &cf.CreateStackInput{
			TimeoutInMinutes: aws.Int32(int32(stack.timeout.Minutes())),
			StackName:        aws.String(stack.name),
			TemplateBody:     aws.String(string(templateBody)),
			Capabilities: []types.Capability{
				types.CapabilityCapabilityIam,
			},
			Parameters: stack.parameters,
		}

		client := cf.NewFromConfig(stack.cfg)
		t1 := time.Now()
		result, err := client.CreateStack(context.TODO(), input)
		t2 := time.Since(t1)
		if err != nil {
			output.PushError([]byte(err.Error()))
			return
		}

		output.PushOK(fmt.Appendf(nil, "stack id %s created for %v", *result.StackId, t2))

	}()

	return output
}
