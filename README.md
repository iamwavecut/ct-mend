[![go.mod](https://img.shields.io/github/go-mod/go-version/iamwavecut/ct-mend)](go.mod)
[![Build Status](https://img.shields.io/github/workflow/status/iamwavecut/ct-mend/build)](https://github.com/iamwavecut/ct-mend/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamwavecut/ct-mend)](https://goreportcard.com/report/github.com/iamwavecut/ct-mend)
[![Codecov](https://codecov.io/gh/iamwavecut/ct-mend/branch/main/graph/badge.svg)](https://codecov.io/gh/iamwavecut/ct-mend)

Mend Code Test
===

by Valeriy Selitskiy

The problem
---
<details>
#### Build a small application exposing CRUD endpoints through a REST API:

>- It can be any kind of data/object that can be created, retrieved, updated and deleted on endpoints
   accepting POST, GET, PUT and DELETE methods.
>- Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB,
   or something else. The business rules are not bound to the database. Provide an implementation
   of SQL and NOSQL (any kind) for the above CRUD operations.
>- Provide tests. The API can be tested without the Database or any other external element.
>- Dockerize the solution, e.g. provide a Dockerfile to run it.
>- The API server should listen on TLS only.
>#### Evaluation Criteria
>- Please push your code to a GitHub repository or send us an archive.
>- Include a Readme helping us run your service.
>- Include a section about your thought process explaining your choices and share other alternative
   designs you considered.
>- Add any information you deem interesting for us to better understand your assignment.
>#### We'll evaluate
>- The readability of your code (including readability of your tests)
>- The correctness of the outputs of the API in accordance with REST standard.
>- Your ability to share your design choices and clearly weight pros & cons of alternative solutions
>- The overall quality of your written technical communication
</details>

The solution
---
### TL;DR
```shell
git clone https://github.com/iamwavecut/ct-mend && cd ct-mend
```
Then
```shell
# build with SQLite storage
# docker compose build server --build-arg STORAGE_TYPE=sqlite STORAGE_ADDR=./db.sqlite # or simply
docker compose build server
docker compose up server # -d
```
_OR_
```shell
# build with MongoDB storage
docker compose build --build-arg STORAGE_TYPE=mongodb --build-arg STORAGE_ADDR=mongodb://mend:mend@mongo
docker compose up # -d
```
And finally, to perform a bunch of requests and to see funky log output
```shell
CLIENT_HOST=127.0.0.1:8443 make run-client
```
#### Configuration options
Server configuration can be set

- at build time for the docker image as `--build-arg`'s
- at run time for the local binaries as `ENV` vars.

##### Client options:
|Name|Default|Comment|
|:---|:---:|---|
|**Server**|||
|`TLS_ADDR`|`:8443`|Exposed on all interfaces by default to avoid routing issues of your environment|
|`STORAGE_TYPE`|`sqlite`|Options: <br> - `sqlite` <br> - `mongodb` |
|`STORAGE_ADDR`|`./db.sqlite`|Path of physical location of db file, or URL of MongoDB instance |
|`LOG_LEVEL`|`trace`|Options: <br>- `trace`<br>- `debug`<br>- `info`<br>- `warning`<br>- `error`<br>- `fatal`<br>- `panic`|
|`GRACEFUL_TIMEOUT`|`10s`| To specify default timeout of connections |
|**Client**|||
|`CLIENT_HOST`|`127.0.0.1:8443`|Can be used to override HTTP client target in case of remote server deployment |

#### Local fun
As fast as
```shell
make build
nohup server & # outputs e.g. [1] 12345
make run-client
kill -9 12345 # or `fg`⏎, Ctrl+C
```
To list all `Makefile` targets, call `make help`:
```shell
build                          build server binary
dev                            generate vet fmt lint test mod-tidy
fmt                            go fmt
generate                       go generate and OPENSSL keys generation
go-clean                       go clean build, test and modules caches
lint                           golangci-lint
mod-tidy                       go mod tidy
run-client                     run client which tries all the api endpoints and prints log output
vet                            go vet
```

#### Boring details, worth mentioning
To not deal with too much or too little boilerplate code I decided to try scaffolding my project from [golang-templates/seed](https://github.com/golang-templates/seed). It's a good starting point for small projects.

I've created a `storage.Adapter` interface to abstract business from storage, two adapters included:
- SQLite v3 _(default)_
- MongoDB

HTTP Server is listening on `:8443` by default. TLS certificates are generated during `make generate` and, of course, on the docker container build and getting embedded into binary to not be easily accessible in the container.

##### Testing
To be honest, there is not much to test on the go side, because it relies on `stdlib` and well-tested 3'rd party components. MongoDB driver is hard to test, in particular, as every move must be mocked and it's a lot of work, comparable with a whole code test in efforts. The only reasonable unit test I made is the HTTP handlers test, which may be run by `make test` and yet, it lacks negative test cases.

Instead, I made integration testing easier. I created another executable - `client`. In fact, its parsing [resources/api.http](resources/api.http) and executing every request found, printing the result to log.
As an alternative, you may use Jetbrains HTTP client to test API using provided [resources/api.http](resources/api.http) file. It also comes with tests bundled.
![](https://i.imgur.com/DMOdeLX.png)
