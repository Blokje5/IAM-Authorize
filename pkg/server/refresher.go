package server

import (
	"context"
	"fmt"
	"time"

	"github.com/blokje5/iam-server/pkg/engine"
	"github.com/blokje5/iam-server/pkg/log"
	"github.com/blokje5/iam-server/pkg/storage"
)

// PolicyRefresher ensures the engine is kept up to date with the relevant data
type PolicyRefresher struct {
	engine *engine.Engine
	store *storage.Storage

	log *log.Logger
}

func NewPolicyRefresher(engine *engine.Engine, store *storage.Storage) *PolicyRefresher {
	return &PolicyRefresher{
		engine: engine,
		store: store,

		log: log.GetLogger(),
	}
}

// Run runs the refresher as a background goroutine
func (p *PolicyRefresher) Run(ctx context.Context) {
	go p.run(ctx)
}

func (p *PolicyRefresher) run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := p.refresh(ctx)
			if err != nil {
				p.log.Errorf("Failed to refresh policies. %w", err)
			}
		}
	}
}

func (p *PolicyRefresher) refresh(ctx context.Context) error {
	refreshMap, err := p.getRefreshMap(ctx)
	if err != nil {
		return err
	}

	p.engine.SetPolicies(ctx, refreshMap)
}

func (p *PolicyRefresher) getRefreshMap(ctx context.Context) (map[int64][]storage.Policy, error) {
	users, err := p.store.ListUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users for refresh: %w", err)
	}

	policyMap := make(map[int64][]storage.Policy)
	
	for _, user := range users {
		policies, err := p.store.GetPoliciesForUser(ctx, &user)
		if err != nil {
			return nil, fmt.Errorf("list policies for refresh: %w", err)
		}

		policyMap[user.ID] = policies
	}

	return policyMap, nil
}