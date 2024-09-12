# URL Shortener

## About this Repo
This is a repository for URL shortener apis. It comes with most sensible plugins already installed:

- [x] [gin](https://gin-gonic.com)(Framework) for handling requests
- [x] [postgresql](https://www.postgresql.org/) for relational database
- [x] gorm ([go-gorm/gorm](https://github.com/go-gorm/gorm)) for ORM library for Golang aims to be developer friendly

##  How to run

Update the _default.yaml_ file

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
