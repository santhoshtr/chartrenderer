package adapter

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (a *EChartsAdapter) convertPieChart(language string) (*charts.Pie, error) {
	// Create a new pie chart
	pie := charts.NewPie()

	// Set global options
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: a.LocalizedTitle(a.Definition.Title, language),
			Subtitle: a.LocalizedTitle(a.Data.Description, language),
		}),
	)

	// Process data series
	items := make([]opts.PieData, 0)
	for _, row := range a.Data.Data {
		if len(row) >= 2 {
			// Convert x-value to string
			name := fmt.Sprintf("%v", row[0])

			// Convert y-value to float64
			if value, ok := row[1].(float64); ok {
				items = append(items, opts.PieData{Name: name, Value: value})
			}
		}
	}

	// Add data to the chart
	pie.SetSeriesOptions(
		charts.WithPieChartOpts(opts.PieChart{
			RoseType: "radius",
		}),
	)

	pie.AddSeries("Data", items)

	return pie, nil
}
