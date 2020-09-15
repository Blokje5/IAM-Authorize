package engine

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

// Input represents the expected query input
type Input struct {
	UserID int
	Action string
	Resource string
}

// Engine is the policy engine that executes authorization queries
type Engine struct {
	compiler *ast.Compiler
	store storage.Store
}

func (e *Engine) Query(ctx context.Context, input Input) (bool, error) {
	query := "data.main.authorized"

	options := []func(r *rego.Rego) {
		rego.Input(input),
		rego.Query(query),
		rego.Compiler(e.compiler),
		rego.Store(e.store),
	}

	resultSet, err := rego.New(options...).Eval(ctx)
	if err != nil {
		return false, fmt.Errorf("query Rego: %w", err)
	}

	var boolList []bool
	for _, result := range resultSet {
		for _, expression := range result.Expressions {
			switch expression.Value.(type) {
			case bool:
				boolList = append(boolList, expression.Value.(bool))
			default:
				// TODO deal with default case
			}
		}
	}

	return reduceBoolList(boolList), nil
}

func reduceBoolList(list []bool) bool {
	for _, v := range list {
		if !v {
			return false
		}
	}

	return true
}