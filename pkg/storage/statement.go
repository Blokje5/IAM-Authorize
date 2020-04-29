package storage

import "strings"

type Statement struct {
	Effect string `json:"effect"`
	Action Action `json:"action"`
	Resource Resource `json:"resource"`
}

type Namespace string

type Action struct {
	Namespace
	Operation string
}

// UnmarshalJSON implements json unmarshalling for an Action reference
func (a *Action) UnmarshalJSON(b []byte) error  {
	split := strings.SplitN(string(b), ":", 2)
	a.Namespace = Namespace(split[0])
	a.Operation = split[1]

	return nil
}

type Resource struct {
	Namespace
	Object string
}

// UnmarshalJSON implements json unmarshalling for an Action reference
func (r *Resource) UnmarshalJSON(b []byte) error  {
	split := strings.SplitN(string(b), ":", 2)
	r.Namespace = Namespace(split[0])
	r.Object = split[1]

	return nil
}