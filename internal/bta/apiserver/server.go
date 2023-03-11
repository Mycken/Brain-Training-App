package apiserver

import (
	"BrainTApp/internal/bta/store"
	"BrainTApp/internal/model/entity"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	sessionName        = "brainTrainApp"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	resstore     store.ResStore
}

func newServer(store store.Store, resStore store.ResStore, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		resstore:     resStore,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/", s.handleLoginPage()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// works only for /private
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	private.HandleFunc("/display", s.handleDisplay()).Methods("GET")
	private.HandleFunc("/shulte", s.handleShulte()).Methods("GET")
	private.HandleFunc("/shulte/results", s.handleShulteResultsCreate()).Methods("POST")
	private.HandleFunc("/arithmetic", s.handleArithmetic()).Methods("GET")
	private.HandleFunc("/arithmetic/results", s.handleArithmeticResultsCreate()).Methods("POST")
	private.HandleFunc("/memorize", s.handleMemorize()).Methods("GET")
	private.HandleFunc("/memorize/start/{num}", s.handleMemorizeStart()).Methods("GET")
	private.HandleFunc("/memorize/results", s.handleMemorizeResultsCreate()).Methods("POST")

}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)

	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*entity.User))
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &entity.User{
			Email:    req.Email,
			Password: req.Password,
			Username: req.Username,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)

	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

/***/
func (s *server) handleLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the HTML file
		html, err := ioutil.ReadFile("frontend/index.html")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Set the content type header
		w.Header().Set("Content-Type", "text/html")

		// Write the HTML file to the response
		_, err = w.Write(html)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}

func (s *server) handleShulte() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the HTML file
		html, err := ioutil.ReadFile("frontend/tablesShulte.html")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Set the content type header
		w.Header().Set("Content-Type", "text/html")

		// Write the HTML file to the response
		_, err = w.Write(html)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}

func (s *server) handleArithmetic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the HTML file
		html, err := ioutil.ReadFile("frontend/arithmetic.html")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Set the content type header
		w.Header().Set("Content-Type", "text/html")

		// Write the HTML file to the response
		_, err = w.Write(html)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}

func (s *server) handleMemorize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the HTML file
		html, err := ioutil.ReadFile("frontend/memorize.html")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Set the content type header
		w.Header().Set("Content-Type", "text/html")

		// Write the HTML file to the response
		_, err = w.Write(html)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}

func (s *server) handleShulteResultsCreate() http.HandlerFunc {
	type request struct {
		Results  []string      `json:"results"`
		Duration time.Duration `json:"duration"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, ok := r.Context().Value(ctxKeyUser).(*entity.User)
		if !ok {
			return
		}

		result := &entity.Result{
			UserID:      user.ID,
			TestID:      1,
			DateTest:    time.Now(),
			ResultInter: req.Duration,
			Results:     strings.Join(req.Results, " "),
		}

		if err := s.resstore.Result().CreateShulte(result); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, result)
	}
}

func (s *server) handleArithmeticResultsCreate() http.HandlerFunc {
	type request struct {
		Interv      time.Duration `json:"duration"`
		Fails       int           `json:"fails"`
		Expressions int           `json:"expressions"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, ok := r.Context().Value(ctxKeyUser).(*entity.User)
		if !ok {
			return
		}

		result := &entity.Result{
			UserID:      user.ID,
			TestID:      2,
			DateTest:    time.Now(),
			ResultInter: req.Interv,
			ResultOne:   req.Fails,
			ResultTwo:   req.Expressions,
		}

		if err := s.resstore.Result().CreateArithmetic(result); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, result)
	}
}

func (s *server) handleMemorizeResultsCreate() http.HandlerFunc {
	type request struct {
		Interv  time.Duration `json:"duration"`
		Correct int           `json:"correct"`
		Results []string      `json:"results"`
		TestSet []string      `json:"test_set"`
		NumSet  int           `json:"num_set"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, ok := r.Context().Value(ctxKeyUser).(*entity.User)
		if !ok {
			return
		}

		result := &entity.Result{
			UserID:      user.ID,
			TestID:      3,
			DateTest:    time.Now(),
			ResultInter: req.Interv,
			ResultOne:   req.Correct,
			ResultTwo:   req.NumSet,
			Results:     strings.Join(req.Results, " "),
			TestSet:     strings.Join(req.TestSet, " "),
		}

		if err := s.resstore.Result().CreateMemorize(result); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, result)
	}
}

func (s *server) handleMemorizeStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		num, err := strconv.Atoi(vars["num"])
		if err != nil {
			http.Error(w, "Invalid number of words", http.StatusBadRequest)
			return
		}
		partSpeech := 1

		testSet, err := s.resstore.Result().TestSetWords(num, partSpeech)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, testSet)
	}
}

func (s *server) handleDisplay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*entity.User)
		if !ok {
			return
		}

		results, err := s.resstore.Result().FindAll(user.ID)

		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, results)
	}
}
