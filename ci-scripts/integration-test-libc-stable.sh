#!/bin/bash
# Runs "stable"-mode tests against a skycoin node configured with a pinned database
# "stable" mode tests assume the blockchain data is static, in order to check API responses more precisely
# $TEST defines which test to run i.e, cli or gui; If empty both are run
 
#Set Script Name variable
SCRIPT=`basename ${BASH_SOURCE[0]}`
PORT="46420"
RPC_PORT="$PORT"
HOST="http://127.0.0.1:$PORT"
RPC_ADDR="http://127.0.0.1:$RPC_PORT"
MODE="stable"
BINARY="skycoin-integration"
TEST=""
UPDATE=""
# run go test with -v flag
VERBOSE=""
# run go test with -run flag
RUN_TESTS=""
# run tests with csrf enabled
USE_CSRF=""
DISABLE_CSRF="-disable-csrf"

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

usage () {
  echo "Usage: $SCRIPT"
  echo "Optional command line arguments"
  echo "-t <string>  -- Test to run, gui or cli; empty runs both tests"
  echo "-r <string>  -- Run test with -run flag"
  echo "-u <boolean> -- Update stable testdata"
  echo "-v <boolean> -- Run test with -v flag"
  echo "-c <boolean> -- Run tests with CSRF enabled"
  exit 1
}

while getopts "h?t:r:uvc" args; do
  case $args in
    h|\?)
        usage;
        exit;;
    t ) TEST=${OPTARG};;
    r ) RUN_TESTS="-run ${OPTARG}";;
    u ) UPDATE="--update";;
    v ) VERBOSE="-v";;
    c ) USE_CSRF="1"; DISABLE_CSRF="";
  esac
done

set -euxo pipefail

DATA_DIR=$(mktemp -d -t skycoin-data-dir.XXXXXX)
WALLET_DIR="${DATA_DIR}/wallets"

if [[ ! "$DATA_DIR" ]]; then
  echo "Could not create temp dir"
  exit 1
fi

# Compile the skycoin node
# We can't use "go run" because this creates two processes which doesn't allow us to kill it at the end
echo "compiling skycoin"
go build -o "$BINARY" -ldflags "${GOLDFLAGS}" cmd/skycoin/skycoin.go

# Run skycoin node with pinned blockchain database
echo "starting skycoin node in background with http listener on $HOST"

./skycoin-integration -disable-networking=true \
                      -web-interface-port=$PORT \
                      -download-peerlist=false \
                      -db-path=./src/api/integration/testdata/blockchain-180.db \
                      -db-read-only=true \
                      -rpc-interface=true \
                      -launch-browser=false \
                      -data-dir="$DATA_DIR" \
                      -enable-wallet-api=true \
                      -wallet-dir="$WALLET_DIR" \
                      $DISABLE_CSRF \
                      -enable-seed-api=true &
SKYCOIN_PID=$!

echo "skycoin node pid=$SKYCOIN_PID"

echo "sleeping for startup"
sleep 3
echo "done sleeping"

set +e


LD_LIBRARY_PATH="/usr/local/lib:build/usr/lib:build/libskycoin" ./bin/test_int_libskycoin_shared
FAIL_STATIC_SHARED=$?

LD_LIBRARY_PATH="/usr/local/lib:build/usr/lib:build/libskycoin" ./bin/test_int_libskycoin_static
FAIL_STATIC=$?

echo "shutting down skycoin node"

# Shutdown skycoin node
kill -s SIGINT $SKYCOIN_PID
wait $SKYCOIN_PID

rm "$BINARY"

if [ $FAIL_STATIC_SHARED -ne 0 ]; then
	exit $FAIL_STATIC_SHARED
elif [ $FAIL_STATIC -ne 0 ]; then
	exit $FAIL_STATIC
else
	exit 0
fi