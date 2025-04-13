#!/bin/bash

# Create output directory if it doesn't exist
mkdir -p auth-service/internal/infra/grpc/links/pb  

# Generate protobuf files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    auth-service/proto/links.proto 