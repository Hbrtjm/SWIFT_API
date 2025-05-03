#!/bin/bash

# Verify curl is installed
if ! command -v curl &> /dev/null; then
    printf "Error: curl is not installed. Please install curl and try again.\n" >&2
    exit 1
fi

# Setup signal handler to gracefully exit on SIGTERM
trap 'echo "Received SIGTERM, exiting child process (PID: $$)"; exit 0' SIGTERM

# Base URL - replace with your actual API base URL
BASE_URL="http://localhost:8080/v1"

TEST_DATA='{
    "countryISO2": "DE",
    "swiftCode": "TESTDE45XXX",
    "bankName": "TEST GERMAN BANK",
    "address": "456 TEST STRASSE",
    "isHeadquarter": true,
    "countryName": "GERMANY"
}'

COUNTRY_CODES=("BG" "PL" "MT" "CL" "LV" "UY")

# Array of swift codes for endpoints 1 and 4
SWIFT_CODES=("AAISALTRXXX" "BPKOPLPWXXX" "BREXPLPWXXX" "BREXPLPWMBK" "BSCHCLR10R5")

# Function to make a GET request to endpoint 1
endpointGetSWIFT() {
    local swift_code=${SWIFT_CODES[$((RANDOM % ${#SWIFT_CODES[@]}))]}
    echo "Calling Endpoint 1: GET /v1/swift-codes/$swift_code"
    curl -s -X GET "$BASE_URL/swift-codes/$swift_code"
    echo
}

endpointGetCountry() {
    local country=${COUNTRY_CODES[$((RANDOM % ${#COUNTRY_CODES[@]}))]}
    echo "Calling Endpoint 2: GET /v1/swift-codes/country/$country"
    curl -s -X GET "$BASE_URL/swift-codes/country/$country"
    echo
}

endpointPost() {
    echo "Calling Endpoint 3: POST /v1/swift-codes"
    curl -s -X POST "$BASE_URL/swift-codes" \
        -H "Content-Type: application/json" \
        -d "$TEST_DATA"
    echo
}

# Function to make a DELETE request to endpoint 4
endpointDelete() {
    local swift_code=${SWIFT_CODES[$((RANDOM % ${#SWIFT_CODES[@]}))]}
    echo "Calling Endpoint 4: DELETE /v1/swift-codes/$swift_code"
    curl -s -X DELETE "$BASE_URL/swift-codes/$swift_code"
    echo
}

# Main function to randomly select and call an endpoint
ask() {
    local option=$((RANDOM % 4 + 1))
    
    case $option in
        1)
            endpointGetSWIFT
            ;;
        2)
            endpointGetCountry
            ;;
        3)
            endpointPost
            ;;
        4)
            endpointDelete
            ;;
        *)
            echo "Invalid option"
            ;;
    esac
}

# Run forever until killed
while true; do
    ask
    sleep $((RANDOM % 2 + 1))
done
