# Interline Transitland <!-- omit in toc -->

`transitland-lib` is a library and command-line tool for reading, writing, and processing transit data in [GTFS](http://gtfs.org) and related formats. The library is structured as a set of data sources, filters, and transformations that can be mixed together in a variety of ways to create processing pipelines. The library supports the [DMFR](https://github.com/transitland/distributed-mobility-feed-registry) format to describe feed resources.

![Test & Release](https://github.com/interline-io/transitland-lib/workflows/Test%20&%20Release/badge.svg) [![GoDoc](https://godoc.org/github.com/interline-io/transitland-lib/tl?status.svg)](https://godoc.org/github.com/interline-io/transitland-lib/tl) ![Go Report Card](https://goreportcard.com/badge/github.com/interline-io/transitland-lib)

## Table of Contents <!-- omit in toc -->
<!-- to update use https://marketplace.visualstudio.com/items?itemName=yzhang.markdown-all-in-one -->
- [Installation](#installation)
	- [Download prebuilt binary](#download-prebuilt-binary)
	- [Install using homebrew](#install-using-homebrew)
	- [Install binary from source](#install-binary-from-source)
- [Usage as a CLI tool](#usage-as-a-cli-tool)
	- [`validate` command](#validate-command)
	- [`copy` command](#copy-command)
	- [`extract` command](#extract-command)
	- [`sync` command](#sync-command)
	- [`fetch` command](#fetch-command)
	- [`import` command](#import-command)
	- [`unimport` command](#unimport-command)
	- [`dmfr` command](#dmfr-command)
- [Usage as a library](#usage-as-a-library)
	- [Key library components](#key-library-components)
	- [Install as a library](#install-as-a-library)
	- [Example of how to use as a library](#example-of-how-to-use-as-a-library)
- [Database migrations](#database-migrations)
- [Usage as a Web Service](#usage-as-a-web-service)
	- [`transitland_server` command](#transitland_server-command)
	- [Hasura](#hasura)
- [Included Readers and Writers](#included-readers-and-writers)
- [Development](#development)
	- [Releases](#releases)
- [Licenses](#licenses)


## Installation

### Download prebuilt binary

The `transitland` binaries for Linux and macOS are attached to each [release](https://github.com/interline-io/transitland-lib/releases).

### Install using homebrew

The `transitland` binary can be installed using homebrew. The executable is code-signed and notarized.

```bash
brew install interline-io/transitland-lib/transitland-lib
```

### Install binary from source

```bash
go get github.com/interline-io/transitland-lib/cmd/transitland
```

This package uses Go Modules and will also install required dependencies.

Main dependencies:
- `twpayne/go-geom`
- `jmoiron/sqlx`
- `Masterminds/squirrel`
- `lib/pq`
- `mattn/go-sqlite3` (requires CGO)

## Usage as a CLI tool

The main subcommands are:
- [validate](#validate-command)
- [copy](#copy-command)
- [extract](#extract-command)
- [sync](#sync-command)
- [fetch](#fetch-command)
- [import](#import-command)
- [unimport](#unimport-command)- 
- [dmfr](#dmfr-command)

### `validate` command

The validate command performs a basic validation on a data source and writes the results to standard out.

```bash
% transitland validate --help
Usage: validate <reader>
  -best-practices
    	Include Best Practices validations
  -ext value
    	Include GTFS Extension
  -o string
    	Write validation report as JSON to file
```

Example: 

```bash
% transitland validate "https://www.bart.gov/dev/schedules/google_transit.zip"
```

### `copy` command

The copy command performs a basic copy from a reader to a writer. By default, any entity with errors will be skipped and not written to output. This can be ignored with `-allow-entity-errors` to ignore simple errors and `-allow-reference-errors` to ignore entity relationship errors, such as a reference to a non-existent stop.

```bash
% transitland copy --help
Usage: copy <reader> <writer>
  -allow-entity-errors
    	Allow entities with errors to be copied
  -allow-reference-errors
    	Allow entities with reference errors to be copied
  -create
    	Create a basic database schema if none exists
  -ext value
    	Include GTFS Extension
  -fvid int
    	Specify FeedVersionID when writing to a database
```

Example:

```bash
% transitland copy --allow-entity-errors "https://www.bart.gov/dev/schedules/google_transit.zip" output.zip

% unzip -p output.zip agency.txt
agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url,agency_email
BART,Bay Area Rapid Transit,https://www.bart.gov/,America/Los_Angeles,,510-464-6000,,
  ```

### `extract` command

The extract command extends the basic copy command with a number of additional options and transformations. It can be used to pull out a single route or trip, interpolate stop times, override a single value on an entity, etc. This is a separate command to keep the basic copy command simple while allowing the extract command to grow and add more features over time.

```bash
% transitland extract --help
Usage: extract <input> <output>
  -allow-entity-errors
    	Allow entities with errors to be copied
  -allow-reference-errors
    	Allow entities with reference errors to be copied
  -create
    	Create a basic database schema if none exists
  -create-missing-shapes
    	Create missing Shapes from Trip stop-to-stop geometries
  -deduplicate-stop-times
    	Deduplicate StopTimes using Journey Patterns
  -ext value
    	Include GTFS Extension
  -extract-agency value
    	Extract Agency
  -extract-calendar value
    	Extract Calendar
  -extract-route value
    	Extract Route
  -extract-route-type value
    	Extract Routes matching route_type
  -extract-stop value
    	Extract Stop
  -extract-trip value
    	Extract Trip
  -fvid int
    	Specify FeedVersionID when writing to a database
  -interpolate-stop-times
    	Interpolate missing StopTime arrival/departure values
  -normalize-service-ids
    	Create any missing Calendar entities for CalendarDate service_ids
  -normalize-timezones
    	Normalize timezones and apply default stop timezones based on agency and parent stops
  -set value
    	Set values on output; format is filename,id,key,value
  -simplify-calendars
    	Attempt to simplify CalendarDates into regular Calendars
  -simplify-shapes float
    	Simplify shapes with this tolerance (ex. 0.000005)
  -use-basic-route-types
    	Collapse extended route_type's into basic GTFS values
```

Example:

```bash
# Extract a single trip from the BART GTFS, and rename the agency to "test".
% transitland extract -extract-trip "3050453" -set "agency.txt,BART,agency_id,test" "https://www.bart.gov/dev/schedules/google_transit.zip" output2.zip

# Note renamed agency
% unzip -p output2.zip agency.txt
agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url,agency_email
test,Bay Area Rapid Transit,https://www.bart.gov/,America/Los_Angeles,,510-464-6000,,

# Only entities related to the specified trip are included in the output.
% unzip -p output2.zip trips.txt
route_id,service_id,trip_id,trip_headsign,trip_short_name,direction_id,block_id,shape_id,wheelchair_accessible,bikes_allowed
1,2020_09_14-DX-MVS-Weekday-15,3050453,San Francisco International Airport,,1,,01_shp,0,0

$ unzip -p output2.zip routes.txt
route_id,agency_id,route_short_name,route_long_name,route_desc,route_type,route_url,route_color,route_text_color,route_sort_order
1,test,YL-S,Antioch to SFIA/Millbrae,,1,http://www.bart.gov/schedules/bylineresults?route=1,FFFF33,,0

% transitland % unzip -p output2.zip stop_times.txt
trip_id,arrival_time,departure_time,stop_id,stop_sequence,stop_headsign,pickup_type,drop_off_type,shape_dist_traveled,timepoint
3050453,04:53:00,04:53:00,CONC,0,,0,0,0.00000,0
3050453,04:58:00,04:58:00,PHIL,2,,0,0,4.06000,0
3050453,05:01:00,05:02:00,WCRK,3,,0,0,5.77000,0
3050453,05:06:00,05:07:00,LAFY,4,,0,0,9.23000,0
3050453,05:11:00,05:12:00,ORIN,5,,0,0,12.99000,0
3050453,05:17:00,05:18:00,ROCK,6,,0,0,17.38000,0
...
```

### `sync` command

The `sync` command loads DMFR files into a database.

```bash
% transitland sync --help
Usage: sync <Filenames...>
  -dburl string
    	Database URL (default: $TL_DATABASE_URL)
  -hide-unseen
    	Hide unseen feeds
  -hide-unseen-operators
    	Hide unseen operators
```

### `fetch` command

The `fetch` command fetches GTFS data and saves feed version records to the database. Use after the `sync` command.

```bash
% transitland fetch --help
Usage: fetch [feed_id...]
  -allow-ftp-fetch
    	Allow fetching from FTP urls
  -allow-local-fetch
    	Allow fetching from filesystem directories/zip files
  -allow-s3-fetch
    	Allow fetching from S3 urls
  -create-feed
    	Create feed records if not found
  -dburl string
    	Database URL (default: $TL_DATABASE_URL)
  -dry-run
    	Dry run; print feeds that would be imported and exit
  -feed-url string
    	Manually fetch a single URL; you must specify exactly one feed_id
  -fetched-at string
    	Manually specify fetched_at value, e.g. 2020-02-06T12:34:56Z
  -storage string
    	GTFS storage location
  -ignore-duplicate-contents
    	Allow duplicate internal SHA1 contents
  -limit int
    	Maximum number of feeds to fetch
  -secrets string
    	Path to DMFR Secrets file
  -workers int
    	Worker threads (default 1)
```

### `import` command

The `import` command imports previously fetched feed versions into the database. Use after the `fetch` command.

```bash
% transitland import --help
Usage: import [feedids...]
  -activate
    	Set as active feed version after import
  -create-missing-shapes
    	Create missing Shapes from Trip stop-to-stop geometries
  -date string
    	Service on date
  -dburl string
    	Database URL (default: $TL_DATABASE_URL)
  -deduplicate-stop-times
    	Deduplicate StopTimes using Journey Patterns
  -dryrun
    	Dry run; print feeds that would be imported and exit
  -ext value
    	Include GTFS Extension
  -fetched-since string
    	Fetched since
  -fv-sha1 value
    	Feed version SHA1
  -fv-sha1-file string
    	Specify feed version IDs by SHA1 in file, one per line
  -fvid value
    	Import specific feed version ID
  -fvid-file string
    	Specify feed version IDs in file, one per line; equivalent to multiple --fvid
  -storage string
    	GTFS storage location
  -interpolate-stop-times
    	Interpolate missing StopTime arrival/departure values
  -latest
    	Only import latest feed version available for each feed
  -limit int
    	Import at most n feeds
  -normalize-timezones
    	Normalize timezones and apply default stop timezones based on agency and parent stops
  -simplify-calendars
    	Attempt to simplify CalendarDates into regular Calendars
  -simplify-shapes float
    	Simplify shapes with this tolerance (ex. 0.000005)
  -workers int
    	Worker threads (default 1)
```

### `unimport` command

The `unimport` command deletes previously imported data from feed versions. The feed version record itself is not deleted. You may optionally specify removal of only schedule data, leaving routes, stops, etc. in place.

```bash
% transitland unimport --help
Usage: unimport [fvids]
  -dburl string
    	Database URL (default: $TL_DATABASE_URL)
  -dryrun
    	Dry run; print feeds that would be imported and exit
  -feed value
    	Feed ID
  -fv-sha1 value
    	Feed version SHA1
  -fv-sha1-file string
    	Specify feed version IDs by SHA1 in file, one per line
  -fvid-file string
    	Specify feed version IDs in file, one per line; equivalent to multiple --fvid
  -schedule-only
    	Unimport stop times, trips, transfers, shapes, and frequencies
```

### `dmfr` command

The `dmfr` command enables validation, linting, and formatting of [Distributed Mobility Feed Registry]([dmfr](https://github.com/transitland/distributed-mobility-feed-registry)) format. Please see [DMFR Command help](dmfr-command.md).

## Usage as a library

### Key library components

- Entity: An `Entity` is entity as specified by GTFS, such as an Agency, Route, Stop, etc.
- Reader: A `Reader` provides streams of GTFS entities over channels. The `tlcsv` and `tldb` modules provide CSV and PostgreSQL/SQLite support, respectively.
- Writer: A `Writer` accepts GTFS entities. As above, `tlcsv` and `tldb` provide basic implementations. Custom writers can also be used to support non-GTFS outputs, such as building a routing graph.
- Copier: A `Copier` reads a stream of GTFS entities from a `Reader`, checks each entity against a `Marker`, performs validation, applies any specified `Filters`, and sends to a `Writer`.
- Marker: A `Marker` selects which GTFS entities will be processed by a `Copier`. For example, selecting only entities related to a single trip or route.
- Filter: A `Filter` applies transformations to GTFS entities, such as converting extended route types to basic values, or modifying entity identifiers.
- Extension: An `Extension` provides support for additional types of GTFS entities.

See [godoc.org](https://godoc.org/github.com/interline-io/transitland-lib/tl) for package documentation.

### Install as a library

```bash
go get github.com/interline-io/transitland-lib
```

### Example of how to use as a library

A simple example of reading and writing GTFS entities from CSV ([full example](https://github.com/interline-io/transitland-lib/raw/master/internal/testreadme/main_test.go)):

```go
package main

import (
	"fmt"
	"testing"

	"github.com/interline-io/transitland-lib/copier"
	"github.com/interline-io/transitland-lib/tl"
	"github.com/interline-io/transitland-lib/tlcsv"
	"github.com/interline-io/transitland-lib/tldb"
)

var URL = "https://github.com/interline-io/transitland-lib/raw/master/test/data/external/bart.zip"

func TestExample1(t *testing.T) {
	// Read stops from a GTFS url
	reader, _ := tlcsv.NewReader(URL)
	reader.Open()
	defer reader.Close()
	// Write to to the current directory
	writer, _ := tlcsv.NewWriter(".")
	writer.Open()
	// Copy stops
	for stop := range reader.Stops() {
		fmt.Println("Read Stop:", stop.StopID)
		eid, _ := writer.AddEntity(&stop)
		fmt.Println("Wrote stop:", eid)
	}
}
```

Database support is handled similary:

```go
func getReader() tl.Reader {
	reader, _ := tlcsv.NewReader(URL)
	return reader
}

func TestExample2(t *testing.T) {
	reader := getReader()
	// Create a SQLite writer, in memory
	dburl := "sqlite3://:memory:"
	dbwriter, err := tldb.NewWriter(dburl)
	if err != nil {
		t.Fatalf("no reader available")
	}
	if err := dbwriter.Open(); err != nil {
		t.Fatalf("could not open database for writing")
	}
	if err := dbwriter.Create(); err != nil {
		t.Fatalf("could not find or create database schema")
	}
	for stop := range reader.Stops() {
		// A database writer AddEntity returns the primary key as a string.
		fmt.Println("Read Stop:", stop.StopID)
		eid, err := dbwriter.AddEntity(&stop)
		if err != nil {
			t.Fatalf("could not write entity to database")
		}
		fmt.Println("wrote stop to database:", eid)
	}
	// Read back from this source.
	dbreader, err := dbwriter.NewReader()
	if err != nil {
		t.Fatalf("could not get a new reader")
	}
	count := 0
	for stop := range dbreader.Stops() {
		fmt.Println("read stop from database:", stop.StopID)
		count++
	}
	if count != 50 {
		t.Errorf("got %d stops, expected 50", count)
	}
}
```

More advanced filtering operations can be performed using a `Copier`, which provides additional hooks for filtering, transformation, and validation:

```go
func TestExample3(t *testing.T) {
	reader := getReader()
	// Create a zip writer
	writer, err := tlcsv.NewWriter("filtered.zip")
	if err != nil {
		t.Fatalf("no writer available")
	}
	// Create a copier to stream, filter, and validate entities
	cp := copier.NewCopier(reader, writer)
	result := cp.Copy()
	if result.WriteError != nil {
		t.Fatalf("fatal copy error")
	}
	for _, err := range result.Errors {
		fmt.Println("error:", err)
	}
	for fn, count := range result.EntityCount {
		fmt.Printf("copied %d entities from %s\n", count, fn)
	}
}
```

See API docs at https://godoc.org/github.com/interline-io/transitland-lib

## Database migrations

Migrations are supported for PostgreSQL, using the schema files in `internal/schema/postgres/migrations`. These files can be read and applied using [golang-migrate](https://github.com/golang-migrate/migrate), which will store the most recently applied migration version in `schema_migrations`. See the `bootstrap.sh` script in that directory for an example, as well as details on how to import Natural Earth data files for associating agencies with places.

SQLite database are intended to be short-lived. They can be created on an as needed basis by passing the `-create` flag to some commands that accept a writer. They use a single executable schema, defined in `internal/schema/sqlite.sql`.

## Usage as a Web Service

`transitland-lib` can be used in a variety of ways to power a web service. Interline currently uses two approaches:

1. Populate a database with one or more feeds using `transitland-lib` and use the `transitland-server` package to serve the Transitland v2 REST and/or v2 GraphQL API endpoints. These API endpoints are primarily read-only and focused on querying and analyzing transit data.

2. Populate a postgres database with one or more feeds using `transitland-lib` and use [Hasura](https://hasura.io/) to provide a basic GraphQL API for reading and writing into the database. 

For more information about how these web services are used within the overall architecture of the Transitland platform, see https://www.transit.land/documentation#transitland-architecture 

### `transitland_server` command

See [transitland-server](https://github.com/interline-io/transitland-server) documentation.

### Hasura

[Hasura](https://hasura.io/) is a web service that can provide an "instant" GraphQL API based on a postgres database and its schema. We combine Hasura with `transitland-lib` for projects that involve creating new or complex queries (since Hasura can be more flexible than the queries provided by `transitland server`) and projects that involve an API with full read and write access (for example, editing GTFS data, which is also not provided by `transitland server`). Note that Hasura's automatically generated database queries are not guaranteed to be efficient (on the other hand, `transitland server` is tuned to provide better performance).

To use Hasura with `transitland-lib` you can either import feeds into an existing postgres database (using the `transitland dmfr` command) and configure Hasura to recognize all the tables and the foreign key relationships between them.

## Included Readers and Writers

| Target                   | Module  | Supports Read | Supports Write |
| ------------------------ | ------- | ------------- | -------------- |
| CSV                      | `tlcsv` | ✅             | ✅              |
| SQLite                   | `tldb`  | ✅             | ✅              |
| PostgreSQL (with PostGIS)  | `tldb`  | ✅             | ✅              |

We welcome the addition of more readers and writers.

## Development

`transitland-lib` follows Go coding conventions.

GitHub Actions runs all tests, stores code coverage reports as artifacts, and prepares releases.

### Releases

Releases follow [Semantic Versioning](https://semver.org/) conventions.

To cut a new release:

1. Update `transitland-lib/tl/tl.go` with the new version.
2. Create a GitHub release. This will create a tag and GitHub Actions will create &amp; attach code-signed binaries.
3. Download the files from the release, and update the [homebrew formula](https://github.com/interline-io/homebrew-transitland-lib/blob/master/transitland-lib.rb) with the updated sha256 hashes and version tag.

## Licenses

`transitland-lib` is released under a "dual license" model:

- open-source for use by all under the [GPLv3](LICENSE) license
- also available under a flexible commercial license from [Interline](mailto:info@interline.io)

