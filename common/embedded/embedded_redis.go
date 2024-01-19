package embedded

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
	"assets/common/db"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateEmbeddedRedis() error {
	if os.Getenv("PROFILE") == "local" {
		containerPool, err := dockertest.NewPool("")
		if err != nil {
			return fmt.Errorf("could not connect to docker: %s", err)
		}

		var exposedPorts []string
		var portBindings map[docker.Port][]docker.PortBinding
		if !config.CoreConfig.Database.Redis.DynamicPort {
			exposedPorts = []string{"6379"}
			portBindings = map[docker.Port][]docker.PortBinding{
				"6379": {{HostIP: "0.0.0.0", HostPort: "6379"}},
			}
		}

		containerOptions := &dockertest.RunOptions{
			Name:         "redis",
			Repository:   "redis",
			Tag:          "6-alpine",
			ExposedPorts: exposedPorts,
			PortBindings: portBindings,
		}

		redisContainer, err := containerPool.RunWithOptions(containerOptions, getDockerHostConfig)
		if err != nil {
			return fmt.Errorf("could not create redis container: %s", err)
		}

		var containerExpirationSec = config.CoreConfig.Database.Redis.ContainerExpirationSec
		if containerExpirationSec != 0 {
			_ = redisContainer.Expire(containerExpirationSec)
		}

		closer.Bind(func() {
			_ = containerPool.Purge(redisContainer)
			log.Info("Redis container closed")
		})

		hostAndPort := strings.Split(redisContainer.GetHostPort("6379/tcp"), ":")
		config.CoreConfig.Database.Redis.Host = hostAndPort[0]
		config.CoreConfig.Database.Redis.Port, err = strconv.Atoi(hostAndPort[1])

		containerPool.MaxWait = time.Second * config.CoreConfig.Database.RetryMaxTimeSec
		if err = containerPool.Retry(func() error {
			_, err := db.GetRedisCachePool()
			return err
		}); err != nil {
			return fmt.Errorf("could not connect to redis: %s", err)
		}

		return nil
	}
	return nil
}
