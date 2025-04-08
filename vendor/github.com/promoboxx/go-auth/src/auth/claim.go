package auth

import "time"

//Claim represents the authorizations of a user
type Claim struct {
	brands             []int64
	brandsUUIDs        []string
	divisions          []int64
	expiration         time.Time
	retailerAccounts   []int64
	retailers          []int64
	roles              []string
	userID             int64
	userUUID           string
	intiatingUser      int64
	initiatingUserUUID string
	businessIDs        []string
}

// NewClaim creates a new claim
func NewClaim(brands []int64, brandsUUIDs []string, divisions []int64, expiration time.Time, retailerAccounts, retailers []int64, roles []string, userID int64, userUUID string, initiatingUser int64, initiatingUserUUID string, businessIDs []string) Claim {
	return Claim{brands: brands, brandsUUIDs: brandsUUIDs, divisions: divisions, expiration: expiration, retailerAccounts: retailerAccounts, retailers: retailers, roles: roles, userID: userID, userUUID: userUUID, intiatingUser: initiatingUser, initiatingUserUUID: initiatingUserUUID, businessIDs: businessIDs}
}

// HasPermission will return true if any checks pass
func (c Claim) HasPermission(checks ...Check) bool {
	for _, v := range checks {
		if v.Pass(c) {
			return true
		}
	}
	return false
}

// HasAllPermissions will return true if all checks pass
func (c Claim) HasAllPermissions(checks ...Check) bool {
	if len(checks) == 0 {
		return false
	}

	for _, v := range checks {
		if !v.Pass(c) {
			return false
		}
	}
	return true
}

//GetBrands returns the set of brands this user has access to
func (c Claim) GetBrands() []int64 {
	return c.brands
}

// GetBrandsUUIDs returns the uuids of the brands the user has access to
func (c Claim) GetBrandsUUIDs() []string {
	return c.brandsUUIDs
}

// GetBusinessIDs will provide a group of retail businesses
func (c Claim) GetBusinessIDs() []string {
	return c.businessIDs
}

//GetDivisions returns the set of divisions this user has access to
func (c Claim) GetDivisions() []int64 {
	return c.divisions
}

//GetExpiration returns the time this claim expires
func (c Claim) GetExpiration() time.Time {
	return c.expiration
}

// GetInitiatingUser returns the initiating user id if any
func (c Claim) GetInitiatingUser() int64 {
	return c.intiatingUser
}

// GetInitiatingUserUUID returns the initiating user uuid if any
func (c Claim) GetInitiatingUserUUID() string {
	return c.initiatingUserUUID
}

//GetRetailerAccounts returns the set of retailer accounts this user has access to
func (c Claim) GetRetailerAccounts() []int64 {
	return c.retailerAccounts
}

//GetRetailers returns the set of retailers this user has access to
func (c Claim) GetRetailers() []int64 {
	return c.retailers
}

//GetRoles returns the user's roles
func (c Claim) GetRoles() []string {
	return c.roles
}

//GetUserID returns the user's ID
func (c Claim) GetUserID() int64 {
	return c.userID
}

// GetUserUUID returns the uuid of the user
func (c Claim) GetUserUUID() string {
	return c.userUUID
}

//IsRole returns true if the user is that role
func (c Claim) IsRole(role string) bool {
	for _, v := range c.roles {
		if v == role {
			return true
		}
	}
	return false
}

//IsAdmin returns true if the user is an admin or system user
func (c Claim) IsAdmin() bool {
	return c.IsRole(RoleAdmin) || c.IsRole(RoleSystem)
}

//IsInternal returns true if the user has the internal role or is a client_services user
func (c Claim) IsInternal() bool {
	return c.IsRole(RoleInternal) || c.IsClientServices()
}

//IsSystem returns true if the user is a system user
func (c Claim) IsSystem() bool {
	return c.IsRole(RoleSystem)
}

//IsApi returns true if the user has the api role
func (c Claim) IsApi() bool {
	return c.IsRole(RoleApi)
}

//IsAccounting returns true if the user has the accounting role
func (c Claim) IsAccounting() bool {
	return c.IsRole(RoleAccounting)
}

//IsDeploy returns true if the user has the deploy role
func (c Claim) IsDeploy() bool {
	return c.IsRole(RoleDeploy)
}

//IsClientServices returns true if the user has the client_services role or is an admin
func (c Claim) IsClientServices() bool {
	return c.IsRole(RoleClientServices) || c.IsAdmin()
}

//IsUser returns true if the user matches the id given
func (c Claim) IsUser(userID int64) bool {
	return c.userID == userID
}

// IsUserByUUID returns true if the user id given matches the one in the claims
func (c Claim) IsUserByUUID(userUUID string) bool {
	return c.userUUID == userUUID
}
