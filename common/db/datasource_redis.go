package db

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
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"time"
)

var redisCacheClient *redis.Client

func GetRedisCachePool() (*redis.Client, error) {
	if redisCacheClient != nil {
		return redisCacheClient, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.CoreConfig.Database.Redis.Host, config.CoreConfig.Database.Redis.Port),
		Password:     config.CoreConfig.Database.Redis.Password,
		MinIdleConns: config.CoreConfig.Database.Redis.MinConnPerDB,
		PoolSize:     config.CoreConfig.Database.Redis.MaxConnPerDB,
		DB:           config.CoreConfig.Database.Redis.CacheDbId,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	redisCacheClient = client

	closer.Bind(func() {
		_ = redisCacheClient.Close()
		log.Info("Redis cache client closed")
	})

	if log.GetLevel() == log.TraceLevel {
		go func() {
			for {
				stats := redisCacheClient.PoolStats()
				log.Tracef("Redis: Connections total: %d, stale: %d, idle: %d\n", stats.TotalConns, stats.StaleConns, stats.IdleConns)
				time.Sleep(5 * time.Second)
			}
		}()
	}

	return redisCacheClient, nil
}

var redisScheduleClient *redis.Client

func GetRedisSchedulePool() (*redis.Client, error) {
	if redisScheduleClient != nil {
		return redisScheduleClient, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.CoreConfig.Database.Redis.Host, config.CoreConfig.Database.Redis.Port),
		Password:     config.CoreConfig.Database.Redis.Password,
		MinIdleConns: config.CoreConfig.Database.Redis.MinConnPerDB,
		PoolSize:     config.CoreConfig.Database.Redis.MaxConnPerDB,
		DB:           config.CoreConfig.Database.Redis.ScheduleDbId,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	redisScheduleClient = client

	closer.Bind(func() {
		_ = redisScheduleClient.Close()
		log.Info("Redis schedule client closed")
	})

	if log.GetLevel() == log.TraceLevel {
		go func() {
			for {
				stats := redisScheduleClient.PoolStats()
				log.Tracef("Redis: Connections total: %d, stale: %d, idle: %d\n", stats.TotalConns, stats.StaleConns, stats.IdleConns)
				time.Sleep(5 * time.Second)
			}
		}()
	}

	return redisScheduleClient, nil
}
