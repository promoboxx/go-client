package auth

import (
	"fmt"
	"time"
)

// Permission holds the set of things a claim has access to
type Permission struct {
	Retailers        []int64  `json:"retailers"`
	RetailerAccounts []int64  `json:"retailer_accounts"`
	Brands           []int64  `json:"brands"`
	BrandsUUIDs      []string `json:"brands_uuids"`
	Divisions        []int64  `json:"divisions"`
	BusinessIDs      []string `json:"business_ids"`
}

// JWTClaim is a jwt.Claim
type JWTClaim struct {
	Audience           string     `json:"aud,omitempty"`
	ExpiresAt          int64      `json:"exp,omitempty"`
	ID                 string     `json:"jti,omitempty"`
	IssuedAt           int64      `json:"iat,omitempty"`
	Issuer             string     `json:"iss,omitempty"`
	NotBefore          int64      `json:"nbf,omitempty"`
	Subject            int64      `json:"sub,omitempty"`
	SubjectUUID        string     `json:"sub_uuid,omitempty"`
	Permissions        Permission `json:"permission,omitempty"`
	Roles              []string   `json:"roles"`
	InitiatingUser     int64      `json:"initiating_user,omitempty"`
	InitiatingUserUUID string     `json:"initiating_user_uuid,omitempty"`
}

// Valid fulfils the jwt.Claim interface
func (j JWTClaim) Valid() error {
	if j.GetExpiration().Before(time.Now()) {
		return fmt.Errorf("JWT is expired")
	}
	if len(j.Roles) == 0 {
		return fmt.Errorf("Role not set")
	}
	if j.Subject == 0 {
		return fmt.Errorf("Sub not set")
	}
	if len(j.SubjectUUID) == 0 {
		return fmt.Errorf("SubjectUUID not set")
	}
	return nil
}

// GetExpiration converts ExpiresAt to a time.Time
func (j JWTClaim) GetExpiration() time.Time {
	return time.Unix(j.ExpiresAt, 0)
}
