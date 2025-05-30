module github.com/abaika-abay/live_sports_project/api-gateway

go 1.23.4

require (
	github.com/abaika-abay/live_sports_project/match-service v0.0.0
	github.com/abaika-abay/live_sports_project/user-service v0.0.0
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.6
)

replace github.com/abaika-abay/live_sports_project/user-service => ../user-service

replace github.com/abaika-abay/live_sports_project/match-service => ../match-service

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)
