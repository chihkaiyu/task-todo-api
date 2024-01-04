package docker

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"

	// init driver
	_ "github.com/lib/pq"
)

type repoInfo struct {
	repo          string
	tag           string
	containerName string
	port          int
	cmd           []string
	env           []string
	isReady       func(port string) error
	clearFunc     func(port string) error
}

var (
	reuseDocker       bool
	runningContainers = []*dockertest.Resource{}
	RepoTags          = map[string]*repoInfo{
		"postgres": {
			"postgres", "14.1-alpine", "LOCAL_POSTGRES", 5432,
			[]string{},
			[]string{"POSTGRES_HOST_AUTH_METHOD=trust"},
			checkPostgres, ClearPostgres,
		},
	}
)

func init() {
	_, reuse := os.LookupEnv("REUSE_DOCKER")
	reuseDocker = reuse
}

func RunExternal(repos []string) ([]string, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error().Err(err).Msg("dockertest.NewPool failed")
		return []string{}, err
	}

	ports := []string{}
	for _, r := range repos {
		info, ok := RepoTags[r]
		if !ok {
			log.Warn().Str("repo", r).Msg("unknown repo")
			return []string{}, fmt.Errorf("unknown repo")
		}

		if reuseDocker {
			c, found := pool.ContainerByName(info.containerName)
			if found {
				// found container but not running and we need to remove it
				if !strings.Contains(c.Container.State.String(), "Up") {
					if err := pool.Purge(c); err != nil {
						log.Error().Err(err).Str("repo", c.Container.ID).Msg("poolPurge failed")
						return []string{}, err
					}
					port, err := run(pool, info)
					if err != nil {
						log.Error().Err(err).Msg("run failed")
						return []string{}, err
					}
					ports = append(ports, port)
				}
				port := c.GetPort(fmt.Sprintf("%d/tcp", info.port))
				err := info.clearFunc(port)
				if err == nil {
					ports = append(ports, port)
					continue
				}
			}
		}

		port, err := run(pool, info)
		if err != nil {
			log.Error().Err(err).Msg("run failed")
			return []string{}, err
		}
		ports = append(ports, port)
	}

	return ports, nil
}

func run(pool *dockertest.Pool, info *repoInfo) (string, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       info.containerName,
		Repository: info.repo,
		Tag:        info.tag,
		Env:        info.env,
		Cmd:        info.cmd,
	})
	if err != nil {
		log.Error().Err(err).Msg("pool.Run failed")
		return "", fmt.Errorf("repo run failed: %s", info.repo)
	}

	var port string
	if err := pool.Retry(func() error {
		port = resource.GetPort(fmt.Sprintf("%d/tcp", info.port))
		if err := info.isReady(port); err != nil {
			return err
		}
		runningContainers = append(runningContainers, resource)
		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("could not connect to docker")
		return "", fmt.Errorf("repo initialize failed: %s", info.repo)
	}

	return port, nil
}

func RemoveExternal() error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error().Err(err).Msg("could not connect to docker")
		return err
	}

	if reuseDocker {
		return nil
	}

	for _, r := range runningContainers {
		if err := pool.Purge(r); err != nil {
			log.Error().Err(err).Str("repo", r.Container.ID).Msg("pool.Purge failed")
			continue
		}
	}

	return nil
}

func checkPostgres(port string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres@localhost:%s/?sslmode=disable", port))
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// ClearPostgres drops all databases
func ClearPostgres(port string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres@localhost:%s/?sslmode=disable", port))
	if err != nil {
		return err
	}
	defer db.Close()

	rows := []string{"gogolook"}
	for _, row := range rows {
		if _, err := db.Exec("DROP DATABASE IF EXISTS " + row + " WITH (FORCE)"); err != nil {
			log.Error().Err(err).Msg("drop database failed")
		}
	}

	return nil
}
