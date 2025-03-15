# Simple Loan Engine

## Prerequisite:

Before running the app, make sure to have these programming language toolchains installed (and/or services running).

- [Golang](https://golang.org) : for the backend server app and API.

## Storage
- [PostgreSQL](https://www.postgresql.org) : for the main application database.

## Getting Started
Ensure you have already registered ssh to github account.
Clone this project, locate to your `GOPATH` folder by running this commands:

```
git clone git@github.com:zainulbr/simple-loan-engine.git #if using ssh

```

### How to running on local

On root directory of project
Ensure all modules downloaded
```
go mod tidy
```
Test codes 
```
go test -short $(go list ./...)
go test -race -short $(go list ./...)
```
Build app

```
cd cmd/loan-engine && go build
```
Rename `.env.exmaple` to `.env`

Edit the config to suit your environment `.env`

Run the app 

```
./loan-engine
```

or running without build
```
go run main.go
```

Build using `docker-compose`

```
docker-compose build dev    # build app
docker-compose up dev    # running app and redis db

```

### Detailed Desgin

[RFC](https://docs.google.com/document/d/1gVI51R14Tqdv2mEHdtK5J4iQZY0SBB8Z3jPomcUaSRw/edit?tab=t.0#heading=h.48g3fskl4cjo)
