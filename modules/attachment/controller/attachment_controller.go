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
	commonUtil "assets/common/util"
	"assets/modules/attachment/model"
	"assets/modules/attachment/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var attchCntr commonController.HttpController

type attachmentController struct {
	service service.AttachmentService
}

func GetAttachmentController() commonController.HttpController {
	if attchCntr != nil {
		return attchCntr
	}
	attchCntr = &attachmentController{service: service.GetAttachmentService()}
	return attchCntr
}

func (controller *attachmentController) RegisterHttpController(router *gin.Engine) {
	attachmentRouter := router.Group("/api/attachment/attachments")

	attachmentRouter.GET(
		"/:id",
		controller.getById,
	)

	attachmentRouter.GET(
		"/:id/meta",
		controller.getMetaById,
	)

	attachmentRouter.GET(
		"/meta",
		commonMiddleware.PaginationHandler,
		controller.getAllMeta,
	)

	attachmentRouter.POST(
		"",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("CREATE_ATTACHMENT"),
		controller.create,
	)

	attachmentRouter.DELETE(
		"/:id",
		commonMiddleware.SecurityHandler,
		commonMiddleware.HasAnyAuthorities("DELETE_ATTACHMENT"),
		controller.deleteById,
	)
}

// attachmentController godoc
// @Security BearerAuth
// @Summary      getById
// @Description  Get attachment by id
// @Tags         Attachment controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Attachment.ID"
// @Success      200	{object}  model.Attachment
// @Failure      400
// @Failure      500
// @Router       /api/attachment/attachments/{id} [GET]
func (controller *attachmentController) getById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("AttachmentController: GetById(id: %s): Start", id)
	data, contentType, fileName, _ := controller.service.GetById(ctx, id)
	ctx.Header("content-disposition", `attachment; filename=`+fileName)
	ctx.Header("cache-control", `max-age=604800`)
	ctx.Data(http.StatusOK, contentType, data)
	ctx.Abort()
	log.WithContext(ctx).Info("AttachmentController: GetById(): End")
}

// attachmentController godoc
// @Security BearerAuth
// @Summary      getMetaById
// @Description  Get attachment meta by id
// @Tags         Attachment controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Attachment.ID"
// @Success      200	{object}  model.Attachment
// @Failure      400
// @Failure      500
// @Router       /api/attachment/attachments/{id}/meta [GET]
func (controller *attachmentController) getMetaById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	fields := commonUtil.MustGetQueryFields(ctx)
	log.WithContext(ctx).Infof("AttachmentController: GetMetaById(id: %s): Start", id)
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetMetaById(ctx, id, fields))
	log.WithContext(ctx).Info("AttachmentController: GetMetaById(): End")
}

// attachmentController godoc
// @Security BearerAuth
// @Summary      getAllMeta
// @Description  Get all attachment
// @Tags         Attachment controller
// @Accept       json
// @Produce      json
// @Success      200	{array}  model.Attachment
// @Failure      400
// @Failure      500
// @Router       /api/attachment/attachments/meta [GET]
func (controller *attachmentController) getAllMeta(ctx *gin.Context) {
	page := commonUtil.MustGetPageable(ctx)
	filter := commonUtil.MustGetFilterObject[model.AttachmentFilter](ctx)
	fields := commonUtil.MustGetQueryFields(ctx)
	log.WithContext(ctx).Info("AttachmentController: GetAllMeta(): Start")
	ctx.AbortWithStatusJSON(http.StatusOK, controller.service.GetAllMeta(ctx, page, filter, fields))
	log.WithContext(ctx).Info("AttachmentController: GetAllMeta(): End")
}

// attachmentController godoc
// @Security BearerAuth
// @Summary      create
// @Description  Create attachment
// @Tags         Attachment controller
// @Accept       multipart/form-data
// @Produce      json
// @Param        file	body	  model.Attachment  true  "Create Attachment"
// @Success      201	{object}  model.Attachment
// @Failure      400
// @Failure      500
// @Router       /api/attachment/attachments/ [POST]
func (controller *attachmentController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("AttachmentController: Create(): Start")
	fileHeader := commonUtil.MustOne(ctx.FormFile("file"))
	ctx.AbortWithStatusJSON(http.StatusCreated, controller.service.Create(ctx, fileHeader))
	log.WithContext(ctx).Info("AttachmentController: Create(): End")
}

// attachmentController godoc
// @Security BearerAuth
// @Summary      deleteById
// @Description  Delete attachment by id
// @Tags         Attachment controller
// @Accept       json
// @Produce      json
// @Param        id		path     string  true  "Attachment.ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /api/attachment/attachments/{id} [DELETE]
func (controller *attachmentController) deleteById(ctx *gin.Context) {
	id := uuid.MustParse(ctx.Param("id"))
	log.WithContext(ctx).Infof("AttachmentController: DeleteById(id: %s): Start", id)
	controller.service.DeleteById(ctx, id)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "OK"})
	log.WithContext(ctx).Info("AttachmentController: DeleteById(): End")
}
