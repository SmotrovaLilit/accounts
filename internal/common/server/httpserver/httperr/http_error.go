package httperr

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func InternalError(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, "INTERNAL", "Internal server error", w, r, http.StatusInternalServerError)
}

func Unauthorised(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, "UNAUTHORISED", "Unauthorised", w, r, http.StatusUnauthorized)
}

func BadRequest(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, "BAD_REQUEST", "Bad request", w, r, http.StatusBadRequest)
}

func ValidationError(err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, "VALIDATION_ERROR", "Validation error", w, r, http.StatusUnprocessableEntity)
}

func httpRespondWithError(err error, code, message string, w http.ResponseWriter, r *http.Request, status int) {
	slog.With("error_code", code, "err", err.Error()).WarnContext(r.Context(), message) // TODO Think about level later
	resp := ErrorResponse{code, message, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err) // TODO think about it
	}
}

type ErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
