# SWIFT_API

[![Test & Deployment](https://github.com/Hbrtjm/SWIFT_API/actions/workflows/cicd.yml/badge.svg)](https://github.com/Hbrtjm/SWIFT_API/actions/workflows/cicd.yml)

A CRUD service for parsing, storing, and exposing international bank SWIFT (BIC) codes via a RESTful API.

## About the project

SWIFT_API ingests SWIFT code data from a CSV file, persists it in a database (MongoDB by default), and provides endpoints to create, read, update, and delete SWIFT entries at both branch and headquarters levels. It is built with Go and Docker for fast bootstrapping and deployment.

## Requirements

International wire transfers rely on SWIFT codes (BICs) to route payments to the correct bank and branch. This service:

- **Parses** SWIFT code entries from a CSV, identifying headquarters (`XXX` suffix) and branches (matching first 8 characters).
- **Normalizes** country ISO codes and names (uppercased).
- **Stores** entries in MongoDB, optimized for low-latency queries by code or country.
- **Exposes** a versioned REST API for retrieval, creation, and deletion of SWIFT entries.

## Required endpoints

### 1. Retrieve SWIFT Code Details

Fetches detailed information for a specific SWIFT code, including all branch information for headquarters.

```
GET /v1/swift-codes/{swift-code}
```

#### Response Structure (Headquarters)

```json
{
    "address": string,
    "bankName": string,
    "countryISO2": string,
    "countryName": string,
    "isHeadquarter": bool,
    "swiftCode": string,
    "branches": [
        {
            "address": string,
            "bankName": string,
            "countryISO2": string,
            "isHeadquarter": bool,
            "swiftCode": string
        },
        {
            "address": string,
            "bankName": string,
            "countryISO2": string,
            "isHeadquarter": bool,
            "swiftCode": string
        },
        ...
    ]
}
```

#### Response Structure (Branch)

```json
{
    "address": string,
    "bankName": string,
    "countryISO2": string,
    "countryName": string,
    "isHeadquarter": bool,
    "swiftCode": string
}
```

### 2. List SWIFT Codes by Country

Returns all SWIFT codes (both headquarters and branches) for a specific country.

```
GET /v1/swift-codes/country/{countryISO2code}
```

#### Response Structure

```json
{
    "countryISO2": string,
    "countryName": string,
    "swiftCodes": [
        {
            "address": string,
            "bankName": string,
            "countryISO2": string,
            "isHeadquarter": bool,
            "swiftCode": string
        },
        {
            "address": string,
            "bankName": string,
            "countryISO2": string,
            "isHeadquarter": bool,
            "swiftCode": string
        },
        ...
    ]
}
```

### 3. Add New SWIFT Code

Creates a new SWIFT code entry in the database.

```
POST /v1/swift-codes
```

#### Request Structure

```json
{
    "address": string,
    "bankName": string,
    "countryISO2": string,
    "countryName": string,
    "isHeadquarter": bool,
    "swiftCode": string
}
```

#### Response Structure

```json
{
    "message": string
}
```

### 4. Delete SWIFT Code

Removes a SWIFT code entry from the database.

```
DELETE /v1/swift-codes/{swift-code}
```

#### Response Structure

```json
{
    "message": string
}
```
## Setup and deploy

### Linux or WSL

#### Using 'deploy' script

I have prepared a launch script to make the deployment easier. This script however doesn't include tests, they should be run in the repository, but I don't guarantee that after changes the code will run properly. The code itself is provided "as is" and except for the test coverage I cannot forsee the future of deprecations or possible vunurabilities in the upcoming versions of GoLang. 

To run the code navigate from the main `SWIFT_API/` directory to the `scripts/deployment` directory:

```bash
cd ./SWIFT_API/scripts/deployment
```

Then to make sure the code runs, execute:

```bash
# This adds permissions to all of the UNIX users, so anyone having access to the file can run it  
chmod +x deploy.sh
./deploy.sh
```

The deployment should start automatically, you can verify that the container is running using:

```bash
docker ps -a
```

You should see the containers of names `mongo:6`, `api:v1`, `grafana/promtail`, `grafana/loki`, and `grafana/grafana`. You can also check that they indeed operate on the same network provided by the `docker compose` configuration:

```bash
docker network list
# Check for the network of the name 'swift_api' or 'swift_api_default' 
docker network inspect swift_api # or swift_api_default
```
You should see the containers that were registered in the network.

#### Manual launch

If you prefer to run the application manually, follow these steps:

1. Navigate to the project root directory:
   ```bash
   cd ./SWIFT_API
   ```

2. Build and start the containers using docker-compose:
   ```bash
   docker-compose build
   docker-compose up -d
   ```

3. To stop the application:
   ```bash
   docker-compose down
   ```

4. To view logs (not recommended, however, unreadable, use logs in a particular container or docker desktop):
   ```bash
   docker-compose logs -f api
   ```

5. View api logs (you will see nothing if the SPEEDUP_MODE is enabled and wihout it the server still logs into file):
    ```bash
    docker logs swift_api-api-1
    ```

### Windows 10/11

#### Using 'deploy' script

You can run the deployment scripts by navigating to the specific folder and executing one of the provided scripts.

Navigate to the scripts directory:

```powershell
cd .\SWIFT_API\scripts\deployment
```

Then execute one of the following scripts:
- For PowerShell: `.\deploy.ps1`
- For Command Prompt: `deploy.bat`

The script will handle the building and deployment of all necessary containers.

#### Manual launch

For manual deployment on Windows:

1. Navigate to the project root directory:
   ```powershell
   cd .\SWIFT_API
   ```

2. Build and start the containers using docker-compose:
   ```powershell
   docker-compose build
   docker-compose up -d
   ```

3. To stop the application:
   ```powershell
   docker-compose down
   ```

### Running form source

Make sure to have go installed and GOPATH set up. Then execute:

```bash
cd ./backend
go mod download
go build -o main ./cmd/api/main.go
chmod +x main
./main
```

Keep in mind this requires MongoDB to be accessible on mongodb://localhost:27017 or mongodb://mongo:27017

## Environment Configuration

The application uses several environment variables that can be configured in the `docker-compose.yaml` file:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| MONGO_URI | MongoDB connection string | mongodb://localhost:27017 |
| DB_NAME | Database name | swiftcodes |
| LOGGER_PREFIX | Prefix for log entries | api |
| LOGGER_DEBUG | Enable detailed logging | false |
| API_DEBUG | Enable API debug responses with full error messages | false |
| BANKS_COLLECTION_NAME | MongoDB collection for banks data | banks |
| COUNTRIES_COLLECTION_NAME | MongoDB collection for countries data | countries |
| LOAD_INITIAL_DATA | Flag to load initial data into the database | true |
| SWIFT_DATA_FILE | Path to the initial data CSV file | configs/swift_data.csv |
| VERSION | API version (used in URL paths) | v1 |
| SPEEDUP_MODE | Discard logs to improve performance | false |

## Logging and Monitoring

The project includes a full observability stack with:

1. **Promtail**: Log collector that reads logs from the API service
2. **Loki**: Log aggregation system
3. **Grafana**: Visualization platform for logs and metrics

To access Grafana:
1. Open your browser and navigate to `http://localhost:3000`
2. Login with:
   - Username: `admin`
   - Password: `admin` (Change this in production environments !!!!)

## API Usage

Once the application is running, the API will be available at `http://localhost:8080`.

The API provides endpoints for managing SWIFT codes and bank information. The base URL includes the version as specified in the environment variables (default: `/v1/`).

## Testing

The project includes various test files for unit and integration testing. To run tests:

```bash
cd ./SWIFT_API/backend
go test ./... -v
```

For custom stress testing, use the scripts in the `scripts/tests/custom` directory:

```bash
cd ./SWIFT_API/scripts/tests/custom
chmod +x stress_test.sh
./stress_test.sh
```

To test all the enpoints run:

```bash
cd ./SWIFT_API/scripts/tests/custom
chmod +x sequence_test.sh
./sequence_test.sh
```

## Docker Configuration

The project uses a multi-stage build process for the API service:

1. Development mode (default in Dockerfile):
   - Uses a single-stage build
   - Mounts the source code directory for live code reloading
   - Recommended for development work

2. Production mode (commented in Dockerfile):
   - Uses a multi-stage build for a smaller final image
   - Builds a static binary
   - More efficient and secure for production deployment

To switch between modes, modify the Dockerfile by commenting/uncommenting the relevant sections.

## Code Structure

The project follows a clean architecture pattern:

- `cmd/api`: Application entry point
- `internal/api`: API layer with handlers and middleware
- `internal/db`: Data access layer
- `internal/service`: Business logic
- `pkg/validators`: Reusable validation components
- `configs`: Configuration files including default data

## Volumes

The docker-compose setup creates and uses the following volumes:

- `mongodb_data`: Persists MongoDB data between container restarts
- `go-mod-cache`: Caches Go modules to speed up builds

## Troubleshooting

If you encounter issues with the application:

1. Check container status:
   ```bash
   docker-compose ps
   ```

2. View container logs:
   ```bash
   docker-compose logs api
   docker-compose logs db
   ```

3. Verify MongoDB connection:
   ```bash
   docker exec -it mongo mongosh
   ```

4. Ensure the API service has access to the required environment variables and volumes.

5. For debugging, set `LOGGER_DEBUG=true` and `API_DEBUG=true` in the docker-compose file.

## Code details

### Code Structure

Architecture pattern of the project prestents itself as follows:

`cmd/api`: The main application
`internal/api`: API layer with request handler and loggers in middleware
`internal/db`: Data access layer, main database connection
`internal/service`: Business logic, as described below
`pkg/validators`: Validators and sanitizers for data handling 
`configs`: Configuration files, in particular inital data, there I could also create a 'settings' file

### Particular parts

#### Main 

Ah, the main. Here, every programmer's journey begins with a simple 'hello world' inside. I implemented the bare minimum logic to set up the server. There, I resided objects used for logging, database connection, data loading, and finally, request handling. 

#### Router

What would we do without a humble router? Its main role is to call the functions defined elsewhere if a certain endpoint is called. It uses a logger to write out requests and responses (if configured). Endpoints are all the same as above, except for one - health check. It was a debug endpoint to verify if the connection to the server is available. I used Postman to ping this port and was relieved to see the "OK" in the response body.

#### Request handler 

If the router is the waiter, the request handler is the kitchen, where the waiter brings the order, and here the magic happens. The handler takes the request's body and passes it further to the service, then replies with whatever the service says. It doesn't do much "thinking" except for the simple error check on the service part.

#### Service

To drag the cliche analogy - the service here is the cook. It verifies the data for database insertion or provides logic for data handling like trimming the excess field values from the database to match the required endpoint response.

#### Parser

Extracts the data from the file base file and uses the service to post it into the database. Remember that the delimiter is set to ';' (semicolon) since there were commas present in the data that could mess up the structure.

#### Mongo repository 

Connects to MongoDB and can be used to Create, Read, Update, and Delete documents (the fields in the database) into and from the collection (like a table in a relational database). Mongo was my primary choice since it's considered to be very fast and in this project, there wasn't much need for relational databases, since no user is required with additional connections that the user has to the data. The API has to store the given data and provide it in a readable format. There is a part of "relations" though, the country code is tied to the country-specific data since I wasn't sure what to do with the TimeZone field for example, which didn't need to be returned in the endpoint but was provided nonetheless. Therefore, I stored country data (country ISO2 code, country name, and time zone) in another collection.    

#### Validators and sanitizers

The code in each one should be self-explanatory, they check if the data complies with set requirements, mostly using Regex and expected lengths checks. There is also a sanitizer in place just in case of a NoSQL injection attack similar to the classic SQL injection.

## Usage of AI

Whoever says that they wouldn't use AI these days is either a liar or above the level of a mid-developer (Because of course Seniors don't write code). 

--- 

I used some help from popular AI-s that being ChatGPT and Claude for comments, documentation, and test generation. As for the tests, I verified the code for the correctness of logic, I made simple adjustments where it was needed and that's how these tools should be used. Also, I used it to essentially translate my bash script, since that's what I'm most proficient with and I needed to to have other scripts to set the application up on the Windows environment. They are there just for the sake of completeness of the project, not for it to be run using actions and they are not recommended, as a dev use Linux as much as you can.