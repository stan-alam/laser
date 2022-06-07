package laser

import "time"

type Token struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
}
