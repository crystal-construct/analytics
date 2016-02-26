package analytics

import (
	"math"
)

//Iterative Noise Removal
func (ts *Series) Smoother(period int) *Series {
	var l int = len(ts.x)

	bufferx := make([]float64, l)
	buffery := make([]float64, l)
	copy(bufferx, ts.x)
	copy(buffery, ts.y)

	for j := 0; j < period; j++ {
		for i := 3; i < l; i++ {
			buffery[i-1] = (buffery[i-2] + buffery[i]) / 2
		}
	}
	return NewSeriesFrom(bufferx, buffery)
}

//Quantization
func (ts *Series) Quantize(grid int) *Series {

	var min = ts.Min
	var max = ts.Max
	var resolution = (max - min) / float64(grid)
	bufferx := make([]float64, ts.Len)
	copy(bufferx, ts.x)
	buffery := make([]float64, ts.Len)
	for i := range ts.y {
		buffery[i] = round(ts.y[1]/resolution) * resolution
	}
	return NewSeriesFrom(bufferx, buffery)
}

//iTrend
func (ts *Series) ITrend(alpha float64) (itrendSeries *Series) {
	l := ts.Len

	var bufferx = make([]float64, ts.Len)
	copy(bufferx, ts.x)
	var buffery = make([]float64, ts.Len)
	copy(buffery, ts.y[:3])
	var triggery = make([]float64, ts.Len)
	copy(triggery, ts.y[:3])

	for i := 3; i < l; i++ {
		y := ts.y[i]
		y1 := ts.y[i-1]
		y2 := ts.y[i-2]
		buffery[i] = (alpha-(alpha*alpha)/4)*y + (0.5 * (alpha * alpha) * y1) - (alpha-0.75*(alpha*alpha))*y2 + 2*(1-alpha)*y1 - (1-alpha)*(1-alpha)*y2
		triggery[i] = 2*y1 - y2
	}
	t := NewSeriesFrom(bufferx, buffery)
	//u := NewSeriesFrom(trigger)
	return t
}

// Standard deviation
func (ts *Series) StDev() float64 {
	if ts.Len == 0 {
		return 0
	}
	var sdsum float64 = 0
	for i := range ts.x {
		sdsum += math.Pow(ts.y[i]-ts.Mean, 2)
	}
	return math.Sqrt(sdsum / float64(ts.Len))
}

// Mean deviation
func (ts *Series) MeanDev() float64 {
	if ts.Len == 0 {
		return 0
	}
	var mdsum float64 = 0
	for j := range ts.y {
		mdsum += math.Abs(ts.y[j] - ts.Mean)
	}
	return mdsum / float64(ts.Len)
}

// Moving Average
func (ts *Series) Ma(period int) *Series {
	var l int = ts.Len
	var sum float64 = 0
	var bufferx = make([]float64, period, ts.Len)

	var buffery = make([]float64, period, ts.Len)
	copy(bufferx, ts.x)
	copy(buffery, ts.y[:period])
	for i := period; i < l; i++ {
		sum = 0
		for j := period; j > 0; j-- {
			sum += ts.y[i-j]
		}
		buffery = append(buffery, sum/float64(period))
	}
	return NewSeriesFrom(bufferx, buffery)
}

//Exponential moving average
func (ts *Series) Ema(period int) *Series {

	var l int = ts.Len

	var bufferx = make([]float64, period, ts.Len)
	copy(bufferx, ts.x)
	var buffery = make([]float64, period, ts.Len)
	copy(buffery, ts.y[:period])
	var m float64 = 2 / (float64(period) + 1) // Multiplier

	for i := period; i < l; i++ {
		buffery[i] = (ts.y[i]-ts.y[i-1])*m + ts.y[i-1]
	}
	return NewSeriesFrom(bufferx, buffery)
}

//Linear weighted moving average
func (ts *Series) Lwma(period int) *Series {

	var l int = ts.Len
	var sum float64 = 0

	var bufferx = make([]float64, period, ts.Len)
	copy(bufferx, ts.x)
	var buffery = make([]float64, period, ts.Len)
	copy(buffery, ts.y[:period])
	for i := period; i < l; i++ {
		sum = 0
		var n int = 0
		for j := period; j > 0; j-- {
			sum += ts.y[i-j] * float64(j)
			n += j
		}
		buffery = append(buffery, sum/float64(n))
	}
	return NewSeriesFrom(bufferx, buffery)
}

//Recent trends
func (ts *Series) RecentTrends(n int) []*Series {
	ret := []*Series{}
	datax := ts.x
	datay := ts.y
	var marker int = ts.Len
	var trend, oldTrend int = 0, 0
	var found int = 0
	for i := ts.Len - 2; i > -1; i-- {
		if datay[i] > datay[i+1] {
			trend = -1
		} else if datay[i] < datay[i+1] {
			trend = 1
		}
		if (trend != oldTrend && oldTrend != 0) || i == 0 {
			newx := make([]float64, marker-i-1)
			newy := make([]float64, marker-i-1)
			for j := range newx {
				newx[j] = datax[i+1+j]
				newy[j] = datay[i+1+j]
			}
			newts := NewSeriesFrom(newx, newy)
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

//Peak and trough data points
func (ts *Series) TrendChanges() *Series {
	bufferx := make([]float64, 500)
	buffery := make([]float64, 500)
	l := ts.Len
	dirup := ts.y[1] > ts.y[0]
	for i := 1; i < l; i++ {
		newdir := ts.y[i] > ts.y[i-1]
		if newdir != dirup {
			bufferx = append(bufferx, ts.x[i-1])
			buffery = append(buffery, ts.y[i-1])
			dirup = newdir
		}
	}
	return NewSeriesFrom(bufferx, buffery)
}
