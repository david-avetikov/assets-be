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

var authorityCntr commonController.HttpController

type authorityController struct {
	service service.AuthorityService
}

func GetAuthorityController() commonController.HttpController {
	if authorityCntr != nil {
		return authorityCntr
	}
	authorityCntr = &authorityController{service: service.GetAuthorityService()}
	return authorityCntr
}

func (controller *authorityController) RegisterHttpController(router *gin.Engine) {
	authorityRouter := router.Group("/api/authorization/authorities", commonMiddleware.SecurityHandler)

	authorityRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_AUTHORITY"),
		controller.getById,
	)

	authorityRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_AUTHORITY"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	authorityRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_AUTHORITY"),
		commonResolver.Resolver[model.Authority],
		controller.create,
	)

	authorityRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_AUTHORITY"),
		commonResolver.Resolver[model.Authority],
		controller.update,
	)

	authorityRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_AUTHORITY"),
		controller.deleteById,
	)
}

// authorityController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get authority by id
// @Tags         Authority controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Authority.ID"
// @Success      200	{object}  model.Authority
// @Failure      400
// @Failure      500
// @Router       /api/authorization/authorities/{id} [GET]
func (controller *authorityController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("AuthorityController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("AuthorityController: GetById(): End")
}

// authorityController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all authority
// @Tags         Authority controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.Authority
// @Failure      400
// @Failure      500
// @Router       /api/authorization/authorities/ [GET]
func (controller *authorityController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("AuthorityController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("AuthorityController: GetAll(): End")
}

// authorityController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create authority
// @Tags         Authority controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Authority  true  "Create Authority"
// @Success      201	{object}  model.Authority
// @Failure      400
// @Failure      500
// @Router       /api/authorization/authorities/ [POST]
func (controller *authorityController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("AuthorityController: Create(): Start")
	authority := ctx.MustGet("RequestBody").(model.Authority)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, authority))
	log.WithContext(ctx).Info("AuthorityController: Create(): End")
}

// authorityController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update authority
// @Tags         Authority controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Authority  true  "Update Authority"
// @Success      200	{object}  model.Authority
// @Failure      400
// @Failure      500
// @Router       /api/authorization/authorities/{id} [PUT]
func (controller *authorityController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("AuthorityController: Update(): Start")
	authority := ctx.MustGet("RequestBody").(model.Authority)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, authority))
	log.WithContext(ctx).Info("AuthorityController: Update(): End")
}

// authorityController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete authority by id
// @Tags         Authority controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Authority.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/authorization/authorities/{id} [DELETE]
func (controller *authorityController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("AuthorityController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("AuthorityController: DeleteById(): End")
}
