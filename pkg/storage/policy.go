package storage

// Policy represents an authorization policy
// It determines whether User/Role/Service x is allowed to perform Action y on Resource z
type Policy struct {
	ID int64 `json:"id"`
	audit Audit
	Version string `json:"version"`
	Statements []Statement `json:"statements"`
}

const (
	Version = "v1"
)

// NewPolicy returns a new Policy
func NewPolicy(statements []Statement) (*Policy) {
	if statements == nil {
		statements = []Statement{}
	}

	return &Policy{
		Version: Version,
		Statements: statements,
	}
}

// Statement is a rule statement within a Policy
type Statement struct {
	Effect Effect
	Actions []string
	Resources []string
}

// NewStatement creates a new statement based on the given list of resources & actions
func NewStatement(effect Effect, resources, actions []string) (*Statement) {
	if actions == nil {
		actions = []string{}
	}
	if resources == nil {
		resources = []string{}
	}

	return &Statement{
		Effect:    effect,
		Actions:   actions,
		Resources: resources,
	}
}

// Effect = Allow/Deny 
type Effect string

const (
	Deny Effect = "deny"
	Allow Effect = "allow"
)

