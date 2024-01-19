package controller

/*
 * Copyright © 2024, "DEADLINE TEAM" LLC
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are not permitted.
 *
 * THIS SOFTWARE IS PROVIDED BY "DEADLINE TEAM" LLC "AS IS" AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL "DEADLINE TEAM" LLC BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * No reproductions or distributions of this code is permitted without
 * written permission from "DEADLINE TEAM" LLC.
 * Do not reverse engineer or modify this code.
 *
 * © "DEADLINE TEAM" LLC, All rights reserved.
 */

import (
	commonController "assets/common/controller"
	commonMiddleware "assets/common/middleware"
	commonResolver "assets/common/resolver"
	commonUtil "assets/common/util"
	"assets/modules/authorization/model"
	"assets/modules/authorization/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var userCntr commonController.HttpController

type userController struct {
	service service.UserService
}

func GetUserController() commonController.HttpController {
	if userCntr != nil {
		return userCntr
	}
	userCntr = &userController{service: service.GetUserService()}
	return userCntr
}

func (controller *userController) RegisterHttpController(router *gin.Engine) {
	userRouter := router.Group("/api/authorization/users")

	userRouter.GET(
		"/:id",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("READ_USER"),
		controller.getById,
	)

	userRouter.GET(
		"",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("READ_USER"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	userRouter.POST(
		"",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("CREATE_USER"),
		commonResolver.Resolver[model.User],
		controller.create,
	)

	userRouter.PUT(
		"",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("UPDATE_USER"),
		commonResolver.Resolver[model.User],
		controller.update,
	)

	userRouter.DELETE(
		"/:id",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("DELETE_USER"),
		controller.deleteById,
	)

	userRouter.GET(
		"/current",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("READ_USER"),
		controller.getCurrent,
	)

	userRouter.PUT(
		"/:id/roles",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("EDIT_ROLE_USER"),
		controller.addRoles,
	)

	userRouter.PUT(
		"/:id/authorities",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("EDIT_AUTHORITY_USER"),
		controller.addAuthorities,
	)

	userRouter.DELETE(
		"/:id/roles",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("EDIT_ROLE_USER"),
		controller.removeRoles,
	)

	userRouter.DELETE(
		"/:id/authorities",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("EDIT_AUTHORITY_USER"),
		controller.removeAuthorities,
	)
}

// userController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get user by id
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "User.ID"
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id} [GET]
func (controller *userController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("UserController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("UserController: GetById(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all users
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/ [GET]
func (controller *userController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("UserController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("UserController: GetAll(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create a user
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.User  true  "Create User"
// @Success      201	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/ [POST]
func (controller *userController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("UserController: Create(): Start")
	user := ctx.MustGet("RequestBody").(model.User)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, user))
	log.WithContext(ctx).Info("UserController: Create(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update a user
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.User  true  "Update User"
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id} [PUT]
func (controller *userController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("UserController: Update(): Start")
	user := ctx.MustGet("RequestBody").(model.User)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, user))
	log.WithContext(ctx).Info("UserController: Update(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete user by id
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "User.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id} [DELETE]
func (controller *userController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("UserController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("UserController: DeleteById(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      getCurrent
// @Description  Get current user
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/current [GET]
func (controller *userController) getCurrent(ctx *gin.Context) {
	log.WithContext(ctx).Info("UserController: getCurrent(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetCurrent(ctx))
	log.WithContext(ctx).Info("UserController: getCurrent(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      addRoles
// @Description  Add roles to user by id
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "User.ID"
// @Param        rolesIds	query	string  true  "User.Roles"
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id}/roles [PUT]
func (controller *userController) addRoles(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	rolesIds, _ := ctx.GetQueryArray("rolesIds")
	roleUuids := commonUtil.Map(rolesIds, func(it string) uuid.UUID { return uuid.MustParse(it) })
	log.WithContext(ctx).Infof("UserController: AddRoles(id: %s, rolesIds: %s): Start", id, roleUuids)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.AddRoles(ctx, id, roleUuids))
	log.WithContext(ctx).Info("UserController: AddRoles(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      removeRoles
// @Description  Remove roles to user by id
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "User.ID"
// @Param        rolesIds	query	string  true  "User.Roles"
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id}/roles [DELETE]
func (controller *userController) removeRoles(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	rolesIds, _ := ctx.GetQueryArray("rolesIds")
	roleUuids := commonUtil.Map(rolesIds, func(it string) uuid.UUID { return uuid.MustParse(it) })
	log.WithContext(ctx).Infof("UserController: RemoveRoles(id: %s, rolesIds: %s): Start", id, roleUuids)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.RemoveRoles(ctx, id, roleUuids))
	log.WithContext(ctx).Info("UserController: RemoveRoles(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      addAuthorities
// @Description  Add authorities to user by id
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "User.ID"
// @Param        authoritiesIds	query	string  true  "User.Authorities"
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id}/authorities [PUT]
func (controller *userController) addAuthorities(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	authoritiesIds, _ := ctx.GetQueryArray("authoritiesIds")
	authorityUuids := commonUtil.Map(authoritiesIds, func(it string) uuid.UUID { return uuid.MustParse(it) })
	log.WithContext(ctx).Infof("UserController: AddAuthorities(id: %s, authoritiesIds: %s): Start", id, authorityUuids)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.AddAuthorities(ctx, id, authorityUuids))
	log.WithContext(ctx).Info("UserController: AddAuthorities(): End")
}

// userController godoc
// @Security BearerAuth
// @Summary      removeAuthorities
// @Description  Remove authorities to user by id
// @Tags         User controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "User.ID"
// @Param        authoritiesIds	query	string  true  "User.Authorities"
// @Success      200	{object}  model.User
// @Failure      400
// @Failure      500
// @Router       /api/authorization/users/{id}/authorities [DELETE]
func (controller *userController) removeAuthorities(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	authoritiesIds, _ := ctx.GetQueryArray("authoritiesIds")
	authorityUuids := commonUtil.Map(authoritiesIds, func(it string) uuid.UUID { return uuid.MustParse(it) })
	log.WithContext(ctx).Infof("UserController: RemoveAuthorities(id: %s, authoritiesIds: %s): Start", id, authorityUuids)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.RemoveAuthorities(ctx, id, authorityUuids))
	log.WithContext(ctx).Info("UserController: RemoveAuthorities(): End")
}
