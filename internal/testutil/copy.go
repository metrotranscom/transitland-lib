package testutil

import (
	"fmt"

	"github.com/interline-io/gotransit"
)

type canCreateFV interface {
	CreateFeedVersion(reader gotransit.Reader) (int, error)
}

// DirectCopy does a direct reader->writer copy, with minimal validation and changes.
func DirectCopy(reader gotransit.Reader, writer gotransit.Writer) error {
	emap := gotransit.NewEntityMap()
	cp := func(ent gotransit.Entity) error {
		sid := ent.EntityID()
		if err := ent.UpdateKeys(emap); err != nil {
			return fmt.Errorf("entity: %#v error: %s", ent, err)
		}
		eid, err := writer.AddEntity(ent)
		if err != nil {
			return fmt.Errorf("entity: %#v error: %s", ent, err)
		}
		return emap.Set(ent, sid, eid)
	}
	// Create any FV
	if w2, ok := writer.(canCreateFV); ok {
		w2.CreateFeedVersion(reader)
	}
	for ent := range reader.Agencies() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.Routes() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	// Two passes required
	for ent := range reader.Stops() {
		if ent.LocationType != 1 {
			continue
		}
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.Stops() {
		if ent.LocationType == 1 {
			continue
		}
		if err := cp(&ent); err != nil {
			return err
		}
	}
	// Done with stops
	for ent := range reader.Calendars() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.CalendarDates() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.Shapes() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.Trips() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ents := range reader.StopTimesByTripID() {
		e2s := []gotransit.Entity{}
		for i := 0; i < len(ents); i++ {
			ents[i].UpdateKeys(emap)
			e2s = append(e2s, &ents[i])
		}
		if err := writer.AddEntities(e2s); err != nil {
			for _, ent := range e2s {
				fmt.Printf("%#v\n", ent)
			}
			return err
		}
	}
	for ent := range reader.Frequencies() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.Transfers() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.FareAttributes() {
		ent.Transfers = "0"
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.FareRules() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	for ent := range reader.FeedInfos() {
		if err := cp(&ent); err != nil {
			return err
		}
	}
	return nil
}
