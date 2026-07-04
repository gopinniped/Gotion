package mapper

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gopinniped/gotion/internal/shared/errs"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func MapError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case errs.TypeNotFound:
			return http.StatusNotFound, appErr.Message
		case errs.TypeConflict:
			return http.StatusConflict, appErr.Message
		case errs.TypeValidation:
			return http.StatusBadRequest, appErr.Message
		case errs.TypeUnauthenticated:
			return http.StatusUnauthorized, appErr.Message
		case errs.TypePermission:
			return http.StatusForbidden, appErr.Message
		case errs.TypeTooManyRequests:
			return http.StatusTooManyRequests, appErr.Message
		case errs.TypeInternal:
			return http.StatusInternalServerError, "internal server error"
		}
	}

	return http.StatusInternalServerError, "internal server error"
}

func SendError(w http.ResponseWriter, err error) {
	statusCode, msg := MapError(err)

	errCode := string(errs.TypeInternal)
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		errCode = string(appErr.Type)
	}

	if statusCode == http.StatusInternalServerError {
		slog.Error("internal server error handled", "err", err.Error())
	}

	resp := ErrorResponse{
		Error: msg,
		Code:  errCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
