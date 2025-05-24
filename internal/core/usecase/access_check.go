package usecase

import (
	"context"
	"fmt"
	"slices"

	"ghorkov32/proletariat-budget-be/internal/common"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"github.com/inv-cloud-platform/hub-com-auth-go/hubauth"
	"github.com/inv-cloud-platform/hub-com-tools-go/hubmiddlewares"
	"github.com/rs/zerolog/log"
)

type AccessChecker struct {
	scopeChecker port.ScopeLookup
}

func NewAccessChecker(scopeLookup port.ScopeLookup) *AccessChecker {
	return &AccessChecker{
		scopeChecker: scopeLookup,
	}
}

func (a AccessChecker) Check(ctx context.Context, action string, ids ...string) (allowedResources []string, err error) {
	reqId := hubmiddlewares.GetRequestId(ctx)

	if has := hubauth.HasAction(ctx, common.Domain, action+"_all"); has {
		return allowedResources, nil
	}

	if has := hubauth.HasAction(ctx, common.Domain, action); !has {
		errValidation := fmt.Errorf("%w: %s", domain.ErrForbiddenAction, action)
		log.Warn().Str("request_id", reqId).Str("email", hubauth.GetEmail(ctx)).Err(errValidation).Str("action", action).Send()

		return allowedResources, errValidation
	}

	if action == common.Create {
		return allowedResources, nil
	}

	scopes, errLookup := a.scopeChecker.GetScopesV1(ctx, common.Domain, action)
	if errLookup != nil {
		errValidation := fmt.Errorf("%w: %w", domain.ErrUnableToFetchScopes, errLookup)
		log.Warn().Str("request_id", reqId).Str("email", hubauth.GetEmail(ctx)).Err(errValidation).Str("action", action).Send()

		return allowedResources, errValidation
	}

	myResourceIdList, errTranslate := a.translateToListOfMyResources(scopes)
	if errTranslate != nil {
		log.Error().Str("request_id", reqId).Str("email", hubauth.GetEmail(ctx)).Err(errTranslate).Str("action", action).Msg("unable to translate scopes")
	}

	allowedResources = append(allowedResources, myResourceIdList...)

	if len(allowedResources) == 0 {
		log.Warn().Str("request_id", reqId).Str("email", hubauth.GetEmail(ctx)).Err(domain.ErrForbiddenResource).Str("action", action).Send()

		return allowedResources, domain.ErrForbiddenResource
	}

	if hasIdsAccess, id := a.hasAccessToAllIds(allowedResources, ids); !hasIdsAccess {
		log.
			Warn().
			Str("request_id", reqId).
			Str("email", hubauth.GetEmail(ctx)).
			Err(domain.ErrForbiddenResource).
			Str("action", action).
			Str("missing_resource", id).
			Send()

		return allowedResources, domain.ErrForbiddenResource
	}

	return allowedResources, nil
}

func (a AccessChecker) GetEmail(ctx context.Context) string {
	return hubauth.GetEmail(ctx)
}

// TODO: Implement your logic
func (a *AccessChecker) translateToListOfMyResources(_ hubauth.ScopeV1) ([]string, error) {
	return []string{}, nil
}

func (a *AccessChecker) hasAccessToAllIds(resources, ids []string) (hasAccess bool, missingId string) {
	for i := range ids {
		if !slices.Contains(resources, ids[i]) {
			return false, ids[i]
		}
	}

	return true, ""
}
