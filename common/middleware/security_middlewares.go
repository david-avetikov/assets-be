package middleware

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
	"assets/common/custom_error"
	"assets/common/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func SecurityHandler(ctx *gin.Context) {
	convertQueryToAccessToken(ctx)
	convertCookieToAccessToken(ctx)

	authorizationHeader := ctx.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		panic(custom_error.NeedAuthorizationHeaderError)
	}

	splitToken := strings.Split(authorizationHeader, " ")
	tokenType, token := splitToken[0], splitToken[1]

	switch strings.ToLower(tokenType) {
	case "bearer":
		util.SetCurrentTokenInfo(ctx, util.ParseJwtToken(token))

	default:
		panic(custom_error.UnsupportedTokenTypeError)
	}
}

func HasAnyRole(roles ...string) func(context *gin.Context) {
	return func(context *gin.Context) {
		existNeededRole := false
		tokenInfo := util.MustGetCurrentTokenInfo(context)
		if util.ArrayContains(tokenInfo.Authorities, "OWNER") {
			return
		}
		for _, authority := range tokenInfo.Roles {
			if util.ArrayContains(roles, authority) {
				existNeededRole = true
				break
			}
		}
		if !existNeededRole {
			panic(custom_error.NotEnoughRightsError)
		}
	}
}

func HasAnyAuthorities(authorities ...string) func(context *gin.Context) {
	return func(context *gin.Context) {
		existNeededAuthority := false
		tokenInfo := util.MustGetCurrentTokenInfo(context)
		if util.ArrayContains(tokenInfo.Authorities, "OWNER") {
			return
		}
		for _, authority := range tokenInfo.Authorities {
			if util.ArrayContains(authorities, authority) {
				existNeededAuthority = true
				break
			}
		}
		if !existNeededAuthority {
			panic(custom_error.NotEnoughRightsError)
		}
	}
}

func convertCookieToAccessToken(ctx *gin.Context) {
	if accessToken, err := ctx.Cookie("access_token"); err == nil && accessToken != "" {
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}
}

func convertQueryToAccessToken(ctx *gin.Context) {
	if accessToken := ctx.Query("access_token"); accessToken != "" {
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}
}
