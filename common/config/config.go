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

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"time"
)

func Load(config interface{}) error {
	viperConfig := viper.New()
	viperConfig.AutomaticEnv()

	viperConfig.AddConfigPath("./config")
	viperConfig.AddConfigPath("../../../config")
	viperConfig.SetConfigName("application")
	viperConfig.SetConfigType("yaml")

	if err := viperConfig.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configuration file: %s\n", err.Error())
	}

	if profile := os.Getenv("PROFILE"); profile != "" {
		viperConfig.SetConfigName("application-" + strings.ToLower(profile))
		if err := viperConfig.MergeInConfig(); err != nil {
			log.WithError(err).Warn("Failed to merge the local configuration file")
		}

		// file watcher
		viperConfig.WatchConfig()
		viperConfig.OnConfigChange(func(event fsnotify.Event) { updateConfig(config, viperConfig) })
	}

	updateConfig(config, viperConfig)
	return nil
}

func updateConfig(config interface{}, viperConfig *viper.Viper) {
	// insert env instead placeholder
	for _, key := range viperConfig.AllKeys() {
		value := viperConfig.GetString(key)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			envName := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
			env := os.Getenv(envName)
			if env == "" {
				log.Panicf("Mandatory env variable not found: %s", envName)
			}
			viperConfig.Set(key, env)
		}
	}

	if err := viperConfig.Unmarshal(config, makeDecoderConfigOption); err != nil {
		log.WithError(err).Error("Failed to unmarshal config after update")
		return
	}

	if err := viperConfig.Unmarshal(&CoreConfig, makeDecoderConfigOption); err != nil {
		log.WithError(err).Error("Failed to unmarshal config after update")
		return
	}

	if level, err := log.ParseLevel(CoreConfig.Logging.Level); err == nil {
		log.SetLevel(level)
	}
	log.Info("Config file updated")
}

func makeDecoderConfigOption(decoderConfig *mapstructure.DecoderConfig) {
	decoderConfig.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		timeDecoderHookFunc,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
}

var timeDecoderHookFunc = func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	asString := data.(string)
	if asString == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, asString)
}
