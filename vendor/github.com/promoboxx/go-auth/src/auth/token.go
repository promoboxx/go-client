package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Roles
const (
	RoleSystem       = "system"
	RoleAdmin        = "admin"
	RoleUser         = "user"
	RoleAccounting   = "accounting"
	RoleApi          = "api"
	RoleInternal     = "internal"
	RoleDeploy       = "deploy"
	RoleClientServices = "client_services"

	// A uuid generated to determin a userID for system jwts
	SystemUUID = "00000000-0000-0000-0000-000000000000"
)

// NewRSAToken will return a token that logs using the provided logger
func NewRSAToken(log Logger, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Token { //TODO - figure out tracing paradigm and add here.
	if log == nil {
		log = new(nullLogger)
	}
	return &rsaToken{log: log, publicKey: publicKey, privateKey: privateKey}
}

//go:generate mockgen -destination=../authmock/token-mock.go -package=authmock github.com/promoboxx/go-auth/src/auth Token

// Token can generate and validate JWTs
type Token interface {
	GenerateJWT(issuer string, userID int64, userUUID string, roles []string, permissions Permission, duration time.Duration) (string, error)
	GenerateSystemToken(issuer string, initiatingUserID int64, duration time.Duration) (string, error)
	ValidateJWT(jwt string) (Claim, error)
}

type rsaToken struct {
	log        Logger
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (r *rsaToken) generate(issuer string, userID int64, userUUID string, initiatingUserID int64, initiatingUserUUID string, roles []string, permissions Permission, duration time.Duration) (string, error) {
	t := jwt.New(jwt.SigningMethodRS512)
	claims := JWTClaim{}
	claims.Subject = userID
	claims.SubjectUUID = userUUID
	claims.Issuer = issuer
	claims.ExpiresAt = time.Now().Add(duration).Unix()
	claims.IssuedAt = time.Now().Unix()
	claims.Permissions = permissions
	claims.Roles = roles
	claims.InitiatingUser = initiatingUserID
	claims.InitiatingUserUUID = initiatingUserUUID
	t.Claims = claims

	return t.SignedString(r.privateKey)
}

func (r *rsaToken) GenerateJWT(issuer string, userID int64, userUUID string, roles []string, permissions Permission, duration time.Duration) (string, error) {
	return r.generate(issuer, userID, userUUID, -1, SystemUUID, roles, permissions, duration)
}

func (r *rsaToken) GenerateSystemToken(issuer string, initiatingUserID int64, duration time.Duration) (string, error) {
	return r.generate(issuer, -1, SystemUUID, initiatingUserID, SystemUUID, []string{RoleSystem}, Permission{}, duration)
}

func (r *rsaToken) ValidateJWT(token string) (Claim, error) {
	result := Claim{}
	claims := JWTClaim{}
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		// validate the alg is what is RSA
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return r.publicKey, nil
	})
	if err != nil {
		return result, err
	}

	// make sure the token is valid
	if !t.Valid {
		return result, fmt.Errorf("Token not valid")
	}

	// generate a Claim object and return
	result = Claim{
		roles:              claims.Roles,
		expiration:         claims.GetExpiration(),
		userID:             claims.Subject,
		userUUID:           claims.SubjectUUID,
		brands:             claims.Permissions.Brands,
		brandsUUIDs:        claims.Permissions.BrandsUUIDs,
		divisions:          claims.Permissions.Divisions,
		retailerAccounts:   claims.Permissions.RetailerAccounts,
		retailers:          claims.Permissions.Retailers,
		intiatingUser:      claims.InitiatingUser,
		initiatingUserUUID: claims.InitiatingUserUUID,
		businessIDs:        claims.Permissions.BusinessIDs,
	}
	return result, nil
}
