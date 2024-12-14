package claims

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"go.quinn.io/dataq/cas"
)

// Claim represents a versioned assertion about an entity.
// It includes the entity's unique identifier, a reference to the previous claim,
// and a reference to the associated data.
type Claim struct {
	UID         string `json:"dataq_entity_uid"`        // unique identifier of the entity. shared across multiple claims, and used to reduce all claims into a single entity.
	ParentUID   string `json:"dataq_entity_parent_uid"` // parent entity UID. This is linked through plugin requests and responses
	Timestamp   int64  `json:"dataq_claim_timestamp"`   // Unix timestamp of creation
	Kind        string `json:"dataq_schema_kind"`       // kind of the entity. Has meaning to DataQ
	Hash        string `json:"hash,omitempty"`          // reference to blob in CAS
	PluginID    string `json:"plugin_id,omitempty"`     // ID of plugin
	PluginKind  string `json:"kind,omitempty"`          // meaningful for the plugin
	ContentType string `json:"content_type,omitempty"`  // content type. provided by plugin
}

// ClaimsService provides methods to store and retrieve claims in a CAS.
type ClaimsService struct {
	cas cas.Storage
}

// NewClaimsService returns a new ClaimsService backed by the given Storage.
func NewClaimsService(s cas.Storage) *ClaimsService {
	return &ClaimsService{
		cas: s,
	}
}

// StoreClaim stores a new claim in the CAS.
// The claim is serialized to JSON and stored.
// Returns the hash of the stored claim.
func (cs *ClaimsService) StoreClaim(ctx context.Context, c *Claim) (string, error) {
	c.Timestamp = time.Now().Unix()
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	hash, err := cs.cas.Store(ctx, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return hash, nil
}

// RetrieveClaim fetches and deserializes a claim from the CAS by its hash.
func (cs *ClaimsService) RetrieveClaim(ctx context.Context, hash string) (*Claim, error) {
	r, err := cs.cas.Retrieve(ctx, hash)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(r)
	var c Claim
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// CreateInitialClaim creates the first claim for an entity, referencing no previous claim.
// dataHash is the CAS hash of the associated data. Returns the claim hash.
func (cs *ClaimsService) CreateInitialClaim(ctx context.Context, UID string, hash string) (string, error) {
	claim := &Claim{
		UID:  UID,
		Hash: hash,
	}
	return cs.StoreClaim(ctx, claim)
}

// CreateNewClaim creates a new claim for an entity that already has at least one version.
// It references the previous claim and increments the version.
// prevClaimHash is the hash of the previous claim, dataHash is the CAS hash of the updated data.
func (cs *ClaimsService) CreateNewClaim(ctx context.Context, prevClaimHash, dataHash string) (string, error) {
	prevClaim, err := cs.RetrieveClaim(ctx, prevClaimHash)
	if err != nil {
		return "", err
	}
	// TODO: claims are not a linked list, no need for parent
	claim := &Claim{
		UID:       prevClaim.UID,
		Hash:      dataHash,
		ParentUID: prevClaimHash,
	}
	return cs.StoreClaim(ctx, claim)
}

// WalkClaimChain walks through the chain of claims backward starting from the given claim hash,
// returning a list of all claims up to the initial one (with no PrevClaimHash).
func (cs *ClaimsService) WalkClaimChain(ctx context.Context, latestClaimHash string) ([]*Claim, error) {
	var chain []*Claim
	currentHash := latestClaimHash
	for currentHash != "" {
		c, err := cs.RetrieveClaim(ctx, currentHash)
		if err != nil {
			return nil, err
		}
		chain = append(chain, c)
		currentHash = c.ParentUID
	}
	// Now chain[0] is the latest claim, chain[len(chain)-1] is the initial claim
	// If you want them in chronological order, reverse the slice
	slices.Reverse(chain)

	return chain, nil
}

// GetLatestClaimForEntity retrieves the latest claim for an entity by scanning through known claims.
// NOTE: Without indexing, we must rely on external knowledge of a "latest" claim hash.
// In a real system, you might store a known claim hash for the entity and walk forward or rely on indexing.
// This function is just a placeholder to illustrate usage.
//
// One approach is that if you somehow have a candidate claim hash (from a known source or hint),
// you can walk backward to verify the chain. Without an index or a pointer-like system, you cannot
// easily "discover" the latest claim just from the CAS.
//
// If you had a list of claim hashes, you could pick the one with the highest version after retrieval.
// That would simulate an indexing step, which is outside the scope of this snippet.
func (cs *ClaimsService) GetLatestClaimForEntity(ctx context.Context, candidateClaimHashes []string) (*Claim, error) {
	if len(candidateClaimHashes) == 0 {
		return nil, errors.New("no candidate claims provided")
	}

	var latest *Claim
	for _, h := range candidateClaimHashes {
		c, err := cs.RetrieveClaim(ctx, h)
		if err != nil {
			continue
		}
		if latest == nil || c.Timestamp > latest.Timestamp {
			latest = c
		}
	}

	if latest == nil {
		return nil, errors.New("no valid claims found")
	}
	return latest, nil
}
