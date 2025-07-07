# Go Datadog Logs API Example

This project provides a simple Go application that demonstrates how to send structured logs to the Datadog Logs API using the official `datadog-api-client-go` library.

The example defines a custom log structure (`appLogBody`), populates it with data, and then maps it to Datadog's `HTTPLogItem` format. It shows how to set standard log attributes like `service` and `ddtags`, and how to include the original, richer log object as custom attributes for full visibility in Datadog.

## Prerequisites

- Go (version 1.18 or later)
- A Datadog account

## Setup and Installation

1.  **Initialize Go Module:**
    If you haven't already, initialize a Go module in the project directory:
    ```sh
    go mod init datadog-logs-api-go
    ```

2.  **Install Dependencies:**
    Tidy will download the necessary Datadog client library.
    ```sh
    go mod tidy
    ```

## Configuration

Before running the application, you must configure your Datadog credentials. The client library automatically reads the API key and site from environment variables.

1.  **Export your Datadog API Key:**
    ```sh
    export DD_API_KEY="<YOUR_DATADOG_API_KEY>"
    ```

2.  **Export your Datadog Site:**
    Your site is the domain you use to access Datadog (e.g., for US1 - `datadoghq.com`, US5 - `us5.datadoghq.com`).
    ```sh
    export DD_SITE="<YOUR_DATADOG_SITE>"
    ```

## Usage

Run the `main.go` application from your terminal:

```sh
go run main.go
```

### Expected Output

If the request is successful, Datadog will return a `202 Accepted` status, and the application will print the following output:

```
Log submission accepted (Status: 202 Accepted)
Response body (should be empty): 
```

Your log will now be available in the Datadog Log Explorer. You can search for it by the service (`payment-gateway`), the tags (`env:staging`), or any of the custom attributes from the `appLogBody` struct.
