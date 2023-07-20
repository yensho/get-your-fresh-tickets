package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type helloResponse struct {
	Value string `json:"value"`
}

type BaseHandler struct {
	db  *sqlx.DB
	log *slog.Logger
}

var currentTokens map[string]string

func init() {
	currentTokens = make(map[string]string)
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
		db:  db,
		log: log,
	}
	router.HandleFunc("/auth", base.AuthCreateHandler)
	//router.With(base.BasicAuth).Post("/space", base.createSpace)
	router.Route("/space", func(r chi.Router) {
		r.Use(base.BasicAuth)
		r.Post("/", base.createSpace)
		r.Route("/{spaceName}", func(r chi.Router) {
			r.Get("/", base.getSpace)
			r.Put("/", base.updateSpace)
			r.Delete("/", base.deleteSpace)
		})
	})

	//router.HandleFunc("/hello", BasicAuth(base.HomeHandler))

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

func (b *BaseHandler) AuthCreateHandler(w http.ResponseWriter, r *http.Request) {
	// token := r.Header["token"]
	txn := b.db.MustBeginTx(r.Context(), nil)
	defer txn.Commit()

	var jsonBody AuthCreateRequest
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("auth request failed"))
		b.log.Error(err.Error())
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(jsonBody.Password), 14)
	_, err = txn.Exec(`INSERT INTO auth (user, token) VALUES ($1, $2)`, jsonBody.Username, string(bytes))
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("auth request failed"))
		b.log.Error(err.Error())
		return
	}
	//currentTokens[jsonBody.Username] = string(bytes)
	w.WriteHeader(200)
}

func (b *BaseHandler) BasicAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		var token string
		b.db.Get(&token, `SELECT token FROM auth WHERE user=$1`, user)
		if err := bcrypt.CompareHashAndPassword([]byte(token), []byte(pass)); err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (b *BaseHandler) createSpace(w http.ResponseWriter, r *http.Request) {
	txn := b.db.MustBeginTx(r.Context(), nil)
	defer txn.Commit()

	var createRequest SpaceCreateRequest
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &createRequest)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not unmarshal request json"))
		b.log.Error(err.Error())
		return
	}
	var entries []Space
	for _, area := range createRequest.Areas {
		seats, err := json.Marshal(area.Seats)
		if err != nil {
			continue
		}
		entries = append(entries, Space{
			Name:    createRequest.Name,
			Section: area.Section,
			Seats:   seats,
		})
	}
	_, err = txn.NamedExec(`INSERT INTO spaces (space_nm, space_section_nm, space_section_seats) VALUES (:space_nm, :space_section_nm, :space_section_seats)`, entries)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error inserting records into db"))
		b.log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

func (b *BaseHandler) getSpace(w http.ResponseWriter, r *http.Request) {
	spaceName := chi.URLParam(r, "spaceName")

	var rows []Space
	err := b.db.Select(&rows, `SELECT * FROM spaces WHERE space_nm=?`, spaceName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error querying db"))
		b.log.Error(err.Error())
		return
	}
	resp := SpaceCreateRequest{
		Name: spaceName,
	}
	for _, area := range rows {
		var seats []string
		err := json.Unmarshal(area.Seats, &seats)
		if err != nil {
			continue
		}
		resp.Areas = append(resp.Areas, Area{
			Section: area.Section,
			Seats:   seats,
		})
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error marshaling json response"))
		b.log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)

}

func (b *BaseHandler) updateSpace(w http.ResponseWriter, r *http.Request) {
	txn := b.db.MustBeginTx(r.Context(), nil)
	defer txn.Commit()
	spaceName := chi.URLParam(r, "spaceName")

	var createRequest SpaceCreateRequest
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &createRequest)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not unmarshal request json"))
		b.log.Error(err.Error())
		return
	}

	_, err = txn.NamedExec(`DELETE * FROM spaces WHERE space_nm = $1`, spaceName)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error accessing db"))
		b.log.Error(err.Error())
		return
	}

	var entries []Space
	for _, area := range createRequest.Areas {
		seats, err := json.Marshal(area.Seats)
		if err != nil {
			continue
		}
		entries = append(entries, Space{
			Name:    createRequest.Name,
			Section: area.Section,
			Seats:   seats,
		})
	}
	_, err = txn.NamedExec(`INSERT INTO spaces (space_nm, space_section_nm, space_section_seats) VALUES (:space_nm, :space_section_nm, :space_section_seats)`, entries)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error inserting records into db"))
		b.log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (b *BaseHandler) deleteSpace(w http.ResponseWriter, r *http.Request) {
	spaceName := chi.URLParam(r, "spaceName")

	_, err := b.db.Exec(`DELETE * FROM spaces WHERE space_nm = $1`, spaceName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error accessing db"))
		b.log.Error(err.Error())
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"name": "%s"}`, spaceName)))
}
