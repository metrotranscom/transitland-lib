package tt

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Key is a nullable foreign key constraint, similar to sql.NullString
type Key struct {
	Val   string
	Valid bool
}

func NewKey(v string) Key {
	return Key{Valid: true, Val: v}
}

func (r *Key) String() string {
	return r.Val
}

func (r Key) Value() (driver.Value, error) {
	if !r.Valid || r.Val == "" {
		return nil, nil
	}
	return r.Val, nil
}

func (r *Key) Scan(src interface{}) error {
	r.Val, r.Valid = "", false
	var err error
	switch v := src.(type) {
	case nil:
		return nil
	case string:
		r.Val = v
	case int:
		r.Val = strconv.Itoa(v)
	case int64:
		r.Val = strconv.Itoa(int(v))
	default:
		err = fmt.Errorf("cant convert %T to Key", src)
	}
	r.Valid = (err == nil && r.Val != "")
	return err
}

func (r *Key) Int() int {
	a, _ := strconv.Atoi(r.Val)
	return a
}

func (r *Key) UnmarshalJSON(v []byte) error {
	r.Val, r.Valid = "", false
	if isEmpty(string(v)) {
		return nil
	}
	if v[0] != '"' && v[len(v)-1] != '"' {
		// Handle unquoted values, e.g. number
		return r.Scan(string(v))
	}
	err := json.Unmarshal(v, &r.Val)
	r.Valid = (err == nil && r.Val != "")
	return err
}

func (r Key) MarshalJSON() ([]byte, error) {
	if !r.Valid {
		return jsonNull(), nil
	}
	return json.Marshal(r.Val)
}

func (r *Key) UnmarshalGQL(v interface{}) error {
	return r.Scan(v)
}

func (r Key) MarshalGQL(w io.Writer) {
	b, _ := r.MarshalJSON()
	w.Write(b)
}
