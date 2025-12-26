package analytics

import (
	"highload/internal/model"
	"math"
	"sync"
)

type Analyzer struct {
	windowSize int
	threshold  float64
	metrics    []model.Metric
	mutex      sync.Mutex
}

func NewAnalyzer(windowSize int, threshold float64) *Analyzer {
	return &Analyzer{
		windowSize: windowSize,
		threshold:  threshold,
		metrics:    make([]model.Metric, 0, windowSize),
	}
}

func (a *Analyzer) Update(m model.Metric) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.metrics = append(a.metrics, m)
	if len(a.metrics) > a.windowSize {
		a.metrics = a.metrics[len(a.metrics)-a.windowSize:]
	}
}

func (a *Analyzer) GetStats() (avgRPS float64, anomalies int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if len(a.metrics) == 0 {
		return 0, 0
	}

	// Rolling Average
	var sum float64
	for _, m := range a.metrics {
		sum += m.RPS
	}
	avgRPS = sum / float64(len(a.metrics))

	// Z-Score
	variance := 0.0
	for _, m := range a.metrics {
		deviation := m.RPS - avgRPS
		variance += deviation * deviation
	}
	variance /= float64(len(a.metrics))
	stdDev := math.Sqrt(variance)

	anomalies = 0
	for _, m := range a.metrics {
		if stdDev != 0 {
			z := math.Abs(m.RPS-avgRPS) / stdDev
			if z > a.threshold {
				anomalies++
			}
		}
	}

	return avgRPS, anomalies
}
