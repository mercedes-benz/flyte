package auth

import (
	"fmt"
	"strings"

	authConfig "github.com/flyteorg/flyte/flyteadmin/auth/config"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
)

var getAuthConfigFunc func() *authConfig.Config

func init() {
	getAuthConfigFunc = getAuthConfig
}

func getAuthConfig() *authConfig.Config {
	return authConfig.GetConfig()
}

// FilterProjectsByEntitlement redacts projects for which current identity has no right to access.
func FilterProjectsByEntitlement(identityContext IdentityContext, projects *admin.Projects) error {
	eligibleProjects, err := eligibleProjectsByIdentity(identityContext)
	if err != nil {
		return err
	}

	filterProjects := []*admin.Project{}

	for _, p := range projects.Projects {
		if _, ok := eligibleProjects[p.Id]; ok {
			filterProjects = append(filterProjects, p)
		}
	}
	projects.Projects = filterProjects

	return nil
}

func eligibleProjectsByIdentity(identityContext IdentityContext) (map[string]any, error) {
	authzCfg := getAuthConfigFunc().ProjectAuthorization

	// App user
	for _, client := range authzCfg.AppAuth.Mappings {
		if client.ClientID == identityContext.AppID() {
			return projectsByEntitlements(client.ProjectSets), nil
		}
	}
	// Human user
	entitlements, err := entitlementsFromClaim(identityContext)
	if err != nil {
		return nil, err
	}
	return projectsByEntitlements(entitlements), nil
}

func projectsByEntitlements(entitlements []string) map[string]any {
	authzCfg := getAuthConfigFunc().ProjectAuthorization
	projects := map[string]any{}
	for _, e := range entitlements {
		// Key (entitlement) must be lowered to match configured ProjectSets.
		// We need this since Viper lowers key of a map by default.
		// GH issue: https://github.com/spf13/viper/issues/411
		if projectSet, ok := authzCfg.ProjectSets[strings.ToLower(e)]; ok {
			for _, p := range projectSet {
				projects[p] = nil
			}
		}
	}
	return projects
}

func entitlementsFromClaim(identityContext IdentityContext) ([]string, error) {
	additionalClaims := identityContext.UserInfo().GetAdditionalClaims().AsMap()
	entitlementClaim := additionalClaims[entitlementClaimKey]

	topLevelEntitlementClaim := identityContext.Claims()[entitlementClaimKey]
	if topLevelEntitlementClaim != nil {
		entitlementClaim = topLevelEntitlementClaim
	}

	return convertClaimToStringSlice(entitlementClaim)
}

func convertClaimToStringSlice(entitlementClaim any) ([]string, error) {
	if entitlementClaim != nil {
		switch claimValue := entitlementClaim.(type) {
		case string:
			return []string{claimValue}, nil
		case []string:
			return claimValue, nil
		case []interface{}:
			entitlements := make([]string, len(claimValue))
			for i, entitlement := range claimValue {
				switch entitlement := entitlement.(type) {
				case string:
					entitlements[i] = entitlement
				default:
					return nil, fmt.Errorf("Failed to convert claim %s", entitlementClaimKey)
				}
			}
			return entitlements, nil
		default:
			return nil, fmt.Errorf("Failed to convert claim %s", entitlementClaimKey)
		}
	}
	return []string{}, nil
}
