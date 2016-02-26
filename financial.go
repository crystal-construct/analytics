package analytics

func (ts *Series) CommonChannelIndex(periodLength float64, numberOfPeriods int) *Series {
	var constant float64 = 0.015
	j := ts.MapReduce(
		func(t *Series) (x float64, y float64) {
			x = t.x[t.Len-1]
			y = (t.Max + t.Min + t.y[t.Len-1]) / 3
			return
		},
		func(xdata []float64, ydata []float64) *Series {
			dataseries := NewSeriesFrom(xdata, ydata)
			sma := dataseries.Ma(3)
			meandev := dataseries.MeanDev()
			ny := make([]float64, ts.Len)
			nx := make([]float64, ts.Len)
			copy(nx, ts.x)
			for i := range xdata {
				ny[i] = (ydata[i] - sma.y[i]) / (constant * meandev)
			}
			s := NewSeriesFrom(nx, ny)
			return s
		},
		periodLength, numberOfPeriods)
	return j
}
