package cliutil

type FieldSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Desc string `json:"desc"`
}

type CommandSchema struct {
	Command       string                 `json:"command"`
	Entity        string                 `json:"entity"`
	Supports      map[string]interface{} `json:"supports,omitempty"`
	DefaultFields []string               `json:"default_fields,omitempty"`
	Fields        []FieldSchema          `json:"fields"`
}

func (s CommandSchema) FieldNames() []string {
	names := make([]string, 0, len(s.Fields))
	for _, field := range s.Fields {
		names = append(names, field.Name)
	}
	return names
}
