BUILD_DIR="build"
LINUX_ARGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DBS_PREFIX="distributed-broad-system"
WINDOWS_ARGS=CGO_ENABLED=0 GOOS=windows GOARCH=amd64
VERSION=1.0

all:
	rm -r build
	@go build -o ${BUILD_DIR}//macos/${DBS_PREFIX}server-${VERSION} cmd/main.go
	@${LINUX_ARGS} go build -o ${BUILD_DIR}//linux/${DBS_PREFIX}server-${VERSION} cmd/main.go
	@${WINDOWS_ARGS} go build -o ${BUILD_DIR}//windows/${DBS_PREFIX}server-${VERSION}.exe cmd/main.go
	@echo distributed_broad_system binary built.

clean:
	rm -r build
	@echo Products cleaned.