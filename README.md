# Mongokaos
Simple mongodb data api-server for local development.

### Install
1. `git clone git@github.com:gummiboll/mongokaos.git`
2. `cd mongokaos`
3. `go mod tidy`
4. `go build`
5. Copy `.env.example` to `.env` and update it
6. `./mongokaos`
7. Test with: `curl -X POST http://localhost:<your port>/action/find --header 'Content-Type: application/json' --header 'api-key: <your api-key>' --data-raw '{ "dataSource": "something", "database": "something", "collection": "somecollection", "filter": { "something": "foo" } }'`
