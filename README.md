# Task Todo API

# Prerequisite
- Docker, docker-compose (must)
- Golang (must)
- (sql-migrate)[https://github.com/rubenv/sql-migrate] (optinal)
- (swaggo)[https://github.com/swaggo/swag] (optional)

# Features
- API Document automatically created from code  
- Basic metrics included (e.g. golang performance, API response time)
- All tools package in docker

# How to Start
1. Build `swaggo` image (you can skip if you have installed it local)
```shell
make build-swaggo
```

2. Build sql-migrate image (you can skip if you have install it local)
```shell
make build-migrate
```

3. Build API server image
```shell
make build-api
```

4. Spin up cluster
```shell
docker-compose up -d
```

5. Migrate database
```shell
make migrate
```

# Test
Run (`export REUSE_DOCKER=1` so that you can reuse the docker container without shutdown it every time you run test)
```shell
make test
```

