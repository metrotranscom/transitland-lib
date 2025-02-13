package tl

import (
	"github.com/interline-io/transitland-lib/tl/causes"
	"github.com/interline-io/transitland-lib/tl/tt"
)

// Transfer transfers.txt
type Transfer struct {
	FromStopID      string `csv:",required" required:"true"`
	ToStopID        string `csv:",required" required:"true"`
	TransferType    int
	MinTransferTime Int
	BaseEntity
}

// Errors for this Entity.
func (ent *Transfer) Errors() (errs []error) {
	// transfer_type is required but can also be empty, so hard to distinguish
	errs = append(errs, ent.BaseEntity.Errors()...)
	errs = append(errs, tt.CheckPresent("from_stop_id", ent.FromStopID)...)
	errs = append(errs, tt.CheckPresent("to_stop_id", ent.ToStopID)...)
	errs = append(errs, tt.CheckInsideRangeInt("transfer_type", ent.TransferType, 0, 3)...)
	errs = append(errs, tt.CheckPositiveInt("min_transfer_time", ent.MinTransferTime.Val)...)
	return errs
}

// Filename transfers.txt
func (ent *Transfer) Filename() string {
	return "transfers.txt"
}

// TableName gtfs_transfers
func (ent *Transfer) TableName() string {
	return "gtfs_transfers"
}

// UpdateKeys updates entity references.
func (ent *Transfer) UpdateKeys(emap *EntityMap) error {
	// Adjust StopIDs
	if fromStopID, ok := emap.GetEntity(&Stop{StopID: ent.FromStopID}); ok {
		ent.FromStopID = fromStopID
	} else {
		return causes.NewInvalidReferenceError("from_stop_id", ent.FromStopID)
	}
	if toStopID, ok := emap.GetEntity(&Stop{StopID: ent.ToStopID}); ok {
		ent.ToStopID = toStopID
	} else {
		return causes.NewInvalidReferenceError("to_stop_id", ent.ToStopID)
	}
	return nil
}
