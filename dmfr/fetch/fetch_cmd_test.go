package fetch

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/interline-io/transitland-lib/internal/testdb"
	"github.com/interline-io/transitland-lib/internal/testutil"
	"github.com/interline-io/transitland-lib/tl"
)

func TestFetchCommand(t *testing.T) {
	ts200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadFile(testutil.ExampleZip.URL)
		if err != nil {
			t.Error(err)
		}
		w.Write(buf)
	}))
	defer ts200.Close()
	ts404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Status-Code", "404")
		w.Write([]byte("not found"))
	}))
	defer ts404.Close()
	// note - Spec==gtfs is required for fetch
	f200 := tl.Feed{FeedID: "f--200", Spec: "gtfs", URLs: tl.FeedUrls{StaticCurrent: ts200.URL}}
	f404 := tl.Feed{FeedID: "f--404", Spec: "gtfs", URLs: tl.FeedUrls{StaticCurrent: ts404.URL}}
	cases := []struct {
		fvcount     int
		errContains string
		feeds       []tl.Feed
		command     []string
	}{
		{1, "", []tl.Feed{f200}, []string{}},
		{1, "", []tl.Feed{f200, f404}, []string{"f--200", "f--404"}},
		{1, "", []tl.Feed{f200, f404}, []string{"f--200"}},
		{0, "", []tl.Feed{f200, f404}, []string{"f--404"}},
	}
	_ = cases
	for _, exp := range cases {
		t.Run("", func(t *testing.T) {
			adapter := testdb.MustOpenWriter("sqlite3://:memory:", true).Adapter
			for _, feed := range exp.feeds {
				testdb.ShouldInsert(t, adapter, &feed)
			}
			c := Command{adapter: adapter}
			tmpDir := t.TempDir()
			withTempDir := []string{"-storage", tmpDir}
			withTempDir = append(withTempDir, exp.command...)
			if err := c.Parse(withTempDir); err != nil {
				t.Fatal(err)
			}
			if err := c.Run(); err != nil && exp.errContains != "" {
				if !strings.Contains(err.Error(), exp.errContains) {
					t.Errorf("got '%s' error, expected to contain '%s'", err.Error(), exp.errContains)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			// Test
			feeds := []tl.Feed{}
			testdb.ShouldSelect(t, adapter, &feeds, "SELECT * FROM current_feeds")
			if len(feeds) != len(exp.feeds) {
				t.Errorf("got %d feeds, expect %d", len(feeds), len(exp.feeds))
			}
			fvs := []tl.FeedVersion{}
			testdb.ShouldSelect(t, adapter, &fvs, "SELECT * FROM feed_versions")
			if len(fvs) != exp.fvcount {
				t.Errorf("got %d feed versions, expect %d", len(fvs), exp.fvcount)
			}
			for _, fv := range fvs {
				fn := filepath.Join(tmpDir, fv.File)
				// fn := fv.File
				st, err := os.Stat(fn)
				if err != nil {
					t.Errorf("got '%s', expected file '%s' to exist", err.Error(), fn)
				} else {
					// TODO: Check SHA1
					if st.Size() == 0 {
						t.Errorf("expected non-empty file")
					}
				}
			}
		})
	}
}
