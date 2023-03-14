package oakoidc

import (
	"encoding/json"
	"fmt"

	"github.com/coreos/go-oidc"
)

// StandardClaims captures fields that are typically included in [oidc.UserInfo] claims by various OIDC providers.
//
// https://developers.google.com/identity/openid-connect/openid-connect#an-id-tokens-payload
type StandardClaims struct {
	Email            string `json:"email"`
	EmailVerified    bool   `json:"email_verified"`
	FamilyName       string `json:"family_name"`
	GivenName        string `json:"given_name"`
	OrganizationName string `json:"hd"`
	Locale           string `json:"locale"`
	Name             string `json:"name"`
	Picture          string `json:"picture"`
	AccountID        string `json:"sub"`
}

func (s *StandardClaims) String() string {
	result, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return "<invalid standard claims>"
	}
	return string(result)
}

func NewStandardClaims(info *oidc.UserInfo) (*StandardClaims, error) {
	var claims *StandardClaims
	if err := info.Claims(&claims); err != nil {
		return nil, fmt.Errorf("cannot decode token claims: %w", err)
	}
	return claims, nil
}
