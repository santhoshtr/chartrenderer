//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/santhoshtr/chartadapter/adapter"
	"github.com/santhoshtr/chartadapter/models"
)

// MWChartOptions combines chart definition and data
type MWChartOptions struct {
	Locale string `json:"locale"`
    Definition models.ChartDefinition `json:"definition"`
    Data       models.ChartData      `json:"data"`
}

// convertToEChartsOptions converts our data format to ECharts options
func convertToEChartsOptions(data MWChartOptions) map[string]interface{} {
    def := data.Definition
    chartData := data.Data
    // Create adapter
    adapter := adapter.NewEChartsAdapter(&def, &chartData)
    if  data.Locale == "" {
        data.Locale = "en"
    }
    // Convert to ECharts
    echartOptions, err := adapter.Convert(data.Locale)
    if err != nil {
        panic(err)
    }


    return echartOptions
}

// getEChartOptions is the exported JavaScript function
func getEChartOptions(this js.Value, args []js.Value) interface{} {
    if len(args) < 1 {
        return map[string]interface{}{
            "error": "No input data provided",
        }
    }
	jsJSON := js.Global().Get("JSON")
	jsonStr := jsJSON.Call("stringify", args[0]).String()

    // Parse input JSON
    var mwChartOptions MWChartOptions
    if err := json.Unmarshal([]byte(jsonStr), &mwChartOptions); err != nil {
        return map[string]interface{}{
            "error": err.Error(),
        }
    }
    // Convert to ECharts options
	result := convertToEChartsOptions(mwChartOptions)
	jsResult, err := json.Marshal(result)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	return js.Global().Get("JSON").Call("parse", string(jsResult))

}

func main() {
    // Create channel to keep main running
    c := make(chan struct{}, 0)

    // Register function
    js.Global().Set("GetEchartOptions", js.FuncOf(getEChartOptions))

    <-c
}