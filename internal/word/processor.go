package word

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diptanw/countd-zmq/internal/platform/async"
)

// Encoder is an interface for data encoder.
type Encoder interface {
	Encode(v interface{}) error
}

// Processor is a service that aggregates word count and schedules the stream
// writer based on flushing strategy. It marshals the data before flushing using
// the provided encoder.
type Processor struct {
	repo      *Repository
	encoder   Encoder
	hotList   map[string]struct{}
	scheduler *async.Scheduler
}

// NewProcessor returns a new instance of Processor.
func NewProcessor(repo *Repository, encoder Encoder, scheduler *async.Scheduler, hotList []string) Processor {
	m := make(map[string]struct{}, len(hotList))
	for _, word := range hotList {
		m[word] = struct{}{}
	}

	return Processor{
		repo:      repo,
		encoder:   encoder,
		scheduler: scheduler,
		hotList:   m,
	}
}

// Aggregate word count and flush it on schedule:
//   - When a never-before-seen word is encountered, it is flushed immediately.
//   - Counting from the time at which a word is first seen, its aggregates are
//     flushed at 1s intervals. Note that this implies that most aggregates are
//     not flushed at the same time, independently of each other.
//   - Words in the hot list are flushed at 10s intervals, instead of the 1s
//     default.
func (p Processor) Aggregate(ctx context.Context, word string) error {
	aggr, err := p.repo.Count(word)
	if err != nil {
		return fmt.Errorf("word: aggregating count: %w", err)
	}

	if aggr == nil || aggr.ID() == "" || aggr.LastSeenAt.IsZero() {
		return errors.New("word: invalid aggregate")
	}

	if aggr.TotalCount <= 1 {
		if err := p.flush(aggr.Word); err != nil {
			return err
		}

		p.scheduler.Schedule(ctx, p.flushInterval(aggr), func(_ context.Context) error {
			return p.flush(aggr.Word)
		})
	}

	return nil
}

func (p Processor) flush(word string) error {
	aggr, err := p.repo.Reset(word)
	if err != nil {
		return err
	}

	return p.encoder.Encode(aggr)
}

func (p Processor) flushInterval(aggr *Aggregate) time.Duration {
	if _, hot := p.hotList[aggr.Word]; hot {
		return 10 * time.Second
	}

	return time.Second
}
