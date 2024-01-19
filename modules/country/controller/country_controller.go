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

var countryCntr commonController.HttpController

type countryController struct {
	service service.CountryService
}

func GetCountryController() commonController.HttpController {
	if countryCntr != nil {
		return countryCntr
	}
	countryCntr = &countryController{service: service.GetCountryService()}
	return countryCntr
}

func (controller *countryController) RegisterHttpController(router *gin.Engine) {
	countryRouter := router.Group("/api/country/countries", commonMiddleware.SecurityHandler)

	countryRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_COUNTRY"),
		controller.getById,
	)

	countryRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_COUNTRY"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	countryRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_COUNTRY"),
		commonResolver.Resolver[model.Country],
		controller.create,
	)

	countryRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_COUNTRY"),
		commonResolver.Resolver[model.Country],
		controller.update,
	)

	countryRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_COUNTRY"),
		controller.deleteById,
	)
}

// countryController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get country by id
// @Tags         Country controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Country.ID"
// @Success      200	{object}  model.Country
// @Failure      400
// @Failure      500
// @Router       /api/country/countries/{id} [GET]
func (controller *countryController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CountryController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("CountryController: GetById(): End")
}

// countryController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all country
// @Tags         Country controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.Country
// @Failure      400
// @Failure      500
// @Router       /api/country/countries/ [GET]
func (controller *countryController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("CountryController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("CountryController: GetAll(): End")
}

// countryController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create country
// @Tags         Country controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Country  true  "Create Country"
// @Success      201	{object}  model.Country
// @Failure      400
// @Failure      500
// @Router       /api/country/countries/ [POST]
func (controller *countryController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("CountryController: Create(): Start")
	country := ctx.MustGet("RequestBody").(model.Country)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, country))
	log.WithContext(ctx).Info("CountryController: Create(): End")
}

// countryController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update country
// @Tags         Country controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Country  true  "Update Country"
// @Success      200	{object}  model.Country
// @Failure      400
// @Failure      500
// @Router       /api/country/countries/{id} [PUT]
func (controller *countryController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("CountryController: Update(): Start")
	country := ctx.MustGet("RequestBody").(model.Country)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, country))
	log.WithContext(ctx).Info("CountryController: Update(): End")
}

// countryController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete country by id
// @Tags         Country controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Country.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/country/countries/{id} [DELETE]
func (controller *countryController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CountryController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("CountryController: DeleteById(): End")
}
