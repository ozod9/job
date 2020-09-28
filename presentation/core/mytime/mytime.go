package mytime

import "time"

type MyTime struct {
	*time.Time
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format("\"2006-01-02 15:04:05\"")), nil
}

func (t *MyTime) UnmarshalJSON(data []byte) (err error) {
	tt, err := time.Parse("\"2006-01-02 15:04:05\"", string(data))
	*t = MyTime{&tt}
	return
}
