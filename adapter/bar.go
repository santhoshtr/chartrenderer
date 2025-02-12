package adapter

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (a *EChartsAdapter) convertBarChart(language string) (*charts.Bar, error) {
	// Create a new bar chart
	bar := charts.NewBar()

	xValues := make([]string, 0)
	for _, row := range a.Data.Data {
		// Convert x-value to string
		xValues = append(xValues, fmt.Sprintf("%s", row[0]))
	}

	// Set global options
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    a.LocalizedTitle(a.Definition.Title, language),
			Subtitle: a.LocalizedTitle(a.Data.Description, language),
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
			bar.AddSeries(a.LocalizedTitle(field.Title, language), generateBarItems(yValues))
		}
	}

	return bar, nil
}

func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, 0)
	for _, v := range values {
		items = append(items, opts.BarData{Value: v})
	}
	return items
}
