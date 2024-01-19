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
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateEmbeddedJaeger() error {
	if os.Getenv("PROFILE") == "local" {
		containerPool, err := dockertest.NewPool("")
		if err != nil {
			return fmt.Errorf("could not connect to docker: %s", err)
		}

		exposedPorts := []string{"5775", "6831", "6832", "5778", "14268", "14250", "16686"}
		portBindings := map[docker.Port][]docker.PortBinding{
			"5775/udp":  {{HostIP: "0.0.0.0", HostPort: "5775"}},
			"6831/udp":  {{HostIP: "0.0.0.0", HostPort: "6831"}},
			"6379/udp":  {{HostIP: "0.0.0.0", HostPort: "6379"}},
			"5778/tcp":  {{HostIP: "0.0.0.0", HostPort: "5778"}},
			"14268/tcp": {{HostIP: "0.0.0.0", HostPort: "14268"}},
			"14250/tcp": {{HostIP: "0.0.0.0", HostPort: "14250"}},
			"16686/tcp": {{HostIP: "0.0.0.0", HostPort: "16686"}},
		}

		containerOptions := &dockertest.RunOptions{
			Name:         "jaeger",
			Repository:   "jaegertracing/all-in-one",
			Tag:          "1.47",
			ExposedPorts: exposedPorts,
			PortBindings: portBindings,
		}

		jaegerContainer, err := containerPool.RunWithOptions(containerOptions, getDockerHostConfig)
		if err != nil {
			return fmt.Errorf("could not create jaeger container: %s", err)
		}

		closer.Bind(func() {
			_ = containerPool.Purge(jaegerContainer)
			log.Info("Jaeger container closed")
		})

		hostAndPort := strings.Split(jaegerContainer.GetHostPort("16686/tcp"), ":")
		config.CoreConfig.Database.Jaeger.Host = hostAndPort[0]
		config.CoreConfig.Database.Jaeger.Port = hostAndPort[1]

		collectorHostAndPort := strings.Split(jaegerContainer.GetHostPort("14268/tcp"), ":")
		config.CoreConfig.Database.Jaeger.CollectorHost = collectorHostAndPort[0]
		config.CoreConfig.Database.Jaeger.CollectorPort = collectorHostAndPort[1]

		containerPool.MaxWait = time.Second * config.CoreConfig.Database.RetryMaxTimeSec
		if err = containerPool.Retry(func() error {
			_, err := http.Get(fmt.Sprintf("http://%s:%s/search", hostAndPort[0], hostAndPort[1]))
			return err
		}); err != nil {
			return fmt.Errorf("could not connect to jaeger: %s", err)
		}

		return nil
	}
	return nil
}
