all: policy.pb.go

policy.pb.go policy_grpc.pb.go: policy.proto
	protoc \
		--go_out=. --go_opt=paths=source_relative\
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		policy.proto
