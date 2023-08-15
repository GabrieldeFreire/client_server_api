SERVER_PATH := server/server.go
CLIENT_PATH := client/client.go

SERVER_LOG := server.log
SERVER_PID := server_pid.txt


run-server:
	@ nohup go run server/server.go > $(SERVER_LOG) 2>&1 & echo $$! > $(SERVER_PID)

run-client:
	go run $(CLIENT_PATH)

shutdown-server:
	@ pkill -P $(shell cat $(SERVER_PID))


