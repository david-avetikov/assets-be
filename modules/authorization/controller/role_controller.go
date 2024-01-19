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

var roleCntr commonController.HttpController

type roleController struct {
	service service.RoleService
}

func GetRoleController() commonController.HttpController {
	if roleCntr != nil {
		return roleCntr
	}
	roleCntr = &roleController{service: service.GetRoleService()}
	return roleCntr
}

func (controller *roleController) RegisterHttpController(router *gin.Engine) {
	roleRouter := router.Group("/api/authorization/roles", commonMiddleware.SecurityHandler)

	roleRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_ROLE"),
		controller.getById,
	)

	roleRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_ROLE"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	roleRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_ROLE"),
		commonResolver.Resolver[model.Role],
		controller.create,
	)

	roleRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_ROLE"),
		commonResolver.Resolver[model.Role],
		controller.update,
	)

	roleRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_ROLE"),
		controller.deleteById,
	)
}

// roleController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get role by id
// @Tags         Role controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Role.ID"
// @Success      200	{object}  model.Role
// @Failure      400
// @Failure      500
// @Router       /api/authorization/roles/{id} [GET]
func (controller *roleController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("RoleController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("RoleController: GetById(): End")
}

// roleController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all role
// @Tags         Role controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.Role
// @Failure      400
// @Failure      500
// @Router       /api/authorization/roles/ [GET]
func (controller *roleController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("RoleController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("RoleController: GetAll(): End")
}

// roleController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create role
// @Tags         Role controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Role  true  "Create Role"
// @Success      201	{object}  model.Role
// @Failure      400
// @Failure      500
// @Router       /api/authorization/roles/ [POST]
func (controller *roleController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("RoleController: Create(): Start")
	role := ctx.MustGet("RequestBody").(model.Role)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, role))
	log.WithContext(ctx).Info("RoleController: Create(): End")
}

// roleController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update role
// @Tags         Role controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Role  true  "Update Role"
// @Success      200	{object}  model.Role
// @Failure      400
// @Failure      500
// @Router       /api/authorization/roles/{id} [PUT]
func (controller *roleController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("RoleController: Update(): Start")
	role := ctx.MustGet("RequestBody").(model.Role)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, role))
	log.WithContext(ctx).Info("RoleController: Update(): End")
}

// roleController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete role by id
// @Tags         Role controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Role.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/authorization/roles/{id} [DELETE]
func (controller *roleController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("RoleController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("RoleController: DeleteById(): End")
}
