package storage

type Policy struct {
	ID string `json:"id"`
	Version string `json:"version"`
	Statements []Statement `json:"statements"`
}