package unversioned

import (
	"database/sql/driver"
	"time"
)

const (
	nullTime      = "0000-00-00 00:00:00"
	sqlTimeLayout = "2030-12-01 22:01:23 -0700 UTC"
)

// implement sql.Scanner.
func (t *Time) Scan(value interface{}) error {
	sqlTime, ok := value.(time.Time)
	if ok {
		t.Time = sqlTime
	}
	return nil
}

//  sql/driver.Value implementation to go from unversioned.time ->time.time.
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}
