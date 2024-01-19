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
	"assets/modules/asset/model"
	"assets/modules/asset/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var cashFlowCntr commonController.HttpController

type cashFlowController struct {
	service service.CashFlowService
}

func GetCashFlowController() commonController.HttpController {
	if cashFlowCntr != nil {
		return cashFlowCntr
	}
	cashFlowCntr = &cashFlowController{service: service.GetCashFlowService()}
	return cashFlowCntr
}

func (controller *cashFlowController) RegisterHttpController(router *gin.Engine) {
	cashFlowRouter := router.Group("/api/cashFlow/cashFlows", commonMiddleware.SecurityHandler)

	cashFlowRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_CASH_FLOW"),
		controller.getById,
	)

	cashFlowRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_CASH_FLOW"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	cashFlowRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_CASH_FLOW"),
		commonResolver.Resolver[model.CashFlow],
		controller.create,
	)

	cashFlowRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_CASH_FLOW"),
		commonResolver.Resolver[model.CashFlow],
		controller.update,
	)

	cashFlowRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_CASH_FLOW"),
		controller.deleteById,
	)
}

// cashFlowController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get cashFlow by id
// @Tags         CashFlow controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "CashFlow.ID"
// @Success      200	{object}  model.CashFlow
// @Failure      400
// @Failure      500
// @Router       /api/cashFlow/cashFlows/{id} [GET]
func (controller *cashFlowController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CashFlowController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("CashFlowController: GetById(): End")
}

// cashFlowController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all cashFlow
// @Tags         CashFlow controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.CashFlow
// @Failure      400
// @Failure      500
// @Router       /api/cashFlow/cashFlows/ [GET]
func (controller *cashFlowController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("CashFlowController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("CashFlowController: GetAll(): End")
}

// cashFlowController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create cashFlow
// @Tags         CashFlow controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.CashFlow  true  "Create CashFlow"
// @Success      201	{object}  model.CashFlow
// @Failure      400
// @Failure      500
// @Router       /api/cashFlow/cashFlows/ [POST]
func (controller *cashFlowController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("CashFlowController: Create(): Start")
	cashFlow := ctx.MustGet("RequestBody").(model.CashFlow)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, cashFlow))
	log.WithContext(ctx).Info("CashFlowController: Create(): End")
}

// cashFlowController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update cashFlow
// @Tags         CashFlow controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.CashFlow  true  "Update CashFlow"
// @Success      200	{object}  model.CashFlow
// @Failure      400
// @Failure      500
// @Router       /api/cashFlow/cashFlows/{id} [PUT]
func (controller *cashFlowController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("CashFlowController: Update(): Start")
	cashFlow := ctx.MustGet("RequestBody").(model.CashFlow)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, cashFlow))
	log.WithContext(ctx).Info("CashFlowController: Update(): End")
}

// cashFlowController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete cashFlow by id
// @Tags         CashFlow controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "CashFlow.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/cashFlow/cashFlows/{id} [DELETE]
func (controller *cashFlowController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("CashFlowController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("CashFlowController: DeleteById(): End")
}
