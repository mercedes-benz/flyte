package authzserver

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/flyteorg/flyte/flyteadmin/auth/config"
	"github.com/flyteorg/flyte/flyteadmin/auth/interfaces/mocks"
)

func Test_tokenEndpoint(t *testing.T) {
	t.Run("Empty request", func(t *testing.T) {
		body := &bytes.Buffer{}
		req := httptest.NewRequest(http.MethodPost, "/token", body)
		rw := httptest.NewRecorder()
		oauth2Provider, _ := newMockProvider(t)
		mockAuthCtx := &mocks.AuthenticationContext{}
		mockAuthCtx.OnOAuth2Provider().Return(oauth2Provider)
		mockAuthCtx.OnOptions().Return(&config.Config{})

		tokenEndpoint(mockAuthCtx, rw, req)
		assert.NotEmpty(t, rw.Body.String())
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("Valid token request authorization code", func(t *testing.T) {
		// create a signer for rsa 256
		tok := jwtgo.New(jwtgo.GetSigningMethod("RS256"))

		oauth2Provider, secrets := newMockProvider(t)

		//tok.Header[KeyIDClaim] = secrets.TokenSigningRSAPrivateKey.KeyID()
		tok.Claims = &CustomClaimsExample{
			StandardClaims: &jwtgo.StandardClaims{},
			ClientID:       "flytectl",
			Scopes:         []string{"access_token", "offline"},
		}

		// Create token string
		authCode, err := tok.SignedString(secrets.TokenSigningRSAPrivateKey)
		assert.NoError(t, err)

		payload := url.Values{
			"code":       {authCode},
			"grant_type": {"authorization_code"},
			"scope":      {"all", "offline"},
		}

		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader([]byte(payload.Encode())))
		req.PostForm = payload
		req.Header.Set("authorization", basicAuth("flytectl", "foobar"))
		mockAuthCtx := &mocks.AuthenticationContext{}
		mockAuthCtx.OnOAuth2Provider().Return(oauth2Provider)
		mockAuthCtx.OnOptions().Return(&config.Config{})

		rw := httptest.NewRecorder()
		tokenEndpoint(mockAuthCtx, rw, req)
		if !assert.Equal(t, http.StatusOK, rw.Code) {
			t.FailNow()
		}

		m := map[string]interface{}{}
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &m))
		assert.True(t, m["expires_in"].(float64) <= 0.5*time.Hour.Seconds())

		assert.NotEmpty(t, m["access_token"])
		// Parse and validate the token.
		parsedToken, err := jwtgo.Parse(m["access_token"].(string), func(tok *jwtgo.Token) (interface{}, error) {
			keySet, err := newJSONWebKeySet([]rsa.PublicKey{secrets.TokenSigningRSAPrivateKey.PublicKey})
			assert.NoError(t, err)

			return findPublicKeyForTokenOrFirst(context.Background(), tok, keySet)
		})

		assert.NoError(t, err)

		expectedAccessTokenExpiry := time.Now().Add(config.DefaultConfig.AppAuth.SelfAuthServer.AccessTokenLifespan.Duration - time.Second).Unix()
		expectedRefreshTokenExpiry := time.Now().Add(config.DefaultConfig.AppAuth.SelfAuthServer.RefreshTokenLifespan.Duration - time.Second).Unix()

		claims := parsedToken.Claims.(jwtgo.MapClaims)
		assert.True(t, claims.VerifyExpiresAt(expectedAccessTokenExpiry, true))

		assert.NotEmpty(t, m["refresh_token"])
		// Parse and validate the token.
		parsedToken, err = jwtgo.Parse(m["refresh_token"].(string), func(tok *jwtgo.Token) (interface{}, error) {
			keySet, err := newJSONWebKeySet([]rsa.PublicKey{secrets.TokenSigningRSAPrivateKey.PublicKey})
			assert.NoError(t, err)

			return findPublicKeyForTokenOrFirst(context.Background(), tok, keySet)
		})

		assert.NoError(t, err)

		claims = parsedToken.Claims.(jwtgo.MapClaims)
		assert.True(t, claims.VerifyExpiresAt(expectedRefreshTokenExpiry, true))
	})

	t.Run("Valid token request client credentials", func(t *testing.T) {
		oauth2Provider, secrets := newMockProvider(t)

		payload := url.Values{
			"grant_type": {"client_credentials"},
			"scope":      {"all", "offline"},
		}

		req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader([]byte(payload.Encode())))
		req.PostForm = payload
		req.Header.Set("authorization", basicAuth("flytepropeller", "foobar"))
		mockAuthCtx := &mocks.AuthenticationContext{}
		mockAuthCtx.OnOAuth2Provider().Return(oauth2Provider)
		mockAuthCtx.OnOptions().Return(&config.Config{})

		rw := httptest.NewRecorder()
		tokenEndpoint(mockAuthCtx, rw, req)
		if !assert.Equal(t, http.StatusOK, rw.Code) {
			t.FailNow()
		}

		m := map[string]interface{}{}
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &m))
		assert.True(t, m["expires_in"].(float64) <= 0.5*time.Hour.Seconds())

		assert.NotEmpty(t, m["access_token"])
		// Parse and validate the token.
		parsedToken, err := jwtgo.Parse(m["access_token"].(string), func(tok *jwtgo.Token) (interface{}, error) {
			keySet, err := newJSONWebKeySet([]rsa.PublicKey{secrets.TokenSigningRSAPrivateKey.PublicKey})
			assert.NoError(t, err)

			return findPublicKeyForTokenOrFirst(context.Background(), tok, keySet)
		})

		assert.NoError(t, err)
		expectedAccessTokenExpiry := time.Now().Add(config.DefaultConfig.AppAuth.SelfAuthServer.AccessTokenLifespan.Duration - time.Second).Unix()

		claims := parsedToken.Claims.(jwtgo.MapClaims)

		assert.Equal(t, "flytepropeller", claims["client_id"].(string))
		assert.True(t, claims.VerifyExpiresAt(expectedAccessTokenExpiry, true))
		assert.NoError(t, err)

	})

}

func basicAuth(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
}
