package repo

import (
	"context"
	"fmt"

	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/schema"
)

type Repo struct {
	index *index.Index
}

func NewRepo(idx *index.Index) *Repo {
	return &Repo{
		index: idx,
	}
}

func (r *Repo) PluginClaims(ctx context.Context) ([]schema.Claim, error) {
	sel := r.index.Q.
		GroupBy("permanode_hash").
		Where("schema_kind = ?", "PluginInstance").
		OrderBy("timestamp DESC")
	return r.index.Query(ctx, sel)
}

func (r *Repo) Plugins(ctx context.Context) ([]schema.PluginInstance, error) {
	claims, err := r.PluginClaims(ctx)
	if err != nil {
		return nil, err
	}

	plugins := make([]schema.PluginInstance, 0, len(claims))
	for _, claim := range claims {
		var plugin schema.PluginInstance
		// TODO: this makes another query to the index which is not necessary, we should
		// already have the content hash at this point
		if err := r.index.GetPermanode(ctx, claim.PermanodeHash, &plugin); err != nil {
			return nil, fmt.Errorf("failed to get plugin: %w", err)
		}
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}
