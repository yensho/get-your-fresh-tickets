package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	gyftdb "github.com/yensho/get-your-fresh-tickets/db"
	"golang.org/x/exp/slog"
)

type helloResponse struct {
	Value string `json:"value"`
}

type BaseHandler struct {
	db  *gyftdb.GyftDB
	log *slog.Logger
}

func NewApiServer(db *sqlx.DB, log *slog.Logger) *http.Server {
	return &http.Server{
		Handler:      NewRouter(db, log),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func NewRouter(db *sqlx.DB, log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()
	base := BaseHandler{
		db:  gyftdb.NewDb(db),
		log: log,
	}
	//router.HandleFunc("/auth", base.AuthCreateHandler)
	//router.With(base.BasicAuth).Post("/space", base.createSpace)
	router.Route("/space", func(r chi.Router) {
		//r.Use(base.BasicAuth)
		r.Post("/", base.createSpace)
		r.Route("/{spaceName}", func(r chi.Router) {
			r.Get("/", base.getSpace)
			r.Put("/", base.updateSpace)
			r.Delete("/", base.deleteSpace)
		})
	})

	router.HandleFunc("/hello", base.HomeHandler)

	return router
}

func (b *BaseHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := helloResponse{Value: "Hello World!"}
	respBytes, err := json.Marshal(response)
	if err != nil {
		b.log.Error(err.Error())
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// func (b *BaseHandler) AuthCreateHandler(w http.ResponseWriter, r *http.Request) {
// 	// token := r.Header["token"]
// 	txn := b.db.MustBeginTx(r.Context(), nil)
// 	defer txn.Commit()

// 	var jsonBody AuthCreateRequest
// 	err := json.NewDecoder(r.Body).Decode(&jsonBody)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("auth request failed"))
// 		b.log.Error(err.Error())
// 		return
// 	}

// 	bytes, err := bcrypt.GenerateFromPassword([]byte(jsonBody.Password), 14)
// 	if err != nil {
// 		txn.Rollback()
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("auth request failed"))
// 		b.log.Error(err.Error())
// 		return
// 	}
// 	_, err = txn.Exec(`INSERT INTO gyft.auth (user_name, pass_token) VALUES ($1, $2)`, jsonBody.Username, string(bytes))
// 	if err != nil {
// 		txn.Rollback()
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("auth request failed"))
// 		b.log.Error(err.Error())
// 		return
// 	}
// 	w.WriteHeader(200)
// }

// func (b *BaseHandler) BasicAuth(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user, pass, _ := r.BasicAuth()
// 		var token string
// 		err := b.db.Get(&token, `SELECT pass_token FROM gyft.auth WHERE user_name=$1`, user)
// 		if err != nil {
// 			b.log.Error(err.Error())
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		if err := bcrypt.CompareHashAndPassword([]byte(token), []byte(pass)); err != nil {
// 			w.WriteHeader(http.StatusForbidden)
// 			return
// 		}
// 		handler.ServeHTTP(w, r)
// 	})
// }
