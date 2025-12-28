package apperrors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/GenkiSugiyama/myapi/api/middlewares"
)

func ErrorHandler(w http.ResponseWriter, req *http.Request, err error) {
	var appErr *MyAppError
	if !errors.As(err, &appErr) {
		appErr = &MyAppError{
			ErrCode: Unknown,
			Message: "internal process failed",
			Err:     err,
		}
	}

	traceID := middlewares.GetTraceID(req.Context())
	log.Printf("[%d] err: %s\n", traceID, appErr)

	var statusCode int
	switch appErr.ErrCode {
	case NAData:
		statusCode = http.StatusNotFound
	case NoTargetData, ReqBodyDecodeFailed, BadParam, StructEncodeFailed:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(appErr)
}
