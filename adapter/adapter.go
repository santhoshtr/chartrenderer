package adapter

import (
	"fmt"

	"github.com/santhoshtr/chartadapter/models"
)

type EChartsAdapter struct {
	Definition *models.ChartDefinition
	Data       *models.ChartData
}

func NewEChartsAdapter(def *models.ChartDefinition, data *models.ChartData) *EChartsAdapter {
	return &EChartsAdapter{
		Definition: def,
		Data:       data,
	}
}

func (a *EChartsAdapter) Convert(language string) (map[string]interface{}, error) {
	switch a.Definition.Type {
	case models.Line:
		chart, err := a.convertLineChart(language)
		return chart.JSON(), err
	case models.Bar:
		chart, err := a.convertBarChart(language)
		return chart.JSON(), err
	case models.Area:
		chart, err := a.convertAreaChart(language)
		return chart.JSON(), err
	case models.Scatter:
		chart, err := a.convertScatterChart(language)
		return chart.JSON(), err
	case models.Pie:
		chart, err := a.convertPieChart(language)
		return chart.JSON(), err
	default:
		return nil, fmt.Errorf("unsupported chart type: %s", a.Definition.Type)
	}
}

func (a *EChartsAdapter) LocalizedTitle(title map[string]string, language string) string {
	localizedTitle, ok := title[language]
	if !ok {
		localizedTitle = title["en"]
	}
	return localizedTitle
}
