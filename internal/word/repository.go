package word

import (
	"errors"
	"sync"
	"time"

	"github.com/diptanw/countd-zmq/internal/platform/storage"
)

// Aggregate is a type that represents the word aggregate details.
type Aggregate struct {
	Word        string    `json:"word"`        // Word is the aggregate group key.
	LastSeenAt  time.Time `json:"last_seen"`   // When the word was first seen.
	FirstSeenAt time.Time `json:"first_seen"`  // When the word was last seen.
	TotalCount  int       `json:"total_count"` // How many times the word has been seen since this aggregate was last flushed.
	DeltaCount  int       `json:"delta_count"` // How many times the word has been seen since the service was started.
}

// ID is the record identifier in the data store.
func (c *Aggregate) ID() storage.ID {
	return storage.ID(c.Word)
}

// Repository is a type that aggregates words counting.
type Repository struct {
	store   *storage.InMemory[*Aggregate]
	storeMu sync.Mutex
}

// NewRepository creates anew instance of Repository.
func NewRepository(db *storage.InMemory[*Aggregate]) *Repository {
	return &Repository{
		store: db,
	}
}

// Count update the aggregate record with the word total/delta count
// and fist/last seen timestamps.
func (r *Repository) Count(word string) (*Aggregate, error) {
	if word == "" {
		return nil, errors.New("empty word value")
	}

	r.storeMu.Lock()
	defer r.storeMu.Unlock()

	aggr, now := &Aggregate{
		Word: word,
	}, time.Now()

	rec, err := r.store.Get(storage.ID(word))
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return nil, err
		}

		// If record does not exist, set the first seen time.
		aggr.FirstSeenAt = now
	}

	if rec != nil {
		aggr = rec
	}

	aggr.DeltaCount++
	aggr.TotalCount++
	aggr.LastSeenAt = now

	if err := r.store.Update(aggr); err != nil {
		return nil, err
	}

	return aggr, nil
}

func (r *Repository) Reset(word string) (*Aggregate, error) {
	if word == "" {
		return nil, errors.New("empty word value")
	}

	r.storeMu.Lock()
	defer r.storeMu.Unlock()

	rec, err := r.store.Get(storage.ID(word))
	if err != nil {
		return nil, err
	}

	var retAggr Aggregate

	if rec != nil {
		retAggr = *rec
		rec.DeltaCount = 0
	}

	if err := r.store.Update(rec); err != nil {
		return nil, err
	}

	return &retAggr, nil
}
