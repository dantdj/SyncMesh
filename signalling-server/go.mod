module github.com/dantdj/syncmesh/signalling-server

go 1.25.1

require (
	github.com/dantdj/syncmesh/api v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
)

replace github.com/dantdj/syncmesh/api => ../api
