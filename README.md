takonews-api
===========

[![Build Status](https://travis-ci.org/takonews/takonews-api.png?branch=master)](https://travis-ci.org/takonews/takonews-api)
[![codecov](https://codecov.io/gh/takonews/takonews-api/branch/master/graph/badge.svg)](https://codecov.io/gh/takonews/takonews-api)

## Development

```
go get -u github.com/takonews/takonews-api
mysql -u root -p
> create database if not exists takonews_development;
mv config/database.yml.sample config/database.yml
godep restore
echo "export GIN_MODE=debug" > ~/.bashrc
go run main.go
```

## Deployment

```
go get -u github.com/takonews/takonews-api
mysql -u root -p
> create database if not exists takonews_production;
mv config/database.yml.sample config/database.yml
godep restore
echo "export GIN_MODE=release" > ~/.bashrc
go run main.go
```

## Local testing

```
go test -v $(go list ./... | grep -v vendor)
```

## LICENSE

MIT.
