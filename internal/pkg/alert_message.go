package pkg

import (
	"crypto-alerts/internal/entity"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

type AlertMessage struct {
	Symbol         string
	Name           string
	Price          float64
	Volume         float64
	Period         string
	Variation      float64
	Threshold      float64
	Direction      string
	IsTargetPrice  bool
	TargetPrice    float64
	FearGreedValue int
	FearGreedClass string
	HistoricalData *entity.HistoricalPriceData
}

func formatLargeNumber(value float64) string {
	abs := math.Abs(value)
	sign := ""
	if value < 0 {
		sign = "-"
	}

	switch {
	case abs >= 1e12:
		return fmt.Sprintf("%s%.1fT", sign, abs/1e12)
	case abs >= 1e9:
		return fmt.Sprintf("%s%.1fB", sign, abs/1e9)
	case abs >= 1e6:
		return fmt.Sprintf("%s%.1fM", sign, abs/1e6)
	case abs >= 1e3:
		return fmt.Sprintf("%s%.1fK", sign, abs/1e3)
	default:
		return fmt.Sprintf("%s%.2f", sign, abs)
	}
}

func FormatEmailSubject(message AlertMessage) string {
	if message.IsTargetPrice {
		return FormatTargetPriceEmailSubject(message)
	}

	direction := "subiu"
	emoji := "🟢"
	if message.Direction == "down" {
		direction = "caiu"
		message.Variation = -message.Variation
		emoji = "🔴"
	}

	return fmt.Sprintf("%s %s %s %.2f%% em %s: Preço atual $%.2f", emoji, message.Symbol, direction, message.Variation, message.Period, message.Price)
}

func FormatEmailBody(message AlertMessage) string {
	if message.IsTargetPrice {
		return FormatTargetPriceEmailBody(message)
	}

	var directionText string
	if message.Direction == "up" {
		directionText = "subiu acima"
	} else {
		directionText = "caiu abaixo"
	}

	content := strings.Builder{}
	content.WriteString("<html><body style='font-family: Arial, sans-serif; line-height: 1.6; color: #333;'>")
	content.WriteString(fmt.Sprintf("<p>Olá,</p><p>Temos um alerta de preço para a criptomoeda <strong>%s (%s)</strong>!</p>", message.Name, message.Symbol))
	content.WriteString(fmt.Sprintf("<p>A variação no período de %s %s do seu alerta configurado de %.2f%%, atingindo <strong>%.2f%%</strong>.</p>",
		message.Period, directionText, message.Threshold, message.Variation))
	content.WriteString("<h3>Detalhes atuais:</h3><ul>")
	content.WriteString(fmt.Sprintf("<li>Preço Atual: <strong>$%.2f USD</strong></li>", message.Price))
	content.WriteString(fmt.Sprintf("<li>Volume negociado nas últimas 24h: <strong>$%s USD</strong></li>", formatLargeNumber(message.Volume)))
	content.WriteString(fmt.Sprintf("<li>Variação no período (%s): <strong>%.2f%%</strong></li></ul>", message.Period, message.Variation))

	if message.HistoricalData != nil {
		historicalChartURL := GenerateHistoricalPriceChartURL(message.HistoricalData)
		if historicalChartURL != "" {
			content.WriteString("<div style='margin: 20px 0; padding: 0; text-align: center;'>")
			content.WriteString("<h3 style='margin-bottom: 5px; color: #333;'>Histórico de Preço (90 dias)</h3>")
			content.WriteString(fmt.Sprintf("<p style='margin-top: 0; margin-bottom: 15px; color: #666; font-size: 14px;'>Mín: $%.2f | Máx: $%.2f | Média: $%.2f</p>",
				message.HistoricalData.MinPrice, message.HistoricalData.MaxPrice, message.HistoricalData.AvgPrice))
			content.WriteString(fmt.Sprintf("<img src='%s' alt='Historical Price Chart' style='max-width: 800px; width: 100%%; height: auto; border-radius: 8px;'/>", historicalChartURL))
			content.WriteString("</div>")
		}

		volumeChartURL := GenerateHistoricalVolumeChartURL(message.HistoricalData)
		if volumeChartURL != "" {
			content.WriteString("<div style='margin: 20px 0; padding: 0; text-align: center;'>")
			content.WriteString("<h3 style='margin-bottom: 5px; color: #333;'>Volume Negociado (90 dias)</h3>")
			content.WriteString(fmt.Sprintf("<p style='margin-top: 0; margin-bottom: 15px; color: #666; font-size: 14px;'>Mín: $%.2fB | Máx: $%.2fB | Média: $%.2fB</p>",
				message.HistoricalData.MinVolume/1e9, message.HistoricalData.MaxVolume/1e9, message.HistoricalData.AvgVolume/1e9))
			content.WriteString(fmt.Sprintf("<img src='%s' alt='Historical Volume Chart' style='max-width: 800px; width: 100%%; height: auto; border-radius: 8px;'/>", volumeChartURL))
			content.WriteString("</div>")
		}
	}

	if message.FearGreedClass != "" {
		fearGreedChartURL := generateFearGreedChartURL(message.FearGreedValue)
		content.WriteString("<div style='margin: 20px 0; padding: 0; text-align: center;'>")
		content.WriteString(fmt.Sprintf("<h3 style='margin-bottom: 5px; color: #333;'>Índice Fear & Greed do Mercado</h3>"))
		content.WriteString(fmt.Sprintf("<p style='margin-top: 0; margin-bottom: 15px; color: #666; font-size: 16px;'>%s</p>", message.FearGreedClass))
		content.WriteString(fmt.Sprintf("<img src='%s' alt='Fear & Greed Index' style='max-width: 450px; width: 100%%; height: auto; border-radius: 8px;'/>", fearGreedChartURL))
		content.WriteString("</div>")
	}

	content.WriteString("<p>Este é um bom momento para verificar seus investimentos e decidir os próximos passos.</p>")
	content.WriteString("<p>Atenciosamente,<br/>Equipe Crypto Alerts</p>")
	content.WriteString("<hr/><p style='font-size: 0.9em; color: #666;'>Este é um e-mail automático. Por favor, não responda.</p>")
	content.WriteString("</body></html>")

	return content.String()
}

func FormatTargetPriceEmailSubject(message AlertMessage) string {
	emoji := "🎯"
	direction := ""

	if message.Direction == "up" {
		direction = "ultrapassou"
	} else {
		direction = "caiu abaixo de"
	}

	return fmt.Sprintf("%s Preço Alvo: %s %s $%.2f (atual: $%.2f)", emoji, message.Symbol, direction, message.TargetPrice, message.Price)
}

func FormatTargetPriceEmailBody(message AlertMessage) string {
	var actionSuggestion, directionText string

	if message.Direction == "up" {
		directionText = "atingiu ou ultrapassou"
		actionSuggestion = "Este pode ser um bom momento para considerar vender, dependendo da sua estrategia."
	} else {
		directionText = "atingiu ou ficou abaixo"
		actionSuggestion = "Este pode ser um bom momento para considerar comprar, dependendo da sua estrategia."
	}

	content := strings.Builder{}
	content.WriteString("<html><body style='font-family: Arial, sans-serif; line-height: 1.6; color: #333;'>")
	content.WriteString("<p>Olá,</p>")
	content.WriteString("<p><strong>Seu alerta de preço alvo foi acionado!</strong></p>")
	content.WriteString(fmt.Sprintf("<p>A criptomoeda <strong>%s (%s)</strong> %s seu preço alvo configurado de <strong>$%.2f USD</strong>.</p>",
		message.Name, message.Symbol, directionText, message.TargetPrice))

	content.WriteString("<h3>Detalhes atuais do mercado:</h3><ul>")
	content.WriteString(fmt.Sprintf("<li>Preço Atual: <strong>$%.2f USD</strong></li>", message.Price))
	content.WriteString(fmt.Sprintf("<li>Volume negociado nas últimas 24h: <strong>$%s USD</strong></li></ul>", formatLargeNumber(message.Volume)))

	if message.HistoricalData != nil {
		historicalChartURL := GenerateHistoricalPriceChartURL(message.HistoricalData)
		if historicalChartURL != "" {
			content.WriteString("<div style='margin: 20px 0; padding: 0; text-align: center;'>")
			content.WriteString("<h3 style='margin-bottom: 5px; color: #333;'>Histórico de Preço (90 dias)</h3>")
			content.WriteString(fmt.Sprintf("<p style='margin-top: 0; margin-bottom: 15px; color: #666; font-size: 14px;'>Mín: $%.2f | Máx: $%.2f | Média: $%.2f</p>",
				message.HistoricalData.MinPrice, message.HistoricalData.MaxPrice, message.HistoricalData.AvgPrice))
			content.WriteString(fmt.Sprintf("<img src='%s' alt='Historical Price Chart' style='max-width: 800px; width: 100%%; height: auto; border-radius: 8px;'/>", historicalChartURL))
			content.WriteString("</div>")
		}

		volumeChartURL := GenerateHistoricalVolumeChartURL(message.HistoricalData)
		if volumeChartURL != "" {
			content.WriteString("<div style='margin: 20px 0; padding: 0; text-align: center;'>")
			content.WriteString("<h3 style='margin-bottom: 5px; color: #333;'>Volume Negociado (90 dias)</h3>")
			content.WriteString(fmt.Sprintf("<p style='margin-top: 0; margin-bottom: 15px; color: #666; font-size: 14px;'>Mín: $%.2fB | Máx: $%.2fB | Média: $%.2fB</p>",
				message.HistoricalData.MinVolume/1e9, message.HistoricalData.MaxVolume/1e9, message.HistoricalData.AvgVolume/1e9))
			content.WriteString(fmt.Sprintf("<img src='%s' alt='Historical Volume Chart' style='max-width: 800px; width: 100%%; height: auto; border-radius: 8px;'/>", volumeChartURL))
			content.WriteString("</div>")
		}
	}

	if message.FearGreedClass != "" {
		fearGreedChartURL := generateFearGreedChartURL(message.FearGreedValue)
		content.WriteString("<div style='margin: 20px 0; padding: 0; text-align: center;'>")
		content.WriteString(fmt.Sprintf("<h3 style='margin-bottom: 5px; color: #333;'>Índice Fear & Greed do Mercado</h3>"))
		content.WriteString(fmt.Sprintf("<p style='margin-top: 0; margin-bottom: 15px; color: #666; font-size: 16px;'>%s</p>", message.FearGreedClass))
		content.WriteString(fmt.Sprintf("<img src='%s' alt='Fear & Greed Index' style='max-width: 450px; width: 100%%; height: auto; border-radius: 8px;'/>", fearGreedChartURL))
		content.WriteString("</div>")
	}

	content.WriteString(fmt.Sprintf("<p>%s</p>", actionSuggestion))
	content.WriteString("<p>Atenciosamente,<br/>Equipe Crypto Alerts</p>")
	content.WriteString("<hr/><p style='font-size: 0.9em; color: #666;'>Este e um e-mail automatico. Por favor, não responda.</p>")
	content.WriteString("</body></html>")

	return content.String()
}

func generateFearGreedChartURL(value int) string {
	chartConfig := fmt.Sprintf(`{
		"type": "gauge",
		"data": {
			"datasets": [{
				"value": %d,
				"data": [0, 25, 50, 75, 100],
				"backgroundColor": ["#EA3943", "#F59E0B", "#EAB308", "#93D900", "#16C784"],
				"borderWidth": 8,
				"borderColor": "#2D3748",
				"spacing": 3,
				"cutout": "75%%",
				"circumference": 180,
				"rotation": 270
			}]
		},
		"options": {
			"responsive": true,
			"maintainAspectRatio": false,
			"needle": {
				"radiusPercentage": 3,
				"widthPercentage": 2,
				"lengthPercentage": 85,
				"color": "#FFFFFF"
			},
			"valueLabel": {
				"display": true,
				"fontSize": 48,
				"fontStyle": "bold",
				"fontFamily": "Arial",
				"color": "#FFFFFF",
				"backgroundColor": "transparent",
				"bottomMarginPercentage": 15
			},
			"plugins": {
				"datalabels": {
					"display": false
				}
			},
			"layout": {
				"padding": {
					"top": 10,
					"bottom": 10
				}
			}
		}
	}`, value)

	encodedChart := url.QueryEscape(chartConfig)
	return fmt.Sprintf("https://quickchart.io/chart?c=%s&width=450&height=280&backgroundColor=%%232D3748", encodedChart)
}

func GenerateHistoricalPriceChartURL(historicalData *entity.HistoricalPriceData) string {
	if historicalData == nil || len(historicalData.Prices) == 0 {
		return ""
	}

	var priceValues []string
	var dateLabels []string

	for i, point := range historicalData.Prices {
		priceValues = append(priceValues, fmt.Sprintf("%.2f", point.Price))

		t := time.Unix(point.Timestamp/1000, 0)
		day := t.Day()

		if i%10 == 0 {
			dateLabels = append(dateLabels, fmt.Sprintf("'%d'", day))
		} else {
			dateLabels = append(dateLabels, "''")
		}
	}

	chartConfig := fmt.Sprintf(`{
		"type": "line",
		"data": {
			"labels": [%s],
			"datasets": [{
				"label": "%s Price (USD)",
				"data": [%s],
				"borderColor": "rgb(99, 102, 241)",
				"backgroundColor": "rgba(99, 102, 241, 0.1)",
				"borderWidth": 2,
				"pointRadius": 0,
				"pointHoverRadius": 4,
				"fill": true,
				"tension": 0.4
			}]
		},
		"options": {
			"responsive": true,
			"plugins": {
				"legend": {
					"display": true,
					"position": "top",
					"labels": {
						"color": "rgb(255, 255, 255)",
						"font": {
							"size": 12,
							"weight": "bold"
						}
					}
				},
				"title": {
					"display": true,
					"text": "90-Day Price History",
					"color": "rgb(255, 255, 255)",
					"font": {
						"size": 14,
						"weight": "bold"
					}
				}
			},
			"scales": {
				"x": {
					"display": true,
					"title": {
						"display": true,
						"text": "Day of Month",
						"color": "rgb(255, 255, 255)",
						"font": {
							"size": 11,
							"weight": "bold"
						}
					},
					"ticks": {
						"color": "rgb(200, 200, 200)",
						"font": {
							"size": 9
						},
						"maxRotation": 0,
						"autoSkip": false
					},
					"grid": {
						"color": "rgba(255, 255, 255, 0.1)"
					}
				},
				"y": {
					"display": true,
					"title": {
						"display": true,
						"text": "Price (USD)",
						"color": "rgb(255, 255, 255)",
						"font": {
							"size": 11,
							"weight": "bold"
						}
					},
					"ticks": {
						"color": "rgb(200, 200, 200)",
						"font": {
							"size": 10
						},
						"callback": "function(value) { return '$' + value.toFixed(2); }"
					},
					"grid": {
						"color": "rgba(255, 255, 255, 0.1)"
					}
				}
			}
		}
	}`,
		strings.Join(dateLabels, ","),
		historicalData.Symbol,
		strings.Join(priceValues, ","))

	encodedChart := url.QueryEscape(chartConfig)
	return fmt.Sprintf("https://quickchart.io/chart?c=%s&width=800&height=400&backgroundColor=%%232D3748", encodedChart)
}

func GenerateHistoricalVolumeChartURL(historicalData *entity.HistoricalPriceData) string {
	if historicalData == nil || len(historicalData.Volumes) == 0 {
		return ""
	}

	var volumeValues []string
	var dateLabels []string

	for i, point := range historicalData.Volumes {
		volumeValueInBillions := point.Volume / 1e9
		volumeValues = append(volumeValues, fmt.Sprintf("%.2f", volumeValueInBillions))

		t := time.Unix(point.Timestamp/1000, 0)
		day := t.Day()

		if i%10 == 0 {
			dateLabels = append(dateLabels, fmt.Sprintf("'%d'", day))
		} else {
			dateLabels = append(dateLabels, "''")
		}
	}

	chartConfig := fmt.Sprintf(`{
		"type": "bar",
		"data": {
			"labels": [%s],
			"datasets": [{
				"label": "%s Volume (USD)",
				"data": [%s],
				"backgroundColor": "rgba(34, 197, 94, 0.6)",
				"borderColor": "rgb(34, 197, 94)",
				"borderWidth": 1
			}]
		},
		"options": {
			"responsive": true,
			"plugins": {
				"legend": {
					"display": true,
					"position": "top",
					"labels": {
						"color": "rgb(255, 255, 255)",
						"font": {
							"size": 12,
							"weight": "bold"
						}
					}
				},
				"title": {
					"display": true,
					"text": "90-Day Trading Volume",
					"color": "rgb(255, 255, 255)",
					"font": {
						"size": 14,
						"weight": "bold"
					}
				}
			},
			"scales": {
				"x": {
					"display": true,
					"title": {
						"display": true,
						"text": "Day of Month",
						"color": "rgb(255, 255, 255)",
						"font": {
							"size": 11,
							"weight": "bold"
						}
					},
					"ticks": {
						"color": "rgb(200, 200, 200)",
						"font": {
							"size": 9
						},
						"maxRotation": 0,
						"autoSkip": false
					},
					"grid": {
						"color": "rgba(255, 255, 255, 0.1)"
					}
				},
				"y": {
					"display": true,
					"title": {
						"display": true,
						"text": "Volume (Billions USD)",
						"color": "rgb(255, 255, 255)",
						"font": {
							"size": 11,
							"weight": "bold"
						}
					},
					"ticks": {
						"color": "rgb(200, 200, 200)",
						"font": {
							"size": 10
						},
						"callback": "function(value) { return '$' + value.toFixed(2) + 'B'; }"
					},
					"grid": {
						"color": "rgba(255, 255, 255, 0.1)"
					}
				}
			}
		}
	}`,
		strings.Join(dateLabels, ","),
		historicalData.Symbol,
		strings.Join(volumeValues, ","))

	encodedChart := url.QueryEscape(chartConfig)
	return fmt.Sprintf("https://quickchart.io/chart?c=%s&width=800&height=400&backgroundColor=%%232D3748", encodedChart)
}
