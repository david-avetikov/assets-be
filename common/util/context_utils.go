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
	"assets/common/model"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	tokenInfoKey    = "TokenInfo"
	pageableKey     = "Pageable"
	queryFieldsKey  = "QueryFields"
	filterObjectKey = "FilterObject"
)

func SetToContext[T any](ctx context.Context, key string, value T) {
	ConvertContext(ctx).Set(key, value)
}

func GetFromContext[T any](ctx context.Context, key string) (T, bool) {
	value, ok := ConvertContext(ctx).Get(key)
	if !ok {
		var zeroValue T
		return zeroValue, ok
	}

	return value.(T), ok
}

func SetCurrentTokenInfo(ctx context.Context, tokenInfo model.TokenInfo) {
	SetToContext(ConvertContext(ctx), tokenInfoKey, tokenInfo)
}

func GetCurrentTokenInfo(ctx context.Context) (model.TokenInfo, bool) {
	return GetFromContext[model.TokenInfo](ConvertContext(ctx), tokenInfoKey)
}

func MustGetCurrentTokenInfo(ctx context.Context) model.TokenInfo {
	result, _ := GetFromContext[model.TokenInfo](ConvertContext(ctx), tokenInfoKey)
	return result
}

func SetPageable(ctx context.Context, pageable model.Pageable) {
	SetToContext(ConvertContext(ctx), pageableKey, pageable)
}

func GetPageable(ctx context.Context) (model.Pageable, bool) {
	return GetFromContext[model.Pageable](ConvertContext(ctx), pageableKey)
}

func MustGetPageable(ctx context.Context) model.Pageable {
	result, _ := GetFromContext[model.Pageable](ConvertContext(ctx), pageableKey)
	return result
}

func SetQueryFields(ctx context.Context, fields string) {
	SetToContext(ConvertContext(ctx), queryFieldsKey, fields)
}

func GetQueryFields(ctx context.Context) (string, bool) {
	return GetFromContext[string](ConvertContext(ctx), queryFieldsKey)
}

func MustGetQueryFields(ctx context.Context) string {
	result, _ := GetFromContext[string](ConvertContext(ctx), queryFieldsKey)
	return result
}

func GetIntQuery(ctx context.Context, key string) int {
	valStr, ok := ConvertContext(ctx).GetQuery(key)
	if !ok {
		valStr = "0"
	}
	return MustOne(strconv.Atoi(valStr))
}

func SetFilterObject[T any](ctx context.Context, filterObject T) {
	SetToContext(ConvertContext(ctx), filterObjectKey, filterObject)
}

func GetFilterObject[T any](ctx context.Context) (T, bool) {
	return GetFromContext[T](ConvertContext(ctx), filterObjectKey)
}

func MustGetFilterObject[T any](ctx context.Context) T {
	result, _ := GetFromContext[T](ConvertContext(ctx), filterObjectKey)
	return result
}

func ConvertContext(ctx context.Context) *gin.Context {
	return ctx.(*gin.Context)
}
