//go:build integration
// +build integration

package entrypoints

// This is an integration test because the token will show up as expired, you will need a live token

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/service"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	endpoint := "localhost:8088"

	var opts []grpc.DialOption

	// May needs absolute path to pem file
	creds, err := credentials.NewClientTLSFromFile("$HOME/flyte/mb-dev-utils/certs/server.pem", "example.com")
	assert.NoError(t, err)
	opts = append(opts, grpc.WithTransportCredentials(creds))

	token := oauth2.Token{
		AccessToken: GetOAuthTokenByJSON(getToken()),
	}
	tokenRpcCredentials := oauth.NewOauthAccess(&token)
	tokenDialOption := grpc.WithPerRPCCredentials(tokenRpcCredentials)
	opts = append(opts, tokenDialOption)

	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		fmt.Printf("Dial error %v\n", err)
	}
	assert.NoError(t, err)
	client := service.NewAdminServiceClient(conn)
	resp, err := client.ListProjects(ctx, &admin.ProjectListRequest{})
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	assert.NoError(t, err)
	fmt.Printf("Response: %v\n", resp)
}

func getToken() string {
	return `
	{
		"iss": "https://ssoalpha.dvb.corpinter.net/v1",
		"sub": "MYUSERA",
		"aud": [
			"sample_aud"
		],
		"exp": 1645544098,
		"iat": 1645540498,
		"azp": "sample_aud",
		"at_hash": "CsJQNnJ_h7ren1dkw1_qEw",
		"email": "user@mercedes-benz.com",
		"email_verified": true,
		"scope": [
			"groups",
			"offline_access",
			"openid",
			"email",
			"profile"
		],
		"name": "My UserA",
		"preferred_username": "My UserA",
		"user_type": "internal",
		"aad_id": "user@mercedes-benz.com",
		"family_name": "UserA",
		"given_name": "My"
	}
	`
}

const tokenPayload = `{"sub":"1234567890","name":"John Doe","iat":1516239022}`

// GetOAuthTokenByJSON expects a json formatted string and returns it as JWT token
func GetOAuthTokenByJSON(payload string) string {
	key := []byte("secret")
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	var body jwtPayload
	if err = json.Unmarshal([]byte(payload), &body); err != nil {
		panic(err)
	}

	raw, err := jwt.Signed(sig).Claims(body).CompactSerialize()
	if err != nil {
		panic(err)
	}

	return raw
}

type jwtPayload struct {
	ISS               string   `json:"iss,omitempty"`
	SUB               string   `json:"sub,omitempty"`
	AUD               []string `json:"aud,omitempty"`
	EXP               int      `json:"exp,omitempty"`
	IAT               int      `json:"iat,omitempty"`
	AZP               string   `json:"azp,omitempty"`
	AtHash            string   `json:"at_hash,omitempty"`
	Email             string   `json:"email,omitempty"`
	EmailVerfified    bool     `json:"email_verified,omitempty"`
	Scope             []string `json:"scpoe,omitempty"`
	Name              string   `json:"name,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	UserType          string   `json:"user_type,omitempty"`
	AADID             string   `json:"aad_id,omitempty"`
	FamilyName        string   `json:"family_name,omitempty"`
	GivenName         string   `json:"given_name,omitempty"`
}
