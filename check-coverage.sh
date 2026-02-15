#!/bin/bash
go test -v -coverpkg=./src/... -coverprofile=coverage.out ./tests >> test_results.log

# Exclude files from coverage
grep -vE "src/cmd/*" coverage.out > coverage.temp
mv coverage.temp coverage.out

if [[ "$1" == "--sort" ]]; then
    go tool cover -func=coverage.out | grep -v "total:" | sort -k3 -n
    echo "--------------------------------------------------------"
    go tool cover -func=coverage.out | grep "total:"
else
    go tool cover -func=coverage.out
fi
