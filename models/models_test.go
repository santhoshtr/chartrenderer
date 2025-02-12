package models_test

import (
	"encoding/json"
	"testing"

	"github.com/santhoshtr/chartadapter/models"
	"github.com/stretchr/testify/assert"
)

func TestChartDefinitionParsing(t *testing.T) {
	jsonData := `{
		"version": 1,
		"license": "CC0-1.0",
		"source": "Chart Example Data.tab",
		"type": "line",
		"title": {
			"en": "Sample Line Chart"
		},
		"xAxis": {
			"title": {
				"en": "Date"
			}
		},
		"yAxis": {
			"title": {
				"en": "support%"
			}
		}
	}`

	var chartDef models.ChartDefinition

	err := json.Unmarshal([]byte(jsonData), &chartDef)
	assert.NoError(t, err)

	assert.Equal(t, 1, chartDef.Version)
	assert.Equal(t, "CC0-1.0", chartDef.License)
	assert.Equal(t, "Chart Example Data.tab", chartDef.Source)
	assert.Equal(t, "line", chartDef.Type.String())
	assert.Equal(t, "Sample Line Chart", chartDef.Title["en"])
	assert.Equal(t, "Date", chartDef.XAxis.Title["en"])
	assert.Equal(t, "support%", chartDef.YAxis.Title["en"])
}

func TestChartDataParsing(t *testing.T) {
	jsonData := `{
		"license": "CC0-1.0",
		"description": {
			"en": "Some meaningless example data about Middle-Earth"
		},
		"schema": {
			"fields": [
				{
					"name": "a1",
					"type": "string",
					"title": {
						"en": "Date"
					}
				},
				{
					"name": "a2",
					"type": "number",
					"title": {
						"en": "Elves"
					}
				},
				{
					"name": "a3",
					"type": "number",
					"title": {
						"en": "Ents"
					}
				},
				{
					"name": "a4",
					"type": "number",
					"title": {
						"en": "Orcs"
					}
				},
				{
					"name": "a5",
					"type": "number",
					"title": {
						"en": "Hobbits"
					}
				},
				{
					"name": "a6",
					"type": "number",
					"title": {
						"en": "Trolls"
					}
				}
			]
		},
		"data": [
			["1993/09/09", 35, 37, 8, 8, 10],
			["1993/09/14", 36, 33, 8, 10, 11],
			["1993/09/20", 35, 35, 6, 11, 11],
			["1993/09/25", 30, 37, 8, 10, 13],
			["1993/09/26", 31, 36, 7, 11, 13],
			["1993/09/26", 28, 34, 7, 12, 15],
			["1993/09/30", 25, 39, 6, 12, 17],
			["1993/10/02", 26, 38, 8, 12, 14],
			["1993/10/08", 22, 37, 8, 12, 18],
			["1993/10/16", 22, 40, 7, 13, 16],
			["1993/10/19", 21, 39, 6, 14, 17],
			["1993/10/22", 18, 43, 7, 14, 18],
			["1993/10/22", 16, 44, 7, 12, 19],
			["1993/10/25", 16, 41, 1, 13, 18]
		]
	}`

	var chartData models.ChartData

	err := json.Unmarshal([]byte(jsonData), &chartData)
	assert.NoError(t, err)
	assert.Equal(t, "CC0-1.0", chartData.License)
	assert.Equal(t, "Some meaningless example data about Middle-Earth", chartData.Description["en"])
	assert.Equal(t, "a1", chartData.Schema.Fields[0].Name)
	assert.Equal(t, "string", chartData.Schema.Fields[0].Type.String())
	assert.Equal(t, "Date", chartData.Schema.Fields[0].Title["en"])
	assert.Equal(t, "a2", chartData.Schema.Fields[1].Name)
	assert.Equal(t, "number", chartData.Schema.Fields[1].Type.String())
	assert.Equal(t, "Elves", chartData.Schema.Fields[1].Title["en"])
	assert.Equal(t, "1993/09/09", chartData.Data[0][0])
	assert.Equal(t, 35.0, chartData.Data[0][1])
	assert.Equal(t, 37.0, chartData.Data[0][2])
	assert.Equal(t, 8.0, chartData.Data[0][3])
	assert.Equal(t, 8.0, chartData.Data[0][4])
	assert.Equal(t, 10.0, chartData.Data[0][5])
}
