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

var assetCntr commonController.HttpController

type assetController struct {
	service service.AssetService
}

func GetAssetController() commonController.HttpController {
	if assetCntr != nil {
		return assetCntr
	}
	assetCntr = &assetController{service: service.GetAssetService()}
	return assetCntr
}

func (controller *assetController) RegisterHttpController(router *gin.Engine) {
	assetRouter := router.Group("/api/asset/assets", commonMiddleware.SecurityHandler)

	assetRouter.GET(
		"/:id",
		commonMiddleware.HasAnyAuthorities("READ_ASSET"),
		controller.getById,
	)

	assetRouter.GET(
		"",
		commonMiddleware.HasAnyAuthorities("READ_ASSET"),
		commonMiddleware.PaginationHandler,
		controller.getAll,
	)

	assetRouter.POST(
		"",
		commonMiddleware.HasAnyAuthorities("CREATE_ASSET"),
		commonResolver.Resolver[model.Asset],
		controller.create,
	)

	assetRouter.PUT(
		"/:id",
		commonMiddleware.HasAnyAuthorities("UPDATE_ASSET"),
		commonResolver.Resolver[model.Asset],
		controller.update,
	)

	assetRouter.DELETE(
		"/:id",
		commonMiddleware.HasAnyAuthorities("DELETE_ASSET"),
		controller.deleteById,
	)
}

// assetController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get asset by id
// @Tags         Asset controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Asset.ID"
// @Success      200	{object}  model.Asset
// @Failure      400
// @Failure      500
// @Router       /api/asset/assets/{id} [GET]
func (controller *assetController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("AssetController: GetById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetById(ctx, id))
	log.WithContext(ctx).Info("AssetController: GetById(): End")
}

// assetController godoc
// @Security BearerAuth
// @Summary      getAll
// @Description  Get all asset
// @Tags         Asset controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.Asset
// @Failure      400
// @Failure      500
// @Router       /api/asset/assets/ [GET]
func (controller *assetController) getAll(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	log.WithContext(ctx).Info("AssetController: GetAll(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAll(ctx, page))
	log.WithContext(ctx).Info("AssetController: GetAll(): End")
}

// assetController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create asset
// @Tags         Asset controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Asset  true  "Create Asset"
// @Success      201	{object}  model.Asset
// @Failure      400
// @Failure      500
// @Router       /api/asset/assets/ [POST]
func (controller *assetController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("AssetController: Create(): Start")
	asset := ctx.MustGet("RequestBody").(model.Asset)
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, asset))
	log.WithContext(ctx).Info("AssetController: Create(): End")
}

// assetController godoc
// @Security BearerAuth
// @Summary      update
// @Description  Update asset
// @Tags         Asset controller
// @Accept       json
// @Produce      json
// @Param        user	body	  model.Asset  true  "Update Asset"
// @Success      200	{object}  model.Asset
// @Failure      400
// @Failure      500
// @Router       /api/asset/assets/{id} [PUT]
func (controller *assetController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("AssetController: Update(): Start")
	asset := ctx.MustGet("RequestBody").(model.Asset)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.Update(ctx, asset))
	log.WithContext(ctx).Info("AssetController: Update(): End")
}

// assetController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete asset by id
// @Tags         Asset controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Asset.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/asset/assets/{id} [DELETE]
func (controller *assetController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("AssetController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("AssetController: DeleteById(): End")
}
