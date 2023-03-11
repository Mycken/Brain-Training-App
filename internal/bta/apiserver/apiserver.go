package apiserver

import (
	"BrainTApp/internal/bta/store/sqlstore"
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return nil
	}

	defer db.Close()

	store := sqlstore.New(db)
	shstore := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	srv := newServer(store, shstore, sessionStore)

	return http.ListenAndServe(config.BingAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

//type APIServer struct {
//	config *Config
//	logger *logrus.Logger
//	router *mux.Router
//	store  *sqlstore.Store
//}
//
//func New(config *Config) *APIServer {
//
//	return &APIServer{
//		config: config,
//		logger: logrus.New(),
//		router: mux.NewRouter(),
//	}
//}
//
//func (s *APIServer) Start() error {
//
//	if err := s.configureLogger(); err != nil {
//		return err
//	}
//
//	s.configureRouter()
//
//	if err := s.configureStore(); err != nil {
//		return err
//	}
//
//	s.logger.Info("starting api server")
//
//	return http.ListenAndServe(s.config.BingAddr, s.router)
//}
//
//func (s *APIServer) configureLogger() error {
//	level, err := logrus.ParseLevel(s.config.LogLevel)
//	if err != nil {
//		return err
//	}
//
//	s.logger.SetLevel(level)
//
//	return nil
//}
//
//func (s *APIServer) configureRouter() {
//	s.router.HandleFunc("/hi", s.handleHi())
//}
//
//func (s *APIServer) handleHi() http.HandlerFunc {
//	//type request struct {
//	//}
//	return func(w http.ResponseWriter, r *http.Request) {
//		io.WriteString(w, "Hi!!!!")
//	}
//}
//
//func (s *APIServer) configureStore() error {
//	st := sqlstore.New(s.config.Store)
//	if err := st.Open(); err != nil {
//		return err
//	}
//
//	s.store = st
//
//	return nil
//}
