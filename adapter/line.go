package adapter

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (a *EChartsAdapter) convertLineChart(language string) (*charts.Line, error) {
	// Create a new line chart
	line := charts.NewLine()

	xValues := make([]string, 0)
	for _, row := range a.Data.Data {
		// Convert x-value to string
		xValues = append(xValues, fmt.Sprintf("%s", row[0]))
	}

	// Set global options
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: a.LocalizedTitle(a.Definition.Title, language),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: a.LocalizedTitle(a.Definition.XAxis.Title, language),
			Type: "category",
			Data: xValues,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: a.LocalizedTitle(a.Definition.YAxis.Title, language),
			Type: "value",
		}),
	)

	for fieldIndex, field := range a.Data.Schema.Fields {
		if fieldIndex == 0 {
			continue
		}
		if field.Type == "number" {
			yValues := make([]float64, 0)
			for _, row := range a.Data.Data {
				if len(row) >= 2 {
					// Convert y-value to float64
					if yVal, ok := row[fieldIndex].(float64); ok {
						yValues = append(yValues, yVal)
					}
				}
			}
			line.AddSeries(a.LocalizedTitle(field.Title, language), generateLineItems(yValues))
		}
	}

	return line, nil
}

func generateLineItems(values []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, v := range values {
		items = append(items, opts.LineData{Value: v})
	}
	return items
}
