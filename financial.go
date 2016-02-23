package analytics

func (ts *Series) CommonChannelIndex(periodLength float64, numberOfPeriods int) *Series {
	var constant float64 = 0.015
	j := ts.MapReduce(
		func(t *Series) []float64 {
			return []float64{t.Data[t.Len-1][0], (t.Max + t.Min + t.Data[t.Len-1][1]) / 3}
		},
		func(data [][]float64) *Series {
			dataseries := NewSeriesFrom(data)
			sma := dataseries.Ma(3)
			meandev := dataseries.MeanDev()
			s := NewSeries()
			for i := range data {
				s.Add(data[i][0], (data[i][1]-sma.Data[i][1])/(constant*meandev))
			}
			return s
		},
		periodLength, numberOfPeriods)
	return j
}
