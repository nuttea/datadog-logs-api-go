// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

// Your existing log data schema
type appLogBody struct {
	StartTime      int64   `json:"start_time"`
	EventTime      int64   `json:"event_time"`
	LogEnv         string  `json:"log_env"`
	ApplogVersion  string  `json:"applog_version"`
	SelfType       string  `json:"self_type"`
	SelfGroup      string  `json:"self_group"`
	SelfSystem     string  `json:"self_system"`
	SelfFunction   string  `json:"self_function"`
	CallDirection  int     `json:"call_direction"`
	CallReqParams  string  `json:"call_req_params"`
	CallReqHeaders string  `json:"call_req_headers"`
	CallReqMethods string  `json:"call_req_methods"`
	CallResBody    string  `json:"call_res_body"`
	CallResStatus  int     `json:"call_res_status"`
	CallResTime    int     `json:"call_res_time"`
	CallSeverity   int     `json:"call_severity"`
	ApplogURL      string  `json:"applog_url"`
	Message        string  `json:"message"`
	CartID         string  `json:"cart_id,omitempty"`
	UserID         string  `json:"user_id,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
}

func main() {
	// Generate 25 unique user IDs
	userIDs := make([]string, 25)
	for i := 0; i < 25; i++ {
		userIDs[i] = fmt.Sprintf("user-%d", i+1)
	}

	// Generate 25 unique cart IDs
	cartIDs := make([]string, 25)
	for i := 0; i < 25; i++ {
		cartIDs[i] = fmt.Sprintf("cart-%d", i+1)
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Configure the Datadog client
	// IMPORTANT: Set DD_API_KEY and DD_SITE environment variables.
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewLogsApi(apiClient)

	// Loop to send logs every 1 second
	for {
		// Create a batch of 50 log entries
		var logBatch []datadogV2.HTTPLogItem
		for i := 0; i < 50; i++ {
			// Randomly select a user ID and cart ID
			userID := userIDs[rand.Intn(len(userIDs))]
			cartID := cartIDs[rand.Intn(len(cartIDs))]

			// Generate a random amount
			amount := rand.Float64() * 1000 // Random amount between 0 and 1000

			// Create the log entry
			myLog := appLogBody{
				StartTime:     time.Now().Unix() - int64(rand.Intn(60)), // Random start time in the last minute
				EventTime:     time.Now().Unix(),
				LogEnv:        "staging",
				ApplogVersion: "v1.2.3",
				SelfSystem:    "payment-gateway",
				SelfFunction:  "processTransaction",
				Message:       fmt.Sprintf("Transaction for user %s, cart %s", userID, cartID),
				CallResStatus: 200,
				CallSeverity:  1,
				CartID:        cartID,
				UserID:        userID,
				Amount:        amount,
			}

			// Convert your struct to a map to use as custom attributes.
			var attributes map[string]interface{}
			data, err := json.Marshal(myLog)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling log data: %v\n", err)
				continue
			}
			json.Unmarshal(data, &attributes)

			// Create the Datadog HTTPLogItem.
			logItem := datadogV2.HTTPLogItem{
				Message:              myLog.Message,
				Service:              datadog.PtrString(myLog.SelfSystem),
				Hostname:             datadog.PtrString("i-012345678"),
				Ddtags:               datadog.PtrString(fmt.Sprintf("env:%s,version:%s", myLog.LogEnv, myLog.ApplogVersion)),
				Ddsource:             datadog.PtrString("go"),
				AdditionalProperties: attributes,
			}
			logBatch = append(logBatch, logItem)
		}

		// Prepare the request body.
		body := logBatch

		// Send the log batch.
		resp, r, err := api.SubmitLog(ctx, body, *datadogV2.NewSubmitLogOptionalParameters())

		// Handle the response.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `LogsApi.SubmitLog`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		} else {
			fmt.Printf("Log batch submission accepted (Status: %s)\n", r.Status)
			fmt.Printf("Response body (should be empty): %s\n", resp)
		}

		// Wait for 1 second before sending the next batch
		time.Sleep(1 * time.Second)
	}
}
