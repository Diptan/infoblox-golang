package rest

import (
	"address-book/internal/addressbook"
	"address-book/pkg/jsonapi"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Controller is a type that provides user handlers
type Controller struct {
	//log logger.Logger
	srv *addressbook.Service
}

// NewController creates a new instance of the Controller type
func NewController(service *addressbook.Service) *Controller {
	return &Controller{
		srv: service,
	}
}

// RegisterHandlers registers all transport routes for the user handler
func (c *Controller) RegisterHandlers() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/users", c.GetAllUsers)
	router.HandlerFunc(http.MethodGet, "/v1/user/:id", c.GetUserById)
	router.HandlerFunc(http.MethodPost, "/v1/user", c.AddUser)

	router.HandlerFunc(http.MethodDelete, "/v1/user/:id", c.DeleteUser)
	router.HandlerFunc(http.MethodPut, "/v1/user/:id", c.UpdateUser)
	return router
}

func (c *Controller) GetUserById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId := params.ByName("id")

	u, err := c.srv.GetUserById(userId)
	if err != nil {
		c.handleError(w, err)
		return
	}

	jsonapi.Write(w, u)
}

func (c *Controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := c.srv.GetAllUsers()
	if err != nil {
		c.handleError(w, err)
		return
	}

	jsonapi.Write(w, u)
}

// type UserPayload struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// 	Address  string `json:"address"`
// 	Phone    string `json:"phone"`
// }

func (c *Controller) AddUser(w http.ResponseWriter, r *http.Request) {
	var user addressbook.User

	if err := jsonapi.Read(r, &user); err != nil {
		c.handleError(w, err)
		return
	}

	_, err := c.srv.AddUser(user)
	if err != nil {
		c.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId := params.ByName("id")

	err := c.srv.DeleteUser(userId)
	if err != nil {
		c.handleError(w, err)
		return
	}
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId := params.ByName("id")
	var u addressbook.User

	if err := jsonapi.Read(r, &u); err != nil {
		c.handleError(w, err)
		return
	}

	_, err := c.srv.UpdateUser(userId, u)
	if err != nil {
		c.handleError(w, err)
		return
	}
}

func (c *Controller) handleError(w http.ResponseWriter, err error) {
	if status := toStatus(err); status != http.StatusOK {
		//c.log.Warnf("request has failed with error: %s", err)

		w.WriteHeader(status)
		jsonapi.Write(w, jsonapi.ErrorResponse{Errors: []string{err.Error()}})
	}
}
