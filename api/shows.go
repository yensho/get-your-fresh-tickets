package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

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
	_, err = txn.NamedExec(`INSERT INTO gyft.spaces (space_nm, space_section_nm, space_section_seats) VALUES (:space_nm, :space_section_nm, :space_section_seats)`, entries)
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
	err := b.db.Select(&rows, `SELECT * FROM gyft.spaces WHERE space_nm=?`, spaceName)
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

	_, err = txn.NamedExec(`DELETE * FROM gyft.spaces WHERE space_nm = $1`, spaceName)
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
	_, err = txn.NamedExec(`INSERT INTO gyft.spaces (space_nm, space_section_nm, space_section_seats) VALUES (:space_nm, :space_section_nm, :space_section_seats)`, entries)
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

	_, err := b.db.Exec(`DELETE * FROM gyft.spaces WHERE space_nm = $1`, spaceName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error accessing db"))
		b.log.Error(err.Error())
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"name": "%s"}`, spaceName)))
}
