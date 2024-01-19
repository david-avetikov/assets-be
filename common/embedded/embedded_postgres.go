package embedded

import (
	"assets/common/config"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateEmbeddedPostgres() error {
	if os.Getenv("PROFILE") == "local" {
		containerPool, err := dockertest.NewPool("")
		if err != nil {
			return fmt.Errorf("could not connect to docker: %s", err)
		}

		var exposedPorts []string
		var portBindings map[docker.Port][]docker.PortBinding
		port := strconv.Itoa(config.CoreConfig.Database.Postgres.Port)
		if port != "" && port != "0" {
			exposedPorts = []string{"5432"}
			portBindings = map[docker.Port][]docker.PortBinding{
				"5432": {{HostIP: "0.0.0.0", HostPort: port}},
			}
		}

		containerOptions := &dockertest.RunOptions{
			Name:         "postgres",
			Repository:   "postgres",
			Tag:          "12.1",
			ExposedPorts: exposedPorts,
			PortBindings: portBindings,
			Env: []string{
				fmt.Sprintf("POSTGRES_DB=%s", config.CoreConfig.Database.Postgres.Database),
				fmt.Sprintf("POSTGRES_USER=%s", config.CoreConfig.Database.Postgres.Username),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", config.CoreConfig.Database.Postgres.Password),
				"listen_addresses = '*'",
			},
		}

		postgresContainer, err := containerPool.RunWithOptions(containerOptions, getDockerHostConfig)
		if err != nil {
			return fmt.Errorf("could not create postgres container: %s", err)
		}

		closer.Bind(func() {
			_ = containerPool.Purge(postgresContainer)
			log.Info("Postgres container closed")
		})

		hostAndPort := strings.Split(postgresContainer.GetHostPort("5432/tcp"), ":")

		containerPool.MaxWait = time.Second * config.CoreConfig.Database.RetryMaxTimeSec
		if err = containerPool.Retry(func() error {
			uri := fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s",
				config.CoreConfig.Database.Postgres.Username,
				config.CoreConfig.Database.Postgres.Password,
				hostAndPort[0],
				hostAndPort[1],
				config.CoreConfig.Database.Postgres.Database,
			)
			dbConfig, err := pgx.ParseConfig(uri)
			if err != nil {
				log.WithError(err).Error("Couldn't parse postgres config")
				return nil
			}

			db := sqlx.NewDb(stdlib.OpenDB(*dbConfig), "pgx")
			if err := db.Ping(); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return fmt.Errorf("couldn't connect to postgres: %s", err)
		}

		return nil
	}
	return nil
}

func getDockerHostConfig(hostConfig *docker.HostConfig) {
	hostConfig.AutoRemove = false
	hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
}
