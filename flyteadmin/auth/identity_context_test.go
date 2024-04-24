package auth

import (
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/service"
	"github.com/flyteorg/flyte/flyteplugins/go/tasks/pluginmachinery/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/sets"
)

func TestGetClaims(t *testing.T) {
	noClaims := map[string]interface{}(nil)
	noClaimsCtx, err := NewIdentityContext("", "", "", time.Now(), nil, nil, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, noClaims, noClaimsCtx.Claims())

	claims := map[string]interface{}{
		"groups":    []string{"g1", "g2"},
		"something": "else",
	}

	userInfo := &service.UserInfoResponse{}
	userInfoClaims := map[string]interface{}{
		"roles": []string{"r1", "r2"},
	}
	userInfo.AdditionalClaims, err = utils.MarshalObjToStruct(&userInfoClaims)
	assert.NoError(t, err)

	withClaimsCtx, err := NewIdentityContext("", "", "", time.Now(), nil, userInfo, claims)
	assert.NoError(t, err)
	assert.EqualValues(t, claims, withClaimsCtx.Claims())
	assert.NotEmpty(t, withClaimsCtx.UserInfo().AdditionalClaims)
}

func TestWithExecutionUserIdentifier(t *testing.T) {
	idctx, err := NewIdentityContext("", "", "", time.Now(), sets.String{}, nil, nil)
	assert.NoError(t, err)
	newIDCtx := idctx.WithExecutionUserIdentifier("byhsu")
	// make sure the original one is intact
	assert.Equal(t, "", idctx.ExecutionIdentity())
	assert.Equal(t, "byhsu", newIDCtx.ExecutionIdentity())
}
