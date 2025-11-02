package cloudformation

type ResourceSerializer interface {
	Serialize() SerializeResult
	AddTag(key, value string)
}

type SerializeResult struct {
	ResourceID   string
	ResourceBody []byte
	Outputs      map[string]StackOutput
}
