# Installation
1. Install golang
https://golang.org/doc/install
2. `go get -u github.com/vidmed/request`
3. `cd $GOPATH/src/github.com/vidmed/request/cmd`
4. `go build && ./cmd -config=config.toml` 
OR 
`go install && cmd -config=config.toml`

# Config file description
1. ListenAddr - address used by web server to listen to
2. ListenPort - port used by web server to listen to
3. LogLevel - logging level (panic = 0, fatal = 1, error = 2, warning = 3, info = 4, debug = 5)

### Testing app
After running app - two endpoints will be opened
- http://localhost:8899/request - random getting one request. You should see HTTP 200 Ok and a request in plain text ("aq" for example) as a response.
- http://localhost:8899/admin/requests - statistics for not 0 request views. If there is no views yet - you will see `views are empty - please call /request`. Please, open http://localhost:8899/request before doing admin requests.
When views are not empty - you should see HTTP 200 Ok and a request statistics in plain text.

To test you can call ab (apache benchmark):

`ab -c 8 -n 10000 "http://localhost:8899/request"`

### Stopping app
To stop application press `CTRL+C`. This will gracefully stop the server.