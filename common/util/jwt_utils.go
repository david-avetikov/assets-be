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
	"assets/common/config"
	"assets/common/custom_error"
	"assets/common/model"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strings"
)

func ParseJwtToken(tokenString string) model.TokenInfo {
	conf := config.CoreConfig.AuthorizationServer
	accessToken := strings.TrimPrefix(tokenString, "Bearer ")
	tokenInto := model.TokenInfo{}

	jwtToken, err := jwt.ParseWithClaims(accessToken, &tokenInto, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != conf.EncodingAlg {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.SignKey), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired") {
			panic(custom_error.TokenExpiredError)
		}
		panic(err)
	}

	if !jwtToken.Valid {
		panic(custom_error.TokenInvalidError)
	}
	return tokenInto
}
