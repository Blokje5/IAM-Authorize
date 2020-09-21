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
			p.log.Debug("Refresh in progress")
			err := p.refresh(ctx)
			if err != nil {
				p.log.Errorf("Failed to refresh policies. %w", err)
			}
			p.log.Debug("Refresh completed")
		}
	}
}

func (p *PolicyRefresher) refresh(ctx context.Context) error {
	userMap, err := p.getUserMap(ctx)
	if err != nil {
		return err
	}

	policyMap, err := p.getPolicyMap(ctx)
	if err != nil {
		return err
	}

	p.engine.SetPolicies(ctx, policyMap, userMap)
	return nil
}

func (p *PolicyRefresher) getUserMap(ctx context.Context) (map[int64][]int64, error) {
	users, err := p.store.ListUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users for refresh: %w", err)
	}

	userMap := make(map[int64][]int64)
	
	for _, user := range users {
		policies, err := p.store.GetPolicyIDsForUser(ctx, &user)
		if err != nil {
			return nil, fmt.Errorf("list policies IDs for refresh: %w", err)
		}

		userMap[user.ID] = policies
	}

	return userMap, nil
}

func (p *PolicyRefresher) getPolicyMap(ctx context.Context) (map[int64]storage.Policy, error) {
	policies, err := p.store.ListPolicies(ctx)
	if err != nil {
		return nil, fmt.Errorf("list policies for refresh: %w", err)
	}

	policyMap := make(map[int64]storage.Policy)
	
	for _, policy := range policies {
		policyMap[policy.ID] = policy
	}

	return policyMap, nil
}