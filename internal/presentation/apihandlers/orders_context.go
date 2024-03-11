package apihandlers

import (
	"applicationDesignTest/internal/presentation/apierrors"
	"applicationDesignTest/internal/presentation/apimodels"
	"context"
)

type ContextKeyApiRequest string

const ApiRequestBodyKey ContextKeyApiRequest = "apiRequestBody"

func getFromContext(ctx context.Context) (*apimodels.OrdersAPIRequest, error) {
	req, ok := ctx.Value(ApiRequestBodyKey).(*apimodels.OrdersAPIRequest)
	if !ok || req == nil {
		return nil, apierrors.NewBadRequestError(InvalidRequestDataMessage, nil)
	}
	return req, nil
}
