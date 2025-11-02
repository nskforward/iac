package cloudformation

type Tagger struct {
	tags []Tag
}

func (t *Tagger) AddTag(key, value string) {
	if t.tags == nil {
		t.tags = make([]Tag, 0, 16)
	}
	t.tags = append(t.tags, Tag{Key: key, Value: value})
}

func (t *Tagger) GetAllTags() []Tag {
	return t.tags
}

type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
