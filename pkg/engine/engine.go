package engine

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
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

func (e *Engine) SetPolicies(ctx context.Context, policyMap map[string]interface{}, roleMap map[string]interface{}) *Engine {
	e.store = inmem.NewFromObject(map[string]interface{}{
		"policies": policyMap,
		"roles":    roleMap,
	})

	return e
}

func Load(paths []string) (*Engine, error) {
	result, err := loader.All(paths)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	if len(result.Modules) == 0 {
		return nil, fmt.Errorf("no policies found in %v", paths)
	}

	compiler, err := result.Compiler()
	if err != nil {
		return nil, fmt.Errorf("get compiler: %w", err)
	}

	store, err := result.Store()
	if err != nil {
		return nil, fmt.Errorf("get store: %w", err)
	}

	e := &Engine{
		store: store,
		compiler: compiler,
	}

	return e, nil
}