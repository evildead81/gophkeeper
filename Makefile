APP_NAME=gophkeeper
OUTPUT_DIR=build
TARGETS=linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64

.PHONY: build clean

build:
	mkdir -p $(OUTPUT_DIR)
	@for target in $(TARGETS); do \
		GOOS=$${target%/*} GOARCH=$${target#*/} go build -o $(OUTPUT_DIR)/$(APP_NAME)_$${target%/*}_$${target#*/} cmd/server/main.go; \
		if [ "$${target%/*}" == "windows" ]; then \
			mv $(OUTPUT_DIR)/$(APP_NAME)_$${target%/*}_$${target#*/} $(OUTPUT_DIR)/$(APP_NAME)_$${target%/*}_$${target#*/}.exe; \
		fi; \
		echo "Собрано: $(OUTPUT_DIR)/$(APP_NAME)_$${target%/*}_$${target#*/}"; \
	done

clean:
	rm -rf $(OUTPUT_DIR)
