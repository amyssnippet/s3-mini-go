BINARY_NAME=s3-mini
INSTALL_DIR=/usr/local/bin
ConfigDir=/etc/s3-mini
DataDir=/mnt/s3-data
.PHONY: all build test test-e2e test-all clean

all: build test test-e2e test-all clean

build:
	go build -o ${BINARY_NAME} main.go


install: build
	cp ${BINARY_NAME} ${INSTALL_DIR}
	mkdir -p ${ConfigDir}
	mkdir -p ${DataDir}
	id -u s3user &>/dev/null || useradd -r -s /bin/false s3user
	chown -R s3user:s3user ${ConfigDir}
	chown -R s3user:s3user ${DataDir}
	cp s3-mini.service /etc/systemd/system
	systemctl daemon-reload
	systemctl start s3-mini
	echo "Installation Complete"

test:
	@echo "Running Unit Tests..."
	go test -v ./internal/... ./cmd/...

test-e2e:
	@echo "Running End-to-End Tests..."
	go test -v ./test/e2e/...

test-all: test test-e2e

clean:
	go clean
	rm ${BINARY_NAME}
	rm -rf my_files/
