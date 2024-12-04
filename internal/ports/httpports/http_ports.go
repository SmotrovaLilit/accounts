package httpports

import (
	"accounts/internal/app"
	"accounts/internal/app/commands"
	"accounts/internal/app/query"
	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"
	"accounts/internal/gen-openapi/server"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// Server is a HTTP server.
type Server struct {
	app app.Application
}

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(application app.Application) Server {
	return Server{
		app: application,
	}
}

// CreateAccount creates a new account.
func (h Server) CreateAccount(w http.ResponseWriter, r *http.Request) {
	requestData := &server.CreateAccountRequest{}
	if err := render.Decode(r, requestData); err != nil {
		respondWithError(err, w, r)
		return
	}
	id, err := account.NewID(*requestData.Id)
	if err != nil {
		respondWithError(err, w, r)
		return
	}
	ownerID, err := owner.NewID(*requestData.OwnerId)
	if err != nil {
		respondWithError(err, w, r)
		return
	}
	currency, err := toCurrency(string(*requestData.Currency))
	if err != nil {
		respondWithError(err, w, r)
		return
	}

	err = h.app.Commands.CreateAccount.Handle(r.Context(), commands.CreateAccount{
		ID:       id,
		OwnerID:  ownerID,
		Currency: currency,
	})
	if err != nil {
		respondWithError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Server) GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := account.NewID(r.URL.Query().Get("id"))
	if err != nil {
		respondWithError(err, w, r)
		return
	}
	acc, err := h.app.Queries.GetAccount.Handle(
		r.Context(),
		query.GetAccount{
			ID: id,
		})
	if err != nil {
		respondWithError(err, w, r)
		return
	}
	render.JSON(w, r, accountToResponse(acc))
}

func accountToResponse(acc account.Account) *server.Account {
	return &server.Account{
		Balance:   acc.BalanceAmount(),
		CreatedAt: acc.CreatedAt(),
		Currency:  toResponseCurrency(acc.BalanceCurrency()),
		Id:        acc.ID().String(),
		OwnerId:   acc.OwnerID().String(),
		UpdatedAt: acc.UpdatedAt(),
	}
}

func toResponseCurrency(currency account.Currency) server.Currency {
	switch currency {
	case account.CurrencyUSD:
		return server.CURRENCYUSD
	case account.CurrencyEUR:
		return server.CURRENCYEUR
	default:
		panic(fmt.Sprintf("unknown currency: %v", currency))
	}
}

func toCurrency(currency string) (account.Currency, error) {
	switch currency {
	case "USD":
		return account.CurrencyUSD, nil
	case "EUR":
		return account.CurrencyEUR, nil
	default:
		return account.Currency{}, errors.New("invalid currency")
	}
}

func respondWithError(err error, w http.ResponseWriter, _ *http.Request) {
	// TODO switch on error type and return appropriate status code
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
