package sync

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/interline-io/transitland-lib/dmfr"
	"github.com/interline-io/transitland-lib/stats"
	"github.com/interline-io/transitland-lib/tldb"
	"github.com/interline-io/transitland-lib/tt"
)

// UpdateFeed .
func UpdateFeed(ctx context.Context, atx tldb.Adapter, rfeed dmfr.Feed) (int, bool, bool, error) {
	// Check if we have the existing Feed
	feedid := 0
	found := false
	updated := false
	var errTx error
	dbfeed := dmfr.Feed{}
	if err := atx.Get(ctx, &dbfeed, "select * from current_feeds where onestop_id = ?", rfeed.FeedID); err == nil {
		// Exists, update key values
		found = true
		feedid = dbfeed.ID
		rfeed.ID = dbfeed.ID
		if !dbfeed.Equal(&rfeed) {
			updated = true
			rfeed.CreatedAt = dbfeed.CreatedAt
			rfeed.DeletedAt = tt.Time{}
			errTx = atx.Update(ctx, &rfeed)
		}
	} else if err == sql.ErrNoRows {
		feedid, errTx = atx.Insert(ctx, &rfeed)
	} else {
		// Error
		errTx = err
	}
	// Create feed state if not exists
	if _, err := stats.GetFeedState(ctx, atx, feedid); err != nil {
		errTx = err
	}
	return feedid, found, updated, errTx
}

// HideUnseedFeeds .
func HideUnseedFeeds(ctx context.Context, atx tldb.Adapter, found []int) (int, error) {
	// Delete unreferenced feeds
	t := tt.NewTime(time.Now().UTC())
	r, err := atx.Sqrl().
		Update("current_feeds").
		Where(sq.NotEq{"id": found}).
		Where(sq.Eq{"deleted_at": nil}).
		Set("deleted_at", t).
		Exec()
	if err != nil {
		return 0, err
	}
	c, err := r.RowsAffected()
	return int(c), err
}

// UpdateFeedGeneratedOperators creates OperatorInFeed records for agencies that are not associated with an operator
func UpdateFeedGeneratedOperators(ctx context.Context, atx tldb.Adapter, found []int) error {
	for _, id := range found {
		feed := dmfr.Feed{}
		if err := atx.Get(ctx, &feed, "select * from current_feeds where id = ?", id); err != nil {
			return err
		}
		if _, err := feedUpdateOifs(ctx, atx, feed); err != nil {
			return err
		}
	}
	return nil
}
