package config

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

import "time"

type DatabaseProperty struct {
	Postgres        PostgresProperty `yaml:"postgres,omitempty"`
	Redis           RedisProperty    `yaml:"redis,omitempty"`
	Jaeger          JaegerProperty   `yaml:"Jaeger,omitempty"`
	RetryMaxTimeSec time.Duration    `yaml:"retryMaxTimeSec,omitempty"`
	BatchSize       int              `yaml:"batchSize,omitempty"`
}

type PostgresProperty struct {
	Host           string
	Port           int
	Database       string
	Username       string
	Password       string
	ShowSql        bool
	MaxConnections int32
	MinConnections int32
}

type RedisProperty struct {
	Host                   string `yaml:"host,omitempty"`
	Port                   int    `yaml:"port,omitempty"`
	Username               string `yaml:"username,omitempty"`
	Password               string `yaml:"password,omitempty"`
	MinConnPerDB           int    `yaml:"minConnPerDB,omitempty"`
	MaxConnPerDB           int    `yaml:"maxConnPerDB,omitempty"`
	CacheDbId              int    `yaml:"cacheDbId,omitempty"`
	ScheduleDbId           int    `yaml:"scheduleDbId,omitempty"`
	ContainerExpirationSec uint   `yaml:"containerExpirationSec,omitempty"`
	DynamicPort            bool   `yaml:"dynamicPort,omitempty"`
}

type JaegerProperty struct {
	Host                   string `yaml:"host,omitempty"`
	Port                   string `yaml:"port,omitempty"`
	CollectorHost          string `yaml:"collectorHost,omitempty"`
	CollectorPort          string `yaml:"collectorPort,omitempty"`
	ContainerExpirationSec uint   `yaml:"containerExpirationSec,omitempty"`
	DynamicPort            bool   `yaml:"dynamicPort,omitempty"`
}
