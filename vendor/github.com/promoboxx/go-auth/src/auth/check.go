package auth

//Check can see if a claim passes it
type Check interface {
	Pass(c Claim) bool
}

// CheckUser will check user IDs
type CheckUser struct {
	UserID int64
}

//Pass returns true if the claim's role is user and the userIDs match
func (u CheckUser) Pass(c Claim) bool {
	return c.IsRole(RoleUser) && c.IsUser(u.UserID)
}

// CheckUserUUID will check the users uuid
type CheckUserUUID struct {
	UserUUID string
}

// Pass returns true if the claims role is user and userUUIDs match
func (u CheckUserUUID) Pass(c Claim) bool {
	return c.IsRole(RoleUser) && c.IsUserByUUID(u.UserUUID)
}

// CheckRole will check role
type CheckRole struct {
	Role string
}

//Pass returns true if the claim's role matches
func (r CheckRole) Pass(c Claim) bool {
	return c.IsRole(r.Role)
}

// CheckRetailer will check retailer access
type CheckRetailer struct {
	RetailerID int64
}

//Pass returns true if the claim' has access to the retailer
func (r CheckRetailer) Pass(c Claim) bool {
	for _, v := range c.GetRetailers() {
		if v == r.RetailerID {
			return true
		}
	}
	return false
}

// CheckRetailerAccount will check retailer account access
type CheckRetailerAccount struct {
	RetailerAccountID int64
}

//Pass returns true if the claim's has access to retailer account
func (r CheckRetailerAccount) Pass(c Claim) bool {
	for _, v := range c.GetRetailerAccounts() {
		if v == r.RetailerAccountID {
			return true
		}
	}
	return false
}

// CheckBrand will check brand access
type CheckBrand struct {
	BrandID int64
}

//Pass returns true if the claim's role matches
func (b CheckBrand) Pass(c Claim) bool {
	for _, v := range c.GetBrands() {
		if v == b.BrandID {
			return true
		}
	}
	return false
}

// CheckBrandUUID checks for brand access
type CheckBrandUUID struct {
	BrandUUID string
}

// Pass returns true if the brand matches
func (b CheckBrandUUID) Pass(c Claim) bool {
	for _, v := range c.GetBrandsUUIDs() {
		if v == b.BrandUUID {
			return true
		}
	}

	return false
}

// CheckDivision will check division access
type CheckDivision struct {
	DivisionID int64
}

//Pass returns true if the claim has division accress
func (d CheckDivision) Pass(c Claim) bool {
	for _, v := range c.GetDivisions() {
		if v == d.DivisionID {
			return true
		}
	}
	return false
}

type CheckBusinessID struct {
	BusinessID string
}

// Pass returns true if the claim has business access
func (b CheckBusinessID) Pass(c Claim) bool {
	for _, v := range c.GetBusinessIDs() {
		if v == b.BusinessID {
			return true
		}
	}

	return false
}
