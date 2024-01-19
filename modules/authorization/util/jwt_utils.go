package util

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
	commonConfig "assets/common/config"
	commonModel "assets/common/model"
	commonUtil "assets/common/util"
	"assets/modules/authorization/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

func NewAuthorizationResponse(user model.User) (model.AuthorizationResponse, error) {
	conf := commonConfig.CoreConfig.AuthorizationServer

	accessToken, id, err := GenerateJwtToken(user, conf.AccessTokenValiditySeconds, "")
	if err != nil {
		return model.NilAuthorizationResponse, err
	}
	refreshToken, _, err := GenerateJwtToken(user, conf.RefreshTokenValiditySeconds, id)
	if err != nil {
		return model.NilAuthorizationResponse, err
	}
	accessTokenString, err := accessToken.SignedString([]byte(conf.SignKey))
	if err != nil {
		return model.NilAuthorizationResponse, err
	}
	refreshTokenString, err := refreshToken.SignedString([]byte(conf.SignKey))
	if err != nil {
		return model.NilAuthorizationResponse, err
	}

	var response model.AuthorizationResponse
	response.TokenInfo = accessToken.Claims.(commonModel.TokenInfo)
	response.AccessToken = accessTokenString
	response.TokenType = "bearer"
	response.RefreshToken = refreshTokenString

	return response, nil
}

func GenerateJwtToken(user model.User, expiredAfterSecond int, parentId string) (*jwt.Token, string, error) {
	conf := commonConfig.CoreConfig.AuthorizationServer

	id := uuid.New().String()
	issuedAt := time.Now().Unix()
	expiresAt := issuedAt + int64(expiredAfterSecond)

	roles := commonUtil.Map(user.Roles, func(role *model.Role) string { return role.Name })

	var authorities []string
	for _, role := range user.Roles {
		authoritiesFromRole := commonUtil.Map(role.Authorities, func(authority *model.Authority) string { return authority.Method })
		authorities = append(authorities, authoritiesFromRole...)
	}
	authoritiesFromUser := commonUtil.Map(user.AdditionalAuthorities, func(authority *model.Authority) string { return authority.Method })
	authorities = append(authorities, authoritiesFromUser...)
	authorities = commonUtil.Unique(authorities)

	claims := commonModel.TokenInfo{
		UserId:      user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Roles:       roles,
		Authorities: authorities,
		ParentId:    parentId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Id:        id,
			IssuedAt:  issuedAt,
			Issuer:    "DeadlineTeam",
		},
	}

	return jwt.NewWithClaims(jwt.GetSigningMethod(conf.EncodingAlg), claims), id, nil
}
