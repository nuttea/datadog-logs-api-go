// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

// Your existing log data schema
type appLogBody struct {
	StartTime      int64  `json:"start_time"`
	EventTime      int64  `json:"event_time"`
	LogEnv         string `json:"log_env"`
	ApplogVersion  string `json:"applog_version"`
	SelfType       string `json:"self_type"`
	SelfGroup      string `json:"self_group"`
	SelfSystem     string `json:"self_system"`
	SelfFunction   string `json:"self_function"`
	CallDirection  int    `json:"call_direction"`
	CallReqParams  string `json:"call_req_params"`
	CallReqHeaders string `json:"call_req_headers"`
	CallReqMethods string `json:"call_req_methods"`
	CallResBody    string `json:"call_res_body"`
	CallResStatus  int    `json:"call_res_status"`
	CallResTime    int    `json:"call_res_time"`
	CallSeverity   int    `json:"call_severity"`
	ApplogURL      string `json:"applog_url"`
	Message        string `json:"message"`
}

func main() {
	// 1. Create an instance of your log data structure.
	myLog := appLogBody{
		StartTime:      1678886400,
		EventTime:      1678886401,
		LogEnv:         "staging",
		ApplogVersion:  "v1.2.3",
		SelfSystem:     "payment-gateway",
		SelfFunction:   "processTransaction",
		Message:        "Transaction processed successfully",
		CallResStatus:  200,
		CallSeverity:   1,
	}

	// 2. Convert your struct to a map to use as custom attributes.
	var attributes map[string]interface{}
	data, err := json.Marshal(myLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log data: %v\n", err)
		return
	}
	json.Unmarshal(data, &attributes)

	// 3. Create the Datadog HTTPLogItem.
	logItem := datadogV2.HTTPLogItem{
		Message:  myLog.Message,
		Service:  datadog.PtrString(myLog.SelfSystem),
		Hostname: datadog.PtrString("i-012345678"),
		Ddtags:   datadog.PtrString(fmt.Sprintf("env:%s,version:%s", myLog.LogEnv, myLog.ApplogVersion)),
		AdditionalProperties: attributes,
	}

	// 4. Prepare the request body.
	body := []datadogV2.HTTPLogItem{logItem}

	// 5. Configure the Datadog client and send the log.
	// IMPORTANT: Set DD_API_KEY and DD_SITE environment variables.
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewLogsApi(apiClient)
	resp, r, err := api.SubmitLog(ctx, body, *datadogV2.NewSubmitLogOptionalParameters())

	// 6. Handle the response.
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LogsApi.SubmitLog`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return
	}

	fmt.Printf("Log submission accepted (Status: %s)\n", r.Status)
	fmt.Printf("Response body (should be empty): %s\n", resp)
}
