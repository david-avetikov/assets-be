package main

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
	"assets/common/config"
	"assets/common/controller"
	"assets/common/embedded"
	"assets/common/server"
	"assets/common/util"
	assetController "assets/modules/asset/controller"
	attachmentController "assets/modules/attachment/controller"
	authorizationController "assets/modules/authorization/controller"
	countryController "assets/modules/country/controller"
	"github.com/gin-gonic/gin"
)

// @title           assets API
// @version         0.1

// @contact.name   API Support
// @contact.url    https://deadline.team
// @contact.email  info@deadline.team

// @host      assets.deadline.team

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	util.InitApplication("assets")
	util.CheckedError(config.Load(&config.CoreConfig))

	util.CheckedError(embedded.CreateEmbeddedPostgres())
	util.CheckedError(embedded.CreateEmbeddedRedis())
	util.CheckedError(embedded.CreateEmbeddedJaeger())

	util.InitTracing()

	util.CheckedRun(server.NewHttpServer(func(router *gin.Engine) {
		controller.Register(router, assetController.GetAssetController())
		controller.Register(router, assetController.GetCashFlowController())

		controller.Register(router, attachmentController.GetAttachmentController())

		controller.Register(router, authorizationController.GetAuthorityController())
		controller.Register(router, authorizationController.GetRoleController())
		controller.Register(router, authorizationController.GetUserController())
		controller.Register(router, authorizationController.GetAuthorizationController())

		controller.Register(router, countryController.GetCityController())
		controller.Register(router, countryController.GetCurrencyController())
		controller.Register(router, countryController.GetCountryController())
	}))
}
