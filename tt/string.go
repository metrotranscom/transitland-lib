package tt

import "strconv"

type String struct {
	Option[string]
}

func (r String) Int() int {
	if !r.Valid {
		return 0
	}
	a, _ := strconv.ParseInt(r.Val, 10, 64)
	return int(a)
}

func NewString(v string) String {
	return String{Option: NewOption(v)}
}
