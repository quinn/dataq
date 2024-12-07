package claims

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"go.quinn.io/dataq/cas"
)

type KV interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type ClaimsIndexer struct {
	cas cas.Storage
	kv  KV
}

// NewClaimsIndexer returns a ClaimsIndexer that can build and query an index of claims.
func NewClaimsIndexer(cas cas.Storage, kv KV) *ClaimsIndexer {
	return &ClaimsIndexer{
		cas: cas,
		kv:  kv,
	}
}

// RebuildIndex scans the entire Storage, finds claims, and indexes the latest claim per entity.
//
// Steps:
// - Iterate over all hashes in Storage.
// - For each hash, try to decode as a Claim.
// - If successful, keep track of the highest version claim per EntityUID.
// - After iteration, write the latest claim hash per entity to KV.
func (ci *ClaimsIndexer) RebuildIndex() error {
	hashes, err := ci.cas.Iterate()
	if err != nil {
		return fmt.Errorf("failed to iterate cas: %w", err)
	}

	// temporary map in memory: entityUID -> (version, claimHash)
	latestClaims := make(map[string]struct {
		version   int
		claimHash string
	})

	for hash := range hashes {
		c, err := ci.tryDecodeClaim(hash)
		if err != nil {
			// Not a claim or some other error; ignore and continue
			continue
		}

		prev, exists := latestClaims[c.EntityUID]
		if !exists || c.Version > prev.version {
			latestClaims[c.EntityUID] = struct {
				version   int
				claimHash string
			}{
				version:   c.Version,
				claimHash: hash,
			}
		}
	}

	// Persist the results to KV
	for entityUID, info := range latestClaims {
		if err := ci.kv.Set(entityUID, info.claimHash); err != nil {
			return fmt.Errorf("failed to set KV entry for entity %s: %w", entityUID, err)
		}
	}

	return nil
}

// tryDecodeClaim attempts to decode the content at the given hash as a Claim.
// Returns an error if the content is not a valid JSON-encoded Claim.
func (ci *ClaimsIndexer) tryDecodeClaim(hash string) (*Claim, error) {
	r, err := ci.cas.Retrieve(hash)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var c Claim
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	// Validate at least the EntityUID and Version fields to ensure it's a "real" claim.
	if c.EntityUID == "" || c.Version == 0 {
		return nil, errors.New("invalid claim structure")
	}

	return &c, nil
}

// GetLatestClaimHash returns the hash of the latest claim for the given entity, if present.
func (ci *ClaimsIndexer) GetLatestClaimHash(entityUID string) (string, error) {
	h, err := ci.kv.Get(entityUID)
	if err != nil {
		return "", fmt.Errorf("failed to get latest claim for entity %s: %w", entityUID, err)
	}
	if h == "" {
		return "", errors.New("no claim found for entity")
	}
	return h, nil
}

// RetrieveLatestClaim returns the latest claim object for the given entity.
// It looks up the entity in KV to find the latest claim hash, then retrieves the claim from Storage.
func (ci *ClaimsIndexer) RetrieveLatestClaim(entityUID string) (*Claim, error) {
	claimHash, err := ci.GetLatestClaimHash(entityUID)
	if err != nil {
		return nil, err
	}

	cs := &ClaimsService{cas: ci.cas}
	claim, err := cs.RetrieveClaim(claimHash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve latest claim for entity %s: %w", entityUID, err)
	}

	return claim, nil
}
