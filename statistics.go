package analytics

import (
	"math"
)

// Iterative Noise Removal
func (ts *Series) Smoother(period int) *Series {
	var l int = len(ts.Data)
	// Reset the buffer
	buffer := make([][]float64, len(ts.Data))
	for i, j := range ts.Data {
		buffer[i] = []float64{j[0], j[1]}
	}

	for j := 0; j < period; j++ {
		for i := 3; i < l; i++ {
			buffer[i-1] = []float64{
				buffer[i-1][0],
				(buffer[i-2][1] + buffer[i][1]) / 2,
			}
		}
	}
	t := &Series{}
	t.Use(buffer)
	return t
}

// Pixelize - Domain reduction
func (ts *Series) Pixelize(grid int) *Series {

	// Calculate the grid values
	var min = ts.Min
	var max = ts.Max
	var tile = (max - min) / float64(grid)
	buffer := make([][]float64, len(ts.Data))
	for i, datapoint := range ts.Data {
		buffer[i] = []float64{datapoint[0], round(datapoint[1]/tile) * tile}
	}
	t := &Series{}
	t.Use(buffer)
	return t
}

// DSL, iTrend
func (ts *Series) ITrend(alpha float64) (itrendSeries *Series) {
	// By Ehler
	// http://www.davenewberg.com/Trading/TS_Code/Ehlers_Indicators/iTrend_Ind.html
	l := len(ts.Data)

	var buffer = make([][]float64, 3)
	var trigger = make([][]float64, 3)
	for i, j := range ts.Data[0:3] {
		buffer[i] = []float64{j[0], j[1]}
	}
	for i, j := range ts.Data[0:3] {
		trigger[i] = []float64{j[0], j[1]}
	}

	for i := 3; i < l; i++ {
		buffer = append(buffer,
			[]float64{
				ts.Data[i][0],
				(alpha-(alpha*alpha)/4)*ts.Data[i][1] + (0.5 * (alpha * alpha) * ts.Data[i-1][1]) - (alpha-0.75*(alpha*alpha))*ts.Data[i-2][1] + 2*(1-alpha)*buffer[i-1][1] - (1-alpha)*(1-alpha)*buffer[i-2][1],
			})

		trigger = append(buffer,
			[]float64{ts.Data[i][0],
				2*buffer[i][1] - buffer[i-2][1],
			})
	}
	t := &Series{}
	t.Use(buffer)
	u := &Series{}
	u.Use(trigger)
	return t
}

func (ts *Series) StDev() float64 {
	var sdsum float64 = 0
	for _, j := range ts.Data {
		sdsum += math.Pow(j[1]-ts.Mean, 2)
	}
	return math.Sqrt(sdsum / float64(len(ts.Data)))
}

// Moving Average
func (ts *Series) Ma(period int) *Series {
	var l int = len(ts.Data)
	var sum float64 = 0
	var buffer = make([][]float64, period)
	for i, j := range ts.Data[0:period] {
		buffer[i] = []float64{j[0], j[1]}
	}
	for i := period; i < l; i++ {
		sum = 0
		for j := period; j > 0; j-- {
			sum += ts.Data[i-j][1]
		}
		buffer = append(buffer, []float64{ts.Data[i][0], sum / float64(period)})
	}
	t := &Series{}
	t.Use(buffer)
	return t
}

func (ts *Series) Ema(period int) *Series {

	var l int = len(ts.Data)

	var buffer = make([][]float64, period)
	for i, j := range ts.Data[0:period] {
		buffer[i] = []float64{j[0], j[1]}
	}
	var m float64 = 2 / (float64(period) + 1) // Multiplier

	for i := period; i < l; i++ {
		buffer = append(buffer, []float64{ts.Data[i][0], (ts.Data[i][1]-ts.Data[i-1][1])*m + ts.Data[i-1][1]})
	}
	t := &Series{}
	t.Use(buffer)
	return t
}

func (ts *Series) Lwma(period int) *Series {

	var l int = len(ts.Data)
	var sum float64 = 0

	var buffer = make([][]float64, period)
	for i, j := range ts.Data[0:period] {
		buffer[i] = []float64{j[0], j[1]}
	}
	for i := period; i < l; i++ {
		sum = 0
		var n int = 0
		for j := period; j > 0; j-- {
			sum += ts.Data[i-j][1] * float64(j)
			n += j
		}
		buffer = append(buffer, []float64{ts.Data[i][0], sum / float64(n)})
	}
	t := &Series{}
	t.Use(buffer)
	return t
}

func (ts *Series) RecentTrends(n int) []*Series {
	ret := []*Series{}
	var marker int = ts.Len
	var trend, oldTrend int = 0, 0
	var found int = 0
	for i := ts.Len - 2; i > -1; i-- {
		if ts.Data[i][1] > ts.Data[i+1][1] {
			trend = -1
		} else if ts.Data[i][1] < ts.Data[i+1][1] {
			trend = 1
		}
		if (trend != oldTrend && oldTrend != 0) || i == 0 {
			newts := &Series{}
			newdata := make([][]float64, marker-i-1)
			for j := range newdata {
				newdata[j] = []float64{ts.Data[i+1+j][0], ts.Data[i+1+j][1]}
			}
			newts.Use(newdata)
			ret = append(ret, newts)
			marker = i + 1
			found++
			if found == n {
				return ret
			}
		}
		oldTrend = trend
	}
	return nil
}

func (ts *Series) TrendChanges() *Series {
	buffer := make([][]float64, 0)
	l := len(ts.Data)
	dirup := ts.Data[1][1] > ts.Data[0][1]
	for i := 1; i < l; i++ {
		newdir := ts.Data[i][1] > ts.Data[i-1][1]
		if newdir != dirup {
			buffer = append(buffer, []float64{ts.Data[i-1][0], ts.Data[i-1][1]})
			dirup = newdir
		}
	}
	t := &Series{}
	t.Use(buffer)
	return t
}
