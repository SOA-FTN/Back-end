package handler

import (
	"context"
	"followings/model"
	"followings/repo"
	"log"
	"net/http"
)

type KeyProduct struct{}

type FollowingHandler struct {
	logger *log.Logger
	// NoSQL: injecting movie repository
	repo *repo.FollowingRepo
}

func NewFollowingHandler(l *log.Logger, r *repo.FollowingRepo) *FollowingHandler {
	return &FollowingHandler{l, r}
}

func (m *FollowingHandler) CreatePerson(rw http.ResponseWriter, h *http.Request) {
	person := h.Context().Value(KeyProduct{}).(*model.User)
	err := m.repo.WritePerson(person)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (m *FollowingHandler) MiddlewarePersonDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		person := &model.User{}
		err := person.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			m.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, person)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
func (m *FollowingHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		m.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
