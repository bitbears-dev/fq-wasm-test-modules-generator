#!/usr/bin/env bash
set -Eeuo pipefail
d="$( cd "$( dirname "$0" )"; pwd -P )"

RESET=$'\e[0m'
BOLD=$'\e[1m'
GREEN=$'\e[0;32m'
RED=$'\e[0;31m'

tmpdir=$(mktemp -d -t fq-wasm-test.XXXXXXXX)
trap 'tear_down' 0

result=1
tear_down() {
    : "Clean up tmpdir" && {
        [[ $tmpdir ]] && rm -rf "$tmpdir"
    }

    : "Report result" && {
        if [ "$result" -eq 0 ]; then
            echo
            echo -e "${GREEN}${BOLD}OK${RESET}"
            echo
        else
            echo
            echo -e "${RED}${BOLD}FAILED${RESET}"
            echo
        fi
        exit $result
    }
}

if ! wat2wasm --version >/dev/null 2>&1; then
  echo
  echo "wat2wasm command is required"
  exit 1
fi

(
  cd "$d"
  go build -o "$tmpdir/test_data_generator"
)

(
  cd "$tmpdir"
  git clone --depth=1 --single-branch --branch=main https://github.com/WebAssembly/spec.git
)

mkdir -p "$d/output/"

find "$tmpdir" -name '*.wast' -print0 |
  while IFS= read -r -d '' f; do
    echo "$f"
    "$tmpdir/test_data_generator" -wat2wasm "wat2wasm" -input "$f" -output-dir "$d/output/"
  done

result=0
