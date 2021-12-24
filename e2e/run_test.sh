#!/bin/bash

set -euo pipefail

GO111MODULE=on
cd gql-server
rm -f test.json loginAPI.json
(go get && go run server.go &)
sleep 2
cd ..
go test -i
go test

exit_status=$?

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

if [ $exit_status -ne 0 ]; then
  exit_status=1
  echo -e $"\n${RED}TESTING FAILED. Review test logs for details.${NC}\n"
else
  echo -e $"\n${GREEN}TESTING COMPLETED SUCCESSFULLY${NC}\n"
fi

echo "Stopping graphql test server"
kill $(lsof -t -i:8080)

exit $exit_status