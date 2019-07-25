package gtdb

import (
	"os"
	"testing"

	"github.com/interline-io/gotransit/copier"
	"github.com/interline-io/gotransit/internal/testutil"
)

func testReader(t *testing.T, adapter Adapter) {
	writer := Writer{Adapter: adapter}
	if err := adapter.Open(); err != nil {
		t.Error(err)
	}
	if err := adapter.Create(); err != nil {
		t.Error(err)
	}
	if err := adapter.Create(); err != nil {
		t.Error(err)
	}
	defer writer.Close()
	// Get mock reader
	me, reader := testutil.NewMinimalExpect()
	// Create FeedVersion - required for foreign key constraints
	if _, err := writer.CreateFeedVersion(reader); err != nil {
		t.Error(err)
	}
	cp := copier.NewCopier(reader, &writer)
	result := cp.Copy()
	if len(result.Errors) > 0 {
		t.Error(result.Errors[0])
	}
	// Now test
	r2, err := writer.NewReader()
	if err != nil {
		t.Error(err)
	}
	testutil.TestExpect(t, *me, r2)
}

// Reader interface tests.

func TestReader_Postgres(t *testing.T) {
	dburl := os.Getenv("GOTRANSIT_TEST_POSTGRES_URL")
	if len(dburl) == 0 {
		t.Skip()
		return
	}
	testReader(t, &PostgresAdapter{DBURL: dburl})
}

func TestReader_SpatiaLite(t *testing.T) {
	dburl := "sqlite3://test.db"
	testReader(t, &SpatiaLiteAdapter{DBURL: dburl})
}
