package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	"strconv"
)

type UserContext struct {
	*Context
	Id   int
	User *model.User
}

var userRouter *web.Router

// @SubApi User API [/users]
func init() {
	userRouter = publicRouter.Subrouter(UserContext{}, "/users")
	userRouter.Middleware((*UserContext).Authorized)
	userRouter.Middleware((*UserContext).FetchUser)
	userRouter.Get("/:id", (*UserContext).GetUser)
	userRouter.Put("/:id", (*UserContext).UpdateUser)
	userRouter.Delete("/:id", (*UserContext).DeleteUser)
}

func (c *UserContext) Authorized(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.CurrentUser == nil {
		c.Unauthorized("a user or admin is required")
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest("a valid id is required")
		return
	}
	c.Id = id

	if (c.Id != 0 && c.CurrentUser.Id == c.Id) || c.CurrentUser.Admin {
		next(rw, r)
	} else {
		c.Unauthorized("a user or admin is required")
	}
}

func (c *UserContext) FetchUser(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.Id == 0 {
		c.BadRequest("a valid id is required")
		return
	}

	user, err := servicer.GetUser(c.Id)
	if err != nil {
		if err == servicer.RecordNotFound {
			c.NotFound("user not found")
		} else {
			c.InternalServerError("internal server error", err)
		}
		return
	}

	c.User = user
	next(rw, r)
}

// @Title getUser
// @Description Get a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Success 200 {object}     model.User                 ""
// @Failure 401 {object}     MessageResponse           	""
// @Router /users/{id} [get]
func (c *UserContext) GetUser(rw web.ResponseWriter, r *web.Request) {
	c.ResponseJson(c.User)
}

// @Title updateUser
// @Description Update a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Param   email            body   string  false       "A new email address"
// @Param   name             body   string  false       "A new name"
// @Param   password         body   string  false       "A new password"
// @Success 200 {object}     model.User                  ""
// @Failure 401 {object}     MessageResponse           	 ""
// @Router /users/{id} [put]
func (c *UserContext) UpdateUser(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest("malformed request")
		return
	}

	err = servicer.UpdateUser(c.User, request.Email, request.Name, request.Password)
	if err != nil {
		c.InternalServerError("internal server error", err)
		return
	}

	c.ResponseJson(c.User)
}

// @Title deleteUser
// @Description Update a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Success 200 {object}     MessageResponse            ""
// @Failure 401 {object}     MessageResponse           	""
// @Router /users/{id} [delete]
func (c *UserContext) DeleteUser(rw web.ResponseWriter, r *web.Request) {
	err := servicer.DeleteUser(c.Id)
	if err != nil {
		c.InternalServerError("internal server error", err)
		return
	}

	c.ResponseJson(&MessageResponse{Message: "user has been deleted"})
}
