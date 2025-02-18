package plus

import (
	"github.com/interline-io/transitland-lib/causes"
	"github.com/interline-io/transitland-lib/gtfs"
	"github.com/interline-io/transitland-lib/tt"
)

// RealtimeStop realtime_stops.txt
type RealtimeStop struct {
	TripID         string `csv:"trip_id"`
	StopID         string `csv:"stop_id"`
	RealtimeStopID string `csv:"realtime_stop_id"`
	tt.BaseEntity
}

// Filename realtime_stops.txt
func (ent *RealtimeStop) Filename() string {
	return "realtime_stops.txt"
}

// TableName ext_plus_realtime_stops
func (ent *RealtimeStop) TableName() string {
	return "ext_plus_realtime_stops"
}

// UpdateKeys updates Entity references.
func (ent *RealtimeStop) UpdateKeys(emap *tt.EntityMap) error {
	if fkid, ok := emap.GetEntity(&gtfs.Trip{TripID: tt.NewString(ent.TripID)}); ok {
		ent.TripID = fkid
	} else {
		return causes.NewInvalidReferenceError("trip_id", ent.TripID)
	}
	if fkid, ok := emap.GetEntity(&gtfs.Stop{StopID: tt.NewString(ent.StopID)}); ok {
		ent.StopID = fkid
	} else {
		return causes.NewInvalidReferenceError("stop_id", ent.StopID)
	}
	// if fkid, ok := emap.Get(&gtfs.Stop{StopID: ent.RealtimeStopID}); ok {
	// 	ent.RealtimeStopID = fkid
	// } else {
	// 	return causes.NewInvalidReferenceError("stop_id", ent.StopID)
	// }
	return nil
}
