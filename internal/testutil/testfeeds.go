package testutil

// ExampleFeed .
var ExampleFeed = ReaderTester{
	URL: "../testdata/example",
	Counts: map[string]int{
		"agency.txt":          1,
		"routes.txt":          5,
		"trips.txt":           11,
		"stops.txt":           9,
		"stop_times.txt":      28,
		"shapes.txt":          3,
		"calendar.txt":        2,
		"calendar_dates.txt":  2,
		"feed_info.txt":       1,
		"fare_rules.txt":      4,
		"fare_attributes.txt": 2,
		"frequency.txt":       11,
		"transfers.txt":       0,
	},
	EntityIDs: map[string][]string{
		"agency.txt":          {"DTA"},
		"routes.txt":          {"AB", "BFC", "STBA", "CITY", "AAMV"},
		"trips.txt":           {"AB1", "AB2", "STBA", "CITY1", "CITY2", "BFC1", "BFC2", "AAMV1", "AAMV2", "AAMV3", "AAMV4"},
		"stops.txt":           {"FUR_CREEK_RES", "BULLFROG"}, // partial
		"shapes.txt":          {"ok", "a", "c"},
		"calendar.txt":        {"FULLW", "WE"},
		"fare_attributes.txt": {"p", "a"},
	},
}

// ExternalTestFeed .
func ExternalTestFeed(key string) (ReaderTester, bool) {
	a, ok := ExternalTestFeeds[key]
	return a, ok
}

// ExternalTestFeeds -- Generated from above commented out code
var ExternalTestFeeds = map[string]ReaderTester{
	"bart.zip": ReaderTester{
		URL: "../testdata/external/bart.zip",
		Counts: map[string]int{"agency.txt": 1,
			"calendar.txt":        3,
			"calendar_dates.txt":  12,
			"fare_attributes.txt": 170,
			"fare_rules.txt":      2304,
			"feed_info.txt":       1,
			"routes.txt":          6,
			"shapes.txt":          12,
			"stop_times.txt":      33167,
			"stops.txt":           50,
			"transfers.txt":       8,
			"trips.txt":           2525},
		EntityIDs: map[string][]string{
			"agency.txt":          []string{"BART"},
			"calendar.txt":        []string{"WKDY", "SAT", "SUN"},
			"fare_attributes.txt": []string{"50", "51", "52", "53", "54", "56", "57", "58", "59", "60"},
			"routes.txt": []string{"01",
				"03",
				"05",
				"07",
				"11",
				"19"},
			"shapes.txt": []string{"01_shp",
				"02_shp",
				"03_shp",
				"04_shp",
				"05_shp",
				"06_shp",
				"07_shp",
				"08_shp",
				"11_shp",
				"12_shp"},
			"stops.txt": []string{"12TH",
				"16TH",
				"19TH",
				"19TH_N",
				"24TH",
				"ANTC",
				"ASHB",
				"BALB",
				"BAYF",
				"CAST"},
			"trips.txt": []string{"3610403WKDY",
				"3730559SAT",
				"3650800SUN",
				"3630418WKDY",
				"3750558SAT",
				"3670758SUN",
				"3650433WKDY",
				"3770618SAT",
				"3690818SUN",
				"3670448WKDY"}}},

	"cdmx.zip": ReaderTester{
		URL: "../testdata/external/cdmx.zip",
		Counts: map[string]int{
			"agency.txt":      8,
			"calendar.txt":    99,
			"frequencies.txt": 1741,
			"routes.txt":      145,
			"shapes.txt":      345,
			"stop_times.txt":  70714,
			"stops.txt":       6021,
			"trips.txt":       1741},
		EntityIDs: map[string][]string{
			"agency.txt": []string{
				"CC",
				"MB",
				"METRO",
				"NCC",
				"RTP",
				"RTP_ESP",
				"STE",
				"SUB"},
			"calendar.txt": []string{
				"14741",
				"16092",
				"16203",
				"16233",
				"26656",
				"26949",
				"28960",
				"36238",
				"36384",
				"36479"},
			"routes.txt": []string{
				"ROUTE_132162",
				"ROUTE_132207",
				"ROUTE_134667",
				"ROUTE_18226",
				"ROUTE_36490",
				"ROUTE_36644",
				"ROUTE_36712",
				"ROUTE_36713",
				"ROUTE_136284",
				"ROUTE_136285"},
			"shapes.txt": []string{
				"14816",
				"14817",
				"14818",
				"14819",
				"14820",
				"14821",
				"14822",
				"14823",
				"14824",
				"14825"},
			"stops.txt": []string{
				"136300",
				"136299",
				"28503",
				"28800",
				"28805",
				"28812",
				"28830",
				"28832",
				"28833",
				"28963"},
			"trips.txt": []string{
				"14743",
				"14840",
				"14841",
				"14842",
				"14843",
				"14844",
				"14845",
				"14846",
				"14848",
				"15171"}}},
	"mbta.zip": ReaderTester{
		URL: "../testdata/external/mbta.zip",
		Counts: map[string]int{
			"agency.txt":         2,
			"calendar.txt":       122,
			"calendar_dates.txt": 31,
			"feed_info.txt":      1,
			"routes.txt":         216,
			"shapes.txt":         897,
			"stop_times.txt":     1291946,
			"stops.txt":          9838,
			"transfers.txt":      1834,
			"trips.txt":          55166},
		EntityIDs: map[string][]string{
			"agency.txt": []string{
				"3",
				"1"},
			"calendar.txt": []string{
				"BUS319-1-Wdy-02",
				"BUS319-2-Wdy-02",
				"BUS319-3-Sa-02",
				"BUS319-4-Su-02",
				"BUS319-5-Wdy-02",
				"BUS319-6-Sa-02",
				"BUS319-7-Su-02",
				"BUS319-8-Wdy-02",
				"BUS319-9-Sa-02",
				"BUS319-A-Su-02"},
			"routes.txt": []string{
				"Red",
				"Mattapan",
				"Orange",
				"Green-B",
				"Green-C",
				"Green-D",
				"Green-E",
				"Blue",
				"741",
				"742"},
			"shapes.txt": []string{
				"010058",
				"010070",
				"040033",
				"040034",
				"040037",
				"040038",
				"050036",
				"050037",
				"070070",
				"070071"},
			"stops.txt": []string{
				"1",
				"10",
				"10000",
				"10003",
				"10005",
				"10006",
				"10007",
				"10008",
				"10009",
				"10010"},
			"trips.txt": []string{
				"40667713",
				"40667714",
				"40667720",
				"40667721",
				"40667723",
				"40667724",
				"40667728",
				"40667730",
				"40667731",
				"40667734"}}},
}

// Regenerate all external feeds feeds:
// NewReaderTesterFromReader for debugging and creating new tests.
// func NewReaderTesterFromReader(reader gotransit.Reader) ReaderTester {
// 	fe := ReaderTester{}
// 	fe.Counts = map[string]int{}
// 	fe.EntityIDs = map[string][]string{}
// 	add := func(ent gotransit.Entity) {
// 		fn, eid := ent.Filename(), ent.EntityID()
// 		fe.Counts[fn]++
// 		if eid != "" && len(fe.EntityIDs[fn]) < 10 {
// 			fe.EntityIDs[fn] = append(fe.EntityIDs[fn], eid)
// 		}
// 	}
// 	AllEntities(reader, add)
// 	return fe
// }
// func RegenerateExternalFeeds(t *testing.T) {
// 	dir := "../testdata/external"
// 	fis, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fes := map[string]testutil.ReaderTester{}
// 	for _, fi := range fis {
// 		fn := path.Join(dir, fi.Name())
// 		r, err := NewReader(fn)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fe := testutil.NewReaderTesterFromReader(r)
// 		fe.URL = fn
// 		fes[fi.Name()] = fe
// 	}
// 	fmt.Printf("%#v\n", fes)
// }
