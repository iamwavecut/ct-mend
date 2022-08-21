
[![Go Reference](https://pkg.go.dev/badge/github.com/iamwavecut/ct-mend.svg)](https://pkg.go.dev/github.com/iamwavecut/ct-mend)
[![go.mod](https://img.shields.io/github/go-mod/go-version/iamwavecut/ct-mend)](go.mod)
[![Build Status](https://img.shields.io/github/workflow/status/iamwavecut/ct-mend/build)](https://github.com/iamwavecut/ct-mend/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamwavecut/ct-mend)](https://goreportcard.com/report/github.com/iamwavecut/ct-mend)
[![Codecov](https://codecov.io/gh/iamwavecut/ct-mend/branch/main/graph/badge.svg)](https://codecov.io/gh/iamwavecut/ct-mend)

Mend Code Test
===

by Valeriy Selitskiy

<details>
<summmary>The problem</summmary>

>#### Build a small application exposing CRUD endpoints through a REST API:
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