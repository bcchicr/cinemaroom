package resources

import "time"

type RefreshTokenResource struct {
	Value     string
	ExpiresAt time.Time
}
