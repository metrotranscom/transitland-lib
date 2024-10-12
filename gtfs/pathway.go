package gtfs

import (
	"errors"
	"fmt"

	"github.com/interline-io/transitland-lib/causes"
	"github.com/interline-io/transitland-lib/tt"
)

// Pathway pathways.txt
type Pathway struct {
	PathwayID           tt.String `csv:",required"`
	FromStopID          tt.String `csv:",required"`
	ToStopID            tt.String `csv:",required"`
	PathwayMode         tt.Int    `csv:",required"`
	IsBidirectional     tt.Int    `csv:",required"`
	Length              tt.Float  `csv:"length" min:"0"`
	TraversalTime       tt.Int    `csv:"traversal_time" min:"0"`
	StairCount          tt.Int    `csv:"stair_count"`
	MaxSlope            tt.Float  `csv:"max_slope"`
	MinWidth            tt.Float  `csv:"min_width"`
	SignpostedAs        tt.String `csv:"signposted_as"`
	ReverseSignpostedAs tt.String `csv:"reversed_signposted_as"`
	tt.BaseEntity
}

// EntityID returns the ID or StopID.
func (ent *Pathway) EntityID() string {
	return entID(ent.ID, ent.PathwayID.Val)
}

// EntityKey returns the GTFS identifier.
func (ent *Pathway) EntityKey() string {
	return ent.PathwayID.Val
}

// Filename pathways.txt
func (ent *Pathway) Filename() string {
	return "pathways.txt"
}

// TableName ext_pathway_pathways
func (ent *Pathway) TableName() string {
	return "gtfs_pathways"
}

// UpdateKeys updates Entity references.
func (ent *Pathway) UpdateKeys(emap *EntityMap) error {
	if fkid, ok := emap.GetEntity(&Stop{StopID: ent.FromStopID.Val}); ok {
		ent.FromStopID.Set(fkid)
	} else {
		return causes.NewInvalidReferenceError("from_stop_id", ent.FromStopID.Val)
	}
	if fkid, ok := emap.GetEntity(&Stop{StopID: ent.ToStopID.Val}); ok {
		ent.ToStopID.Set(fkid)
	} else {
		return causes.NewInvalidReferenceError("to_stop_id", ent.ToStopID.Val)
	}
	return nil
}

// GetString returns the string representation of an field.
func (ent *Pathway) GetString(key string) (string, error) {
	v := ""
	switch key {
	case "pathway_id":
		v = ent.PathwayID.String()
	case "from_stop_id":
		v = ent.FromStopID.String()
	case "to_stop_id":
		v = ent.ToStopID.String()
	case "pathway_mode":
		v = ent.PathwayMode.String()
	case "is_bidirectional":
		v = ent.IsBidirectional.String()
	case "length":
		if ent.Length.Val > 0 {
			v = fmt.Sprintf("%0.5f", ent.Length.Val)
		}
	case "traversal_time":
		if ent.TraversalTime.Val > 0 {
			v = ent.TraversalTime.String()
		}
	case "stair_count":
		if ent.StairCount.Val != 0 && ent.StairCount.Val != -1 {
			v = ent.StairCount.String()
		}
	case "max_slope":
		if ent.MaxSlope.Val != 0 {
			v = fmt.Sprintf("%0.2f", ent.MaxSlope.Val)
		}
	case "min_width":
		if ent.MinWidth.Val != 0 {
			v = fmt.Sprintf("%0.2f", ent.MinWidth.Val)
		}
	case "signposted_as":
		v = ent.SignpostedAs.String()
	case "reversed_signposted_as":
		v = ent.ReverseSignpostedAs.String()
	default:
		return v, errors.New("unknown key")
	}
	return v, nil
}
