# REDIS_URL as a global env screws with redis Ruby gem.
go build ./... && REDIS_URL="localhost:6379" ./platform-ws-services
