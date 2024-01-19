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
	commonError "assets/common/custom_error"
	"assets/modules/authorization/model"
	"assets/modules/authorization/service"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var authorizationCntr commonController.HttpController

type authorizationController struct {
	service service.AuthorizationService
}

func GetAuthorizationController() commonController.HttpController {
	if authorizationCntr != nil {
		return authorizationCntr
	}
	authorizationCntr = &authorizationController{service: service.GetAuthorizationService()}
	return authorizationCntr
}

func (controller *authorizationController) RegisterHttpController(router *gin.Engine) {
	authorizationRouter := router.Group("/api/authorization")
	authorizationRouter.POST("/oauth/token", controller.authorize)
}

// authorizationController godoc
// @Summary      authorize
// @Description  Authorize
// @Tags         Authorization controller
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        request	body	string	true	"Authorize request grant_type=enum(password, refresh_token)&username=test&password=test"
// @Success      200		{object}  model.AuthorizationResponse
// @Failure      500
// @Router       /api/authorization/oauth/token [POST]
func (controller *authorizationController) authorize(ctx *gin.Context) {
	log.WithContext(ctx).Info("AuthorizationController: authorize(): Start")

	var response model.AuthorizationResponse

	grantType := ctx.PostForm("grant_type")
	switch grantType {
	case "password":
		username := strings.ToLower(ctx.PostForm("username"))
		password := ctx.PostForm("password")
		response = controller.service.GenerateToken(ctx, username, password)

	case "refresh_token":
		refreshToken := ctx.PostForm("refresh_token")
		response = controller.service.RefreshToken(ctx, refreshToken)

	default:
		panic(commonError.UnknownGrantTypeError)
	}

	ctx.Header("Set-Cookie", fmt.Sprintf("access_token=%s; Expires=%s; Path=/; HttpOnly", response.AccessToken, time.Unix(response.ExpiresAt, 0).UTC().Format(time.RFC1123)))
	ctx.AbortWithStatusJSON(http.StatusOK, response)

	log.WithContext(ctx).Info("AuthorizationController: authorize(): End")
}
