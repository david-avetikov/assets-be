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
	"assets/modules/country/model"
	"assets/modules/country/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var cityCntr commonController.HttpController

type cityController struct {
	service service.CityService
}

func GetCityController() commonController.HttpController {
	if cityCntr != nil {
		return cityCntr
	}
	cityCntr = &cityController{service: service.GetCityService()}
	return cityCntr
}

func (controller *cityController) RegisterHttpController(router *gin.Engine) {
	cityRouter := router.Group("/api/city/cities", commonMiddleware.SecurityHandler)

	cityRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_CITY"),
		controller.getById,
	)

	cityRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_CITY"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	cityRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_CITY"),
		commonResolver.Resolver[model.City],
		controller.create,
	)

	cityRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_CITY"),
		commonResolver.Resolver[model.City],
		controller.update,
	)

	cityRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_CITY"),
		controller.deleteById,
	)
}

// cityController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get city by id
// @Tags         City controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "City.ID"
// @Success      200	{object}  model.City
// @Failure      400
// @Failure      500
// @Router       /api/city/cities/{id} [GET]
func (controller *cityController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CityController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("CityController: GetById(): End")
}

// cityController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all city
// @Tags         City controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.City
// @Failure      400
// @Failure      500
// @Router       /api/city/cities/ [GET]
func (controller *cityController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("CityController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("CityController: GetAll(): End")
}

// cityController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create city
// @Tags         City controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.City  true  "Create City"
// @Success      201	{object}  model.City
// @Failure      400
// @Failure      500
// @Router       /api/city/cities/ [POST]
func (controller *cityController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("CityController: Create(): Start")
	city := ctx.MustGet("RequestBody").(model.City)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, city))
	log.WithContext(ctx).Info("CityController: Create(): End")
}

// cityController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update city
// @Tags         City controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.City  true  "Update City"
// @Success      200	{object}  model.City
// @Failure      400
// @Failure      500
// @Router       /api/city/cities/{id} [PUT]
func (controller *cityController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("CityController: Update(): Start")
	city := ctx.MustGet("RequestBody").(model.City)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, city))
	log.WithContext(ctx).Info("CityController: Update(): End")
}

// cityController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete city by id
// @Tags         City controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "City.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/city/cities/{id} [DELETE]
func (controller *cityController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CityController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("CityController: DeleteById(): End")
}
