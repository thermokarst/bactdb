package handlers

import (
	"fmt"
	"net/http"

	"github.com/thermokarst/bactdb/types"
)

type errorHandler func(http.ResponseWriter, *http.Request) *types.AppError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		w.WriteHeader(err.Status)
		fmt.Fprintln(w, err.Error.Error())
	}
}
