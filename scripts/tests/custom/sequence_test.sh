#!/bin/bash

# Check if curl is installed
curl -V > /dev/null 2>&1
if [ $? -ne 0 ] ; then 
    printf "curl is not installed\n"
    printf "install curl with:\n"
    printf "\tsudo apt install curl\n"
    exit 1
fi

printf "Running deployment script...\n"

# This will re-deploy the application, so it's only used for initial test
chmod +x ../../deployment/deploy.sh
../../deployment/deploy.sh
if [ $? -ne 0 ] ; then
    printf "Deployment failed\n"
    exit 2
fi

while ! curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health | grep -q 200; do
    printf "Waiting for the application to start...\n"
    sleep 5
done

printf "Application is up and running\n"

cleanup()
{
    # Clean up any existing containers from previous runs
    echo "Cleaning up existing containers..."
    docker compose down -v --remove-orphans 2>/dev/null || true

    # If it was built as a part of the test, remove any test containers taking up the port 27017
    containers=$(docker ps -a -q --filter "name=mongo" 2>/dev/null) || true
    if [ -n "$containers" ]; then
      echo "Stopping and removing existing MongoDB containers..."
      docker stop $containers || true
      docker rm $containers || true
    fi
    printf "Cleanup complete\n"
    printf "Running containers:\n"
    docker ps
}

trap cleanup EXIT

FAILED=0
PASSED=0
NORESPONSE=0

printf "Posting the branch swift code...\n"

# Curl request from postman
# Post the branch
curl -s --location 'http://localhost:8080/v1/swift-codes' \
--header 'Content-Type: application/json' \
--data '{
	"address":       "123 TEST BOULEVARD",
	"bankName":      "TEST FRENCH BANK",
	"countryISO2":   "FR",
	"countryName":   "FRANCE",
	"isHeadquarter": false,
	"swiftCode":     "TESTFR22FDF"
} 
' > post_branch_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    # Check if the response is as expected
    expected_response='{"message":"Bank with SWIFT code created successfully"}'
    response=$(cat post_branch_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully posted the branch swift code\n"
    fi
fi

printf "Posting the headquarter swift code...\n"

# Post the headquarter
curl -s --location 'http://localhost:8080/v1/swift-codes' \
--header 'Content-Type: application/json' \
--data '{
	"address":       "123 TEST BOULEVARD",
	"bankName":      "TEST FRENCH BANK",
	"countryISO2":   "FR",
	"countryName":   "FRANCE",
	"isHeadquarter": true,
	"swiftCode":     "TESTFR22XXX"
}
' > post_headquarter_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    # Check if the response is as expected
    expected_response='{"message":"Bank with SWIFT code created successfully"}'
    response=$(cat post_headquarter_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully posted the headquarter swift code\n"
    fi
fi

printf "Getting the branch swift code...\n"

# Get the branch swift code
curl -s --location 'http://localhost:8080/v1/swift-codes/TESTFR22FDF' \
--header 'Content-Type: application/json' > get_branch_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response='{"address":"123 TEST BOULEVARD","bankName":"TEST FRENCH BANK","countryISO2":"FR","countryName":"FRANCE","isHeadquarter":false,"swiftCode":"TESTFR22FDF"}'
    response=$(cat get_branch_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully got the branch swift code\n"
    fi
fi

printf "Getting the headquarter swift code...\n"

# Get headquarter swift code
curl -s --location 'http://localhost:8080/v1/swift-codes/TESTFR22XXX' \
--header 'Content-Type: application/json' > get_headquarter_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response='{"address":"123 TEST BOULEVARD","bankName":"TEST FRENCH BANK","countryISO2":"FR","countryName":"FRANCE","isHeadquarter":true,"swiftCode":"TESTFR22XXX","branches":[{"address":"123 TEST BOULEVARD","bankName":"TEST FRENCH BANK","countryISO2":"FR","isHeadquarter":false,"swiftCode":"TESTFR22FDF"}]}'
    response=$(cat get_headquarter_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully got the headquarter swift code\n"
    fi
fi

printf "Getting the swift codes by country ISO2 code (MT - Malta)...\n"

curl -s --location --request GET 'http://localhost:8080/v1/swift-codes/country/MT' \
--header 'Content-Type: application/json' > get_country_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response=$(cat expected_country_test.json)
    response=$(cat get_country_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully got the swift codes by country ISO2 code (MT - Malta)\n"
    fi
fi

printf "Deleting the branch swift code TESTFR22FDF...\n"

curl -s --location --request DELETE 'http://localhost:8080/v1/swift-codes/TESTFR22FDF' \
--header 'Content-Type: application/json' > delete_branch_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response='{"message":"SWIFT code deleted successfully"}'
    response=$(cat delete_branch_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully deleted the branch swift code TESTFR22FDF\n"
    fi
fi

printf "Checking if the branch swift code TESTFR22FDF is deleted...\n"

curl -s --location 'http://localhost:8080/v1/swift-codes/TESTFR22FDF' \
--header 'Content-Type: application/json' > get_branch_error_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response='{"message":"no bank found with the given SWIFT code"}'
    response=$(cat get_branch_error_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully checked if the branch swift code TESTFR22FDF is deleted\n"
    fi
fi

printf "Deleting the headquarter swift code TESTFR22XXX...\n"

curl -s --location --request DELETE 'http://localhost:8080/v1/swift-codes/TESTFR22XXX' \
--header 'Content-Type: application/json' > delete_headquarter_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response='{"message":"SWIFT code deleted successfully"}'
    response=$(cat delete_headquarter_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully deleted the headquarter swift code TESTFR22XXX\n"
    fi
fi

printf "Checking if the headquarter swift code TESTFR22XXX is deleted...\n"

curl -s --location 'http://localhost:8080/v1/swift-codes/TESTFR22XXX' \
--header 'Content-Type: application/json' > get_headquarter_error_response.json

if [ $? -ne 0 ] ; then
    printf "Request failed\n"
    NORESPONSE=$((NORESPONSE + 1))
else
    expected_response='{"message":"no bank found with the given SWIFT code"}'
    response=$(cat get_headquarter_error_response.json)
    if [ "$response" != "$expected_response" ]; then
        printf "Response does not match expected response\n"
        printf "Expected: $expected_response\n"
        printf "Got: $response\n"
        FAILED=$((FAILED + 1))
    else
        PASSED=$((PASSED + 1))
        printf "Successfully checked if the headquarter swift code TESTFR22XXX is deleted\n"
    fi
fi


printf "Testing complete\n"
printf "Number of tests failed: $FAILED\n"
printf "Number of tests passed: $PASSED\n"
printf "Tests without response: $NORESPONSE\n"
printf "Total tests executed: $((FAILED + PASSED + NORESPONSE))\n"

# Exit with error code if any tests failed or had no response
if [ $FAILED -gt 0 ] || [ $NORESPONSE -gt 0 ]; then
    exit 1
else
    exit 0
fi