package main

import (
	"encoding/json"
	"os"

	"github.com/santhoshtr/chartadapter/adapter"
	"github.com/santhoshtr/chartadapter/models"
)

type MWChartOptions struct {
	Locale     string                 `json:"locale"`
	Definition models.ChartDefinition `json:"definition"`
	Data       models.ChartData       `json:"data"`
}

func main() {
	if len(os.Args) < 2 {
		panic("Please provide the JSON file name as a command line argument")
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var MWChartOptions MWChartOptions
	if err := json.NewDecoder(file).Decode(&MWChartOptions); err != nil {
		panic(err)
	}

	def := MWChartOptions.Definition
	data := MWChartOptions.Data
	// Create adapter
	adapter := adapter.NewEChartsAdapter(&def, &data)
	if MWChartOptions.Locale == "" {
		MWChartOptions.Locale = "en"
	}
	// Convert to ECharts
	echartOptions, err := adapter.Convert(MWChartOptions.Locale)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(echartOptions)
}
