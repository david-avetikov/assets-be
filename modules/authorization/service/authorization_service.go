package service

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
	commonError "assets/common/custom_error"
	commonUtil "assets/common/util"
	"assets/modules/authorization/model"
	"assets/modules/authorization/util"
	"context"
)

var authorizationSrv AuthorizationService

type AuthorizationService interface {
	GenerateToken(ctx context.Context, username string, password string) model.AuthorizationResponse
	RefreshToken(ctx context.Context, refreshToken string) model.AuthorizationResponse
}

type authorizationService struct {
	userService UserService
	roleService RoleService
}

func GetAuthorizationService() AuthorizationService {
	if authorizationSrv != nil {
		return authorizationSrv
	}
	authorizationSrv = &authorizationService{
		userService: GetUserService(),
		roleService: GetRoleService(),
	}
	return authorizationSrv
}

func (service *authorizationService) GenerateToken(ctx context.Context, username string, password string) model.AuthorizationResponse {
	if username == "" || password == "" {
		panic(commonError.UsernameAndPasswordMastNotBeNullError)
	}

	if isAuthorize := anyAuthorize(ctx, username, password,
		service.userService.Authorize,
	); !isAuthorize {
		panic(commonError.UsernameOrPasswordIsIncorrectError)
	}

	user := service.userService.FindByUsername(ctx, username)
	if commonUtil.IsZeroObject(user.ID) {
		panic(commonError.NotFoundError)
	}

	return commonUtil.MustOne(util.NewAuthorizationResponse(user))
}
func (service *authorizationService) RefreshToken(ctx context.Context, refreshToken string) model.AuthorizationResponse {
	token := commonUtil.ParseJwtToken(refreshToken)
	if token.Username == "" {
		panic(commonError.UsernameAndPasswordMastNotBeNullError)
	}

	user := service.userService.FindByUsername(ctx, token.Username)
	if commonUtil.IsZeroObject(user.ID) {
		panic(commonError.NotFoundError)
	}

	return commonUtil.MustOne(util.NewAuthorizationResponse(user))
}

func anyAuthorize(ctx context.Context, username string, password string, funcs ...func(ctx context.Context, username string, password string) bool) bool {
	for _, function := range funcs {
		if function(ctx, username, password) {
			return true
		}
	}
	return false
}
