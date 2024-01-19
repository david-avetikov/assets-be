package cache

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
	"assets/common/db"
	"assets/common/util"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type Cache[T any] struct {
	Get   func(context context.Context) (T, error)
	Set   func(context context.Context, obj T) T
	Evict func(context context.Context)
}

func NewCache[T any](cacheName string, ttl time.Duration) *Cache[T] {
	return &Cache[T]{Get: GetCache[T](cacheName), Set: SetCache[T](cacheName, ttl), Evict: EvictCache(cacheName)}
}

func GetCache[T any](cacheName string) func(context context.Context) (T, error) {
	cache, err := db.GetRedisCachePool()
	if err != nil {
		log.WithError(err).Error("Couldn't get redis connection pool")
	}

	return func(context context.Context) (T, error) {
		var result T
		ctx := util.ConvertContext(context)

		if resultFromCache, err := cache.Get(ctx, cacheName+":"+ctx.Request.URL.String()).Result(); err == nil {
			return result, json.Unmarshal([]byte(resultFromCache), &result)
		}

		return result, custom_error.CacheNotFound
	}
}

func SetCache[T any](cacheName string, ttl time.Duration) func(context context.Context, obj T) T {
	cache, err := db.GetRedisCachePool()
	if err != nil {
		log.WithError(err).Error("Couldn't get redis connection pool")
	}

	return func(context context.Context, obj T) T {
		ctx := util.ConvertContext(context)
		log.WithContext(context).WithError(cache.Set(ctx, cacheName+":"+ctx.Request.URL.String(), string(util.MustOne(json.Marshal(obj))), ttl).Err())
		return obj
	}
}

func EvictCache(cacheNames ...string) func(context context.Context) {
	cache, err := db.GetRedisCachePool()
	if err != nil {
		log.WithError(err).Error("Couldn't get redis connection pool")
	}

	return func(context context.Context) {
		var allKeys []string

		for _, cacheName := range cacheNames {
			if keys, err := cache.Keys(context, cacheName+":*").Result(); err == nil {
				allKeys = append(allKeys, keys...)
			}
		}

		if len(allKeys) != 0 {
			cache.Del(context, allKeys...)
		}
	}
}
