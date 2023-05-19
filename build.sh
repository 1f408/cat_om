#!/bin/bash
set -eu
export GO111MODULE=on
OUT_DIR=$(pwd)/bin

echo "output: ${OUT_DIR}"
mkdir -p "${OUT_DIR}"
echo -n "install CATOM: "
(
  cd "src/catom" || exit 1
  go build -o "${OUT_DIR}" -ldflags '-s -w' || exit
) || continue
echo done
