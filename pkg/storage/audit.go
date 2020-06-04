package storage

import "time"

// Audit is a composition object that can be insterted into other objects to maintain audit information
type Audit struct {
	createdBy      string
	lastModifiedBy string
	createdAt      time.Time
	lastModifiedAt time.Time
}

// NewAudit returns a new Audit object
func NewAudit(creator string) Audit {
	ts := time.Now()

	return Audit{
		createdBy:      creator,
		lastModifiedBy: creator,
		createdAt:      ts,
		lastModifiedAt: ts,
	}
}

// Update updates the audit object and edits the lastModifiedAt state
func (a Audit) Update(modifier string) {
	ts := time.Now()

	a.lastModifiedBy = modifier
	a.lastModifiedAt = ts
}
