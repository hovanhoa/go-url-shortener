# URL Shortener

## About this Repo
This is a repository for URL shortener apis. It comes with most sensible plugins already installed:

- [x] [gin](https://gin-gonic.com)(Framework) for handling requests
- [x] [postgresql](https://www.postgresql.org/) for relational database
- [x] gorm ([go-gorm/gorm](https://github.com/go-gorm/gorm)) for ORM library for Golang aims to be developer friendly
- [x] redis ([redis/go-redis](https://github.com/redis/go-redis)) for caching query
- [x] rate limiter ([internal\middleware\ratelimit](\internal\middleware\ratelimit)) middleware for prevent too many requests

##  How to run

Create the _default.yaml_ file by copy _example.yaml_ file in deployments folder

Download dependencies
```shell
go mod download
```

Run server
```shell
make server
```

or

Run with docker
```shell
cd deployments
POSTGRES_PASSWORD=... REDIS_PASSWORD=... docker-compose up -d
```

Open [Makefile](Makefile) for more details.

## License
Distributed under the MIT License.
