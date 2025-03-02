#!/bin/bash

APP_NAME="gophkeeper"
OUTPUT_DIR="build"

mkdir -p $OUTPUT_DIR

TARGETS=("linux/amd64" "linux/arm64" "windows/amd64" "darwin/amd64" "darwin/arm64")

echo "Начинаем сборку для ОС..."
for TARGET in "${TARGETS[@]}"; do
    GOOS=${TARGET%/*}
    GOARCH=${TARGET#*/}

    OUTPUT_FILE="$OUTPUT_DIR/${APP_NAME}_${GOOS}_${GOARCH}"
    if [ "$GOOS" == "windows" ]; then
        OUTPUT_FILE+=".exe"
    fi

    echo "Сборка для $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT_FILE cmd/server/main.go
    if [ $? -eq 0 ]; then
        echo "Собрано: $OUTPUT_FILE"
    else
        echo "Ошибка сборки для $GOOS/$GOARCH"
    fi
done

echo "Сборка завершена! Бинарные файлы находятся в папке $OUTPUT_DIR"