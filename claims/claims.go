package claims

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"go.quinn.io/dataq/cas"
)

var PREFIX = byte('{')

// Claim represents a versioned assertion about an entity.
// It includes the entity's unique identifier, a reference to the previous claim,
// and a reference to the associated data.
type Claim struct {
	DQV       string `json:"dataq_version"`     // version of dataq claim schema
	UID       string `json:"entity_uid"`        // unique identifier of the entity. shared across multiple claims, and used to reduce all claims into a single entity.
	ParentUID string `json:"entity_parent_uid"` // parent entity UID. This is linked through plugin requests and responses
	Timestamp int64  `json:"claim_timestamp"`   // Unix timestamp of creation
	Kind      string `json:"schema_kind"`       // kind of the entity. Has meaning to DataQ. Determines table of index
	Payload   any    `json:"payload"`           // payload of the entity. This what is stored in the Index
}

type ClaimPayload struct {
	Kind    string          `json:"schema_kind"`
	Payload json.RawMessage `json:"payload"`
}

type Claimer interface {
	StoreClaim(context.Context, *Claim) (hash string, err error)
	RetrieveClaim(ctx context.Context, hash string) (*Claim, *ClaimPayload, error)
}

type Claimable interface {
	GetUid() string
	GetClaimTimestamp() int64
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
	c.DQV = "1"
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
func (cs *ClaimsService) RetrieveClaim(ctx context.Context, hash string) (*Claim, *ClaimPayload, error) {
	r, err := cs.cas.Retrieve(ctx, hash)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read claim bytes: %w", err)
	}

	if len(bytes) == 0 || bytes[0] != PREFIX {
		return nil, nil, errors.New("not a claim")
	}

	var c Claim
	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, nil, err
	}

	if c.DQV != "1" || c.UID == "" || c.Timestamp == 0 {
		return nil, nil, errors.New("invalid claim structure")
	}

	var p ClaimPayload
	if err := json.Unmarshal(bytes, &p); err != nil {
		return nil, nil, err
	}

	return &c, nil, nil
}
