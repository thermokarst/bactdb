package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/api"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/types"
)

func handleGetter(g api.Getter) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		id, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 0)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		e, appErr := g.Get(id, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.Marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleLister(l api.Lister) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		opt := r.URL.Query()
		opt.Add("Genus", mux.Vars(r)["genus"])

		claims := helpers.GetClaims(r)

		es, appErr := l.List(&opt, &claims)
		if appErr != nil {
			return appErr
		}
		data, err := es.Marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleUpdater(u api.Updater) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		id, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 0)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		e, err := u.Unmarshal(bodyBytes)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		appErr := u.Update(id, &e, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.Marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleCreater(c api.Creater) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		e, err := c.Unmarshal(bodyBytes)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		appErr := c.Create(&e, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.Marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleDeleter(d api.Deleter) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		id, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 0)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		appErr := d.Delete(id, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		return nil
	}
}
