#!/bin/bash

Connections=("10" "100" "200" "500" "1000" "2000" "5000" "10000")

THREADS=4

DURATION_PER_TEST=30

BASE_URL="http://localhost:8080/"

TIMEOUT=5

FILE="asker_sim.lua"

# wrk -version exists with a goddamn error, whyyy
# wrk --version > /dev/null 2>&1
# if [ $? -ne 0 ]; then
#     echo "wrk is not installed. Please install wrk to run this script."
#     echo "You can install it by running:"
#     echo "sudo apt install wrk"
#     exit 1
# fi


# Function to run the wrk test
run_wrk_test() {
    local connection=$1
    local duration=$2

    echo "Running wrk test with $connection connections for $duration seconds..."

    wrk -t$THREADS -c$connection -d$duration --timeout $TIMEOUT -s $FILE $BASE_URL
}

main() {
    for connection in "${Connections[@]}"; do
        run_wrk_test "$connection" "${DURATION_PER_TEST}"
        echo "Test with $connection connections completed"
        echo "----------------------------------------"
    done
}

main