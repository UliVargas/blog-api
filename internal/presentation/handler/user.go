package handler

import (
	"net/http"
	"strconv"

	services "github.com/UliVargas/blog-go/internal/application/service"
	domainService "github.com/UliVargas/blog-go/internal/domain/service"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/UliVargas/blog-go/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domainService.UserServiceInterface
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.HandleError(c, appErrors.ErrInvalidID)
		return
	}
	user, err := h.userService.GetByID(uint(idUint))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
