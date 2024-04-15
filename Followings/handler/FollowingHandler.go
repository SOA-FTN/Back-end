package handler

import (
	"context"
	"encoding/json"
	"followings/model"
	"followings/repo"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	m.logger.Println("Usao u metodu")
	person := h.Context().Value(KeyProduct{}).(*model.User)
	err := m.repo.WritePerson(person)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (m *FollowingHandler) FollowPerson(rw http.ResponseWriter, h *http.Request) {
	m.logger.Println("Entered FollowPerson method")

	// Decode the request body into a FollowingRelationship struct
	var following repo.FollowingRelationship
	err := json.NewDecoder(h.Body).Decode(&following)
	if err != nil {
		m.logger.Println("Error decoding request body:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Call the repository method to create the following relationship
	err = m.repo.FollowPerson(&following)
	if err != nil {
		m.logger.Println("Database exception:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with a success status code
	rw.WriteHeader(http.StatusCreated)
}

func (m *FollowingHandler) GetFollowRecommendations(rw http.ResponseWriter, r *http.Request) {
	// Extract the username from the URL path
	vars := mux.Vars(r)
	username := vars["username"]

	// Use the username as needed in the handler logic
	// For example, you can pass it to the repository method to retrieve recommendations

	// Call the repository method to retrieve follow recommendations for the specified username
	recommendations, err := m.repo.GetFollowRecommendations(username)
	if err != nil {
		m.logger.Println("Error retrieving follow recommendations:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the list of recommendations to JSON and send it in the response
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(recommendations)
	if err != nil {
		m.logger.Println("Error encoding follow recommendations to JSON:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the status code to OK
	rw.WriteHeader(http.StatusOK)
}

func (m *FollowingHandler) GetFollowedUsers(rw http.ResponseWriter, h *http.Request) {
	m.logger.Println("Entered GetFollowedUsers method")

	// Extract the username from the request parameters
	vars := mux.Vars(h)
	username := vars["username"]

	// Call the repository method to retrieve followed users
	followedUsers, err := m.repo.GetFollowedUsers(username)
	if err != nil {
		m.logger.Println("Database exception:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the list of followed users to JSON
	err = json.NewEncoder(rw).Encode(followedUsers)
	if err != nil {
		m.logger.Println("Error encoding JSON:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with a success status code
	rw.WriteHeader(http.StatusOK)
}

func (fh *FollowingHandler) IsFollowing(rw http.ResponseWriter, req *http.Request) {
	var requestBody struct {
		FollowerUsername string `json:"followerUsername"`
		FollowedUsername string `json:"followedUsername"`
	}

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	isFollowing, err := fh.repo.IsFollowing(requestBody.FollowerUsername, requestBody.FollowedUsername)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the boolean value to a string representation
	var result string
	if isFollowing {
		result = "true"
	} else {
		result = "false"
	}

	// Write the response
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(result))
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

func (m *FollowingHandler) GetUsersExcept(rw http.ResponseWriter, h *http.Request) {
	m.logger.Println("Entered GetUsersExcept method")

	// Extract the username from the request parameters
	vars := mux.Vars(h)
	username := vars["username"]

	// Call the repository method to retrieve users except the specified username
	users, err := m.repo.GetUsersExcept(username)
	if err != nil {
		m.logger.Println("Error retrieving users except:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the list of users to JSON
	err = json.NewEncoder(rw).Encode(users)
	if err != nil {
		m.logger.Println("Error encoding JSON:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with a success status code
	rw.WriteHeader(http.StatusOK)
}
