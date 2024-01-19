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

var currencyCntr commonController.HttpController

type currencyController struct {
	service service.CurrencyService
}

func GetCurrencyController() commonController.HttpController {
	if currencyCntr != nil {
		return currencyCntr
	}
	currencyCntr = &currencyController{service: service.GetCurrencyService()}
	return currencyCntr
}

func (controller *currencyController) RegisterHttpController(router *gin.Engine) {
	currencyRouter := router.Group("/api/currency/currencies", commonMiddleware.SecurityHandler)

	currencyRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_CURRENCY"),
		controller.getById,
	)

	currencyRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_CURRENCY"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	currencyRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_CURRENCY"),
		commonResolver.Resolver[model.Currency],
		controller.create,
	)

	currencyRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_CURRENCY"),
		commonResolver.Resolver[model.Currency],
		controller.update,
	)

	currencyRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_CURRENCY"),
		controller.deleteById,
	)
}

// currencyController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get currency by id
// @Tags         Currency controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Currency.ID"
// @Success      200	{object}  model.Currency
// @Failure      400
// @Failure      500
// @Router       /api/currency/currencies/{id} [GET]
func (controller *currencyController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CurrencyController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("CurrencyController: GetById(): End")
}

// currencyController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all currency
// @Tags         Currency controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.Currency
// @Failure      400
// @Failure      500
// @Router       /api/currency/currencies/ [GET]
func (controller *currencyController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("CurrencyController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("CurrencyController: GetAll(): End")
}

// currencyController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create currency
// @Tags         Currency controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Currency  true  "Create Currency"
// @Success      201	{object}  model.Currency
// @Failure      400
// @Failure      500
// @Router       /api/currency/currencies/ [POST]
func (controller *currencyController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("CurrencyController: Create(): Start")
	currency := ctx.MustGet("RequestBody").(model.Currency)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, currency))
	log.WithContext(ctx).Info("CurrencyController: Create(): End")
}

// currencyController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update currency
// @Tags         Currency controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Currency  true  "Update Currency"
// @Success      200	{object}  model.Currency
// @Failure      400
// @Failure      500
// @Router       /api/currency/currencies/{id} [PUT]
func (controller *currencyController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("CurrencyController: Update(): Start")
	currency := ctx.MustGet("RequestBody").(model.Currency)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, currency))
	log.WithContext(ctx).Info("CurrencyController: Update(): End")
}

// currencyController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete currency by id
// @Tags         Currency controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Currency.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/currency/currencies/{id} [DELETE]
func (controller *currencyController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CurrencyController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("CurrencyController: DeleteById(): End")
}
