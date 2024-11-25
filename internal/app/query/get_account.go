package query

import (
	"context"

	"accounts/internal/common/decorator"
	"accounts/internal/domain/account"
)

// GetAccount is a query to get an account.
type GetAccount struct {
	ID account.ID
}

// GetAccountQueryHandler is a query handler for GetAccount.
type GetAccountQueryHandler decorator.QueryHandler[GetAccount, account.Account]

// GetAccountReadModel is a read model for GetAccount.
type GetAccountReadModel interface {
	Get(ctx context.Context, id account.ID) (account.Account, error)
}

// getAccountQueryHandler is a query handler for GetAccount.
type getAccountQueryHandler struct {
	readModel GetAccountReadModel
}

// NewGetAccountQueryHandler creates a new GetAccountQueryHandler.
func NewGetAccountQueryHandler(readModel GetAccountReadModel) GetAccountQueryHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return decorator.ApplyQueryDecorators(&getAccountQueryHandler{readModel: readModel})
}

// Handle handles the GetAccount query.
func (h *getAccountQueryHandler) Handle(ctx context.Context, q GetAccount) (account.Account, error) {
	return h.readModel.Get(ctx, q.ID)
}
