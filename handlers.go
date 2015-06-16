package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/thermokarst/jwt"
)

type Claims struct {
	Name string
	Iss  string
	Sub  int64
	Role string
	Iat  int64
	Exp  int64
}

func Handler() http.Handler {
	claimsFunc := func(email string) (map[string]interface{}, error) {
		currentTime := time.Now()
		user, err := dbGetUserByEmail(email)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"name": user.Name,
			"iss":  "bactdb",
			"sub":  user.Id,
			"role": user.Role,
			"iat":  currentTime.Unix(),
			"exp":  currentTime.Add(time.Minute * 60 * 24).Unix(),
		}, nil
	}

	verifyClaims := func(claims []byte, r *http.Request) error {
		currentTime := time.Now()
		var c Claims
		err := json.Unmarshal(claims, &c)
		if err != nil {
			return err
		}
		if currentTime.After(time.Unix(c.Exp, 0)) {
			return errors.New("this token has expired")
		}
		context.Set(r, "claims", c)
		return nil
	}

	config := &jwt.Config{
		Secret: os.Getenv("SECRET"),
		Auth:   dbAuthenticate,
		Claims: claimsFunc,
	}

	j, err := jwt.New(config)
	if err != nil {
		panic(err)
	}

	m := mux.NewRouter()

	// Non-auth routes
	m.Handle("/authenticate", tokenHandler(j.GenerateToken())).Methods("POST")

	// Auth routes
	m.Handle("/users", j.Secure(http.HandlerFunc(serveUsersList), verifyClaims)).Methods("GET")
	m.Handle("/users/{Id:.+}", j.Secure(http.HandlerFunc(serveUser), verifyClaims)).Methods("GET")

	// Path-based pattern matching subrouter
	s := m.PathPrefix("/{genus}").Subrouter()

	type r struct {
		f http.HandlerFunc
		m string
		p string
	}

	routes := []r{
		r{handleLister(StrainService{}), "GET", "/strains"},
		r{handleCreater(StrainService{}), "POST", "/strains"},
		r{handleGetter(StrainService{}), "GET", "/strains/{Id:.+}"},
		r{handleUpdater(StrainService{}), "PUT", "/strains/{Id:.+}"},
		r{handleLister(MeasurementService{}), "GET", "/measurements"},
		r{handleGetter(MeasurementService{}), "GET", "/measurements/{Id:.+}"},
		r{handleLister(CharacteristicService{}), "GET", "/characteristics"},
		r{handleGetter(CharacteristicService{}), "GET", "/characteristics/{Id:.+}"},
		r{handleLister(SpeciesService{}), "GET", "/species"},
		r{handleCreater(SpeciesService{}), "POST", "/species"},
		r{handleGetter(SpeciesService{}), "GET", "/species/{Id:.+}"},
		r{handleUpdater(SpeciesService{}), "PUT", "/species/{Id:.+}"},
		r{handleLister(CharacteristicTypeService{}), "GET", "/characteristicTypes"},
		r{handleGetter(CharacteristicTypeService{}), "GET", "/characteristicTypes/{Id:.+}"},
	}

	for _, route := range routes {
		s.Handle(route.p, j.Secure(http.HandlerFunc(route.f), verifyClaims)).Methods(route.m)
	}

	return corsHandler(m)
}

func handleGetter(g getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e, err := g.get(id, mux.Vars(r)["genus"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := e.marshal()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func handleLister(l lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opt := r.URL.Query()
		opt.Add("Genus", mux.Vars(r)["genus"])

		es, err := l.list(&opt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := es.marshal()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func handleUpdater(u updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e, err := u.unmarshal(bodyBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := context.Get(r, "claims")
		var claims Claims = c.(Claims)

		err = u.update(id, &e, claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := e.marshal()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func handleCreater(c creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e, err := c.unmarshal(bodyBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		con := context.Get(r, "claims")
		var claims Claims = con.(Claims)

		err = c.create(&e, claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := e.marshal()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func tokenHandler(h http.Handler) http.Handler {
	token := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Hackish, but we want the token in a JSON object
		w.Write([]byte(`{"token":"`))
		h.ServeHTTP(w, r)
		w.Write([]byte(`"}`))
	}
	return http.HandlerFunc(token)
}

func corsHandler(h http.Handler) http.Handler {
	cors := func(w http.ResponseWriter, r *http.Request) {
		domains := os.Getenv("DOMAINS")
		allowedDomains := strings.Split(domains, ",")
		if origin := r.Header.Get("Origin"); origin != "" {
			for _, s := range allowedDomains {
				if s == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
				}
			}
		}
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(cors)
}
