package auth

import (
	"testing"

	authConfig "github.com/flyteorg/flyte/flyteadmin/auth/config"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/service"
	"github.com/flyteorg/flyte/flytestdlib/utils"
	"github.com/stretchr/testify/assert"
)

func TestEntitlementsFromIndentityContext(t *testing.T) {

	t.Run("entitlements as list", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": []string{"r1", "r2"}}
		idCtx := IdentityContext{claims: &claims}
		entitlements, err := entitlementsFromClaim(idCtx)
		assert.NoError(t, err)
		assert.EqualValues(t, []string{"r1", "r2"}, entitlements)
	})

	t.Run("entitlement as string", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": "r2"}
		idCtx := IdentityContext{claims: &claims}
		entitlements, err := entitlementsFromClaim(idCtx)
		assert.NoError(t, err)
		assert.EqualValues(t, []string{"r2"}, entitlements)
	})

	t.Run("entitlement arbitrarily nested value isn't supported yet", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": map[string]interface{}{"entitlement_group": "r2"}}
		idCtx := IdentityContext{claims: &claims}
		entitlements, err := entitlementsFromClaim(idCtx)
		assert.ErrorContains(t, err, "Failed to convert claim entitlement_group")
		assert.Nil(t, entitlements)
	})

	t.Run("entitlements being part of additional claims", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": []interface{}{"r1", "r2"}}
		additionalClaimsPbStruct, _ := utils.MarshalObjToStruct(claims)
		idCtx := IdentityContext{
			userInfo: &service.UserInfoResponse{
				AdditionalClaims: additionalClaimsPbStruct,
			},
		}
		entitlements, err := entitlementsFromClaim(idCtx)
		assert.NoError(t, err)
		assert.EqualValues(t, []string{"r1", "r2"}, entitlements)
	})

	t.Run("entitlement arbitrarily nested value beneath AdditionalClaims isn't supported", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": []interface{}{map[string]interface{}{"entitlement_group": []interface{}{"r1", "r2"}}, "r2"}}
		additionalClaimsPbStruct, _ := utils.MarshalObjToStruct(claims)
		idCtx := IdentityContext{
			userInfo: &service.UserInfoResponse{
				AdditionalClaims: additionalClaimsPbStruct,
			},
		}
		entitlements, err := entitlementsFromClaim(idCtx)
		assert.ErrorContains(t, err, "Failed to convert claim entitlement_group")
		assert.Nil(t, entitlements)
	})
}

func TestProjectsByEntitlements(t *testing.T) {

	getAuthConfigFunc = func() *authConfig.Config {
		return &authConfig.Config{
			ProjectAuthorization: authConfig.ProjectAuthorizationConfig{
				ProjectSets: map[string][]string{
					"r1": {"p1", "p2"},
					"r2": {"p3"},
				},
			},
		}
	}

	t.Run("join entitlements and projects", func(t *testing.T) {
		entitlements := []string{"r1", "r2"}

		projects := projectsByEntitlements(entitlements)
		expected := map[string]any{
			"p1": nil,
			"p2": nil,
			"p3": nil,
		}
		assert.EqualValues(t, expected, projects)
	})

	t.Run("using nil for entitlements", func(t *testing.T) {
		assert.NotPanics(t, func() {
			projects := projectsByEntitlements(nil)
			expected := map[string]any{}
			assert.EqualValues(t, expected, projects)
		})
	})

	t.Run("empty ProjectSets", func(t *testing.T) {
		getAuthConfigFunc = func() *authConfig.Config {
			return &authConfig.Config{}
		}
		entitlements := []string{"r1", "r2"}

		assert.NotPanics(t, func() {
			projects := projectsByEntitlements(entitlements)
			expected := map[string]any{}
			assert.EqualValues(t, expected, projects)
		})
	})
}

func TestEligibleProjectsByIdentity(t *testing.T) {

	getAuthConfigFunc = func() *authConfig.Config {
		return &authConfig.Config{
			ProjectAuthorization: authConfig.ProjectAuthorizationConfig{
				ProjectSets: map[string][]string{
					"r1": {"p1", "p2"},
					"r2": {"p3"},
				},
				AppAuth: authConfig.ProjectAuthorizationAppAuth{
					Mappings: []authConfig.ProjectAuthorizationClientIDProjectSetMapping{
						{ClientID: "cleanerApp", ProjectSets: []string{"r1", "r2"}},
					},
				},
				UserAuth: authConfig.ProjectAuthorizationUserAuth{
					Claim: "entitlement_group",
				},
			},
		}
	}

	t.Run("App user", func(t *testing.T) {
		idCtx := IdentityContext{appID: "cleanerApp"}
		projects, err := eligibleProjectsByIdentity(idCtx)
		expected := map[string]any{
			"p1": nil,
			"p2": nil,
			"p3": nil,
		}
		assert.NoError(t, err)
		assert.EqualValues(t, expected, projects)
	})

	t.Run("Human user", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": []string{"r1", "r2"}}
		idCtx := IdentityContext{claims: &claims}
		projects, err := eligibleProjectsByIdentity(idCtx)
		expected := map[string]any{
			"p1": nil,
			"p2": nil,
			"p3": nil,
		}
		assert.NoError(t, err)
		assert.EqualValues(t, expected, projects)
	})
}

func TestFilterProjectsByEntitlement(t *testing.T) {

	getAuthConfigFunc = func() *authConfig.Config {
		return &authConfig.Config{
			ProjectAuthorization: authConfig.ProjectAuthorizationConfig{
				ProjectSets: map[string][]string{
					"r1": {"p1", "p2"},
					"r2": {"p3"},
				},
				AppAuth: authConfig.ProjectAuthorizationAppAuth{
					Mappings: []authConfig.ProjectAuthorizationClientIDProjectSetMapping{
						{ClientID: "cleanerApp", ProjectSets: []string{"r1", "r2"}},
					},
				},
				UserAuth: authConfig.ProjectAuthorizationUserAuth{
					Claim: "entitlement_group",
				},
			},
		}
	}

	t.Run("App user", func(t *testing.T) {
		idCtx := IdentityContext{appID: "cleanerApp"}
		projects := admin.Projects{
			Projects: []*admin.Project{
				{Id: "p1"},
				{Id: "p2"},
				{Id: "p3"},
				{Id: "p4"},
			},
		}
		err := FilterProjectsByEntitlement(idCtx, &projects)
		expected := admin.Projects{
			Projects: []*admin.Project{
				{Id: "p1"},
				{Id: "p2"},
				{Id: "p3"},
			},
		}
		assert.NoError(t, err)
		assert.EqualValues(t, expected.Projects, projects.Projects)
	})

	t.Run("Human user", func(t *testing.T) {
		claims := map[string]interface{}{"entitlement_group": []string{"r1", "r2"}}
		idCtx := IdentityContext{claims: &claims}
		projects := admin.Projects{
			Projects: []*admin.Project{
				{Id: "p1"},
				{Id: "p2"},
				{Id: "p3"},
				{Id: "p4"},
			},
		}
		err := FilterProjectsByEntitlement(idCtx, &projects)
		expected := admin.Projects{
			Projects: []*admin.Project{
				{Id: "p1"},
				{Id: "p2"},
				{Id: "p3"},
			},
		}
		assert.NoError(t, err)
		assert.EqualValues(t, expected.Projects, projects.Projects)
	})
}
