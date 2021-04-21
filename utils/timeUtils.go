package utils

import (
	"encoding/json"
	"time"
)

type Long int64

func (l Long) ToTime() time.Time {
	return time.Unix(int64(l/1000), int64(1e6*(l%1000)))
}

func (l Long) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(l))
}
func (l *Long) UnmarshalJSON(b []byte) error {
	var i int64
	err := json.Unmarshal(b, &i)
	if err == nil {
		*l = Long(i)
	}
	return err
}
