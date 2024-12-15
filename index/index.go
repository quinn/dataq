package index

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/claims"
	"go.quinn.io/dataq/schema"
	"gorm.io/gorm"
)

type ClaimsIndexer struct {
	cl  claims.Claimer
	cas cas.Storage
	db  *gorm.DB
}

// NewClaimsIndexer returns a ClaimsIndexer that can build and query an index of claims.
func NewClaimsIndexer(cas cas.Storage, db *gorm.DB) *ClaimsIndexer {
	return &ClaimsIndexer{
		cas: cas,
		cl:  claims.NewClaimsService(cas),
		db:  db,
	}
}

func (ci *ClaimsIndexer) StoreClaim(ctx context.Context, c *claims.Claim) (string, error) {
	hash, err := ci.cl.StoreClaim(ctx, c)
	if err != nil {
		return "", fmt.Errorf("failed to store claim: %w", err)
	}

	if err := ci.IndexClaim(ctx, hash); err != nil {
		return "", fmt.Errorf("failed to index claim: %w", err)
	}

	return hash, nil
}

func (ci *ClaimsIndexer) IndexClaim(ctx context.Context, hash string) error {
	claim, payload, err := ci.cl.RetrieveClaim(ctx, hash)
	if err != nil {
		return fmt.Errorf("failed to retrieve claim: %w", err)
	}

	var newModel claims.Claimable
	var existingModel claims.Claimable
	switch payload.Kind {
	case "task":
		newModel = &schema.Task{}
		existingModel = &schema.Task{}
	default:
		return errors.New("unsupported claim kind: " + payload.Kind)
	}

	if err := json.Unmarshal(payload.Payload, existingModel); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := ci.db.WithContext(ctx).First(existingModel, "uid = ?", claim.UID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			existingModel = nil
		} else {
			return fmt.Errorf("failed to query model: %w", err)
		}
	}

	if existingModel == nil {
		if err := ci.db.WithContext(ctx).Create(newModel).Error; err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}
	} else {
		if existingModel.GetClaimTimestamp() >= claim.Timestamp {
			return nil
		}

		// TODO: gorm apparently ignores zero values when updating so this could cause problems down the line
		if err := ci.db.WithContext(ctx).Updates(newModel).Error; err != nil {
			return fmt.Errorf("failed to update model: %w", err)
		}
	}

	return nil
}

// RebuildIndex scans the entire Storage, finds claims, and indexes the latest claim per entity.
//
// Steps:
// - Iterate over all hashes in Storage.
// - For each hash, try to decode as a Claim.
// - If successful, keep track of the highest version claim per EntityUID.
// - After iteration, write the latest claim hash per entity to KV.
func (ci *ClaimsIndexer) RebuildIndex(ctx context.Context) error {
	hashes, err := ci.cas.Iterate(ctx)
	if err != nil {
		return fmt.Errorf("failed to iterate cas: %w", err)
	}

	// TODO: increase performance by refactoring IndexClaim to not check for existing model and timestamp
	// temporary map in memory: entityUID -> (version, claimHash)
	// latestClaims := make(map[string]struct {
	// 	Timestamp int64
	// 	claimHash string
	// })

	for hash := range hashes {
		if err := ci.IndexClaim(ctx, hash); err != nil {
			return fmt.Errorf("failed to index claim: %w", err)
		}

		// prev, exists := latestClaims[c.UID]
		// if !exists || c.Timestamp > prev.Timestamp {
		// 	latestClaims[c.UID] = struct {
		// 		Timestamp int64
		// 		claimHash string
		// 	}{
		// 		Timestamp: c.Timestamp,
		// 		claimHash: hash,
		// 	}
		// }
	}

	return nil
}

// RetrieveClaim attempts to decode the content at the given hash as a Claim.
// Returns an error if the content is not a valid JSON-encoded Claim.
func (ci *ClaimsIndexer) RetrieveClaim(ctx context.Context, hash string) (*claims.Claim, *claims.ClaimPayload, error) {
	return ci.cl.RetrieveClaim(ctx, hash)
}
