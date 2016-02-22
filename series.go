package timeseries

type Series struct {
	Data [][]float64
	Max  float64
	Min  float64
	Mean float64
	Len  int
	sum  float64
}

func NewSeries() *Series {
	ts := &Series{}
	ts.Clear()
	return ts
}

func (ts *Series) Clear() {
	ts.Data = [][]float64{}
	ts.UpdateStats()
}

// Convert the data to a 1D array
func (ts *Series) ToArray() []float64 {
	outArray := make([]float64, len(ts.Data))
	for i, j := range ts.Data {
		outArray[i] = j[1]
	}
	return outArray
}

// Stats: Min, Max, Mean
func (ts *Series) UpdateStats() {
	ts.Len = len(ts.Data)
	if ts.Len == 0 {
		ts.Min = 0
		ts.Mean = 0
		ts.Max = 0
		ts.sum = 0
		return
	}

	ts.Min = ts.Data[0][1]
	ts.Max = ts.Data[0][1]
	ts.sum = 0
	for _, j := range ts.Data {
		ts.minmax(j[1])
		ts.sum += j[1]
	}
	ts.calculatemean()
}

func (ts *Series) calculatemean() {
	ts.Mean = ts.sum / float64(ts.Len)
}

func (ts *Series) minmax(value float64) {
	if ts.Min > value {
		ts.Min = value
	}
	if ts.Max < value {
		ts.Max = value
	}
}

func (ts *Series) Use(data [][]float64) {
	ts.Data = data
	ts.UpdateStats()
}

func (ts *Series) Add(time float64, value float64) {
	ts.Data = append(ts.Data, []float64{time, value})
	ts.sum += value
	ts.Len += 1
	if ts.Len == 1 {
		ts.Min = value
		ts.Max = value
	} else {
		ts.minmax(value)
	}
	ts.calculatemean()
}

func (ts *Series) Set(pos int, value float64) {
	oldValue := ts.Data[pos][1]
	ts.Data[pos][1] = value
	ts.minmax(value)
	ts.sum -= oldValue
	ts.sum += value

	if pos == 0 && ts.Len == 1 {
		ts.Min = value
		ts.Max = value
		ts.Mean = value
		return
	}
	if oldValue == ts.Max || oldValue == ts.Min {
		ts.UpdateStats()
	} else {
		ts.minmax(value)
		ts.calculatemean()
	}
}

func (ts *Series) Last(n int) *Series {
	data := make([][]float64, n)
	for i, j := range ts.Data[ts.Len-n:] {
		data[i] = []float64{j[0], j[1]}
	}
	newts := &Series{}
	newts.Use(data)
	return newts
}

func (ts *Series) FromVolumeTrades(data []VolumeTrade) {
	ts.Data = make([][]float64, len(data))
	for i, j := range data {
		ts.Data[i] = []float64{float64(j.Time), j.Price, j.Volume}
	}
}

func (ts *Series) ToValueArray(length int, offset int) []float64 {
	slice := ts.Data[ts.Len-length-offset : ts.Len-offset]
	ret := make([]float64, length)
	for i := range ret {
		ret[i] = slice[i][1]
	}
	return ret
}

func (ts *Series) Move(xoffset float64, yoffset float64) *Series {
	newdata := make([][]float64, ts.Len)
	for i, j := range ts.Data {
		newdata[i] = []float64{j[0] + xoffset, j[1] + yoffset}
	}
	newts := &Series{}
	newts.Use(newdata)
	return newts
}

func (ts *Series) From(time float64) *Series {
	var i int
	newts := &Series{}
	for i = ts.Len - 1; i > 0; i-- {

		if ts.Data[i-1][0] < time {
			newts.Use(ts.Data[i:])
			return newts
		}
	}
	newts.Use(ts.Data)
	return newts
}

func (ts *Series) Append(toAdd *Series) *Series {
	ctr := 0
	newdata := make([][]float64, ts.Len+toAdd.Len)
	for i := range ts.Data {
		newdata[ctr] = []float64{ts.Data[i][0], ts.Data[i][1]}
		ctr++
	}
	for i := range toAdd.Data {
		newdata[ctr] = []float64{toAdd.Data[i][0], toAdd.Data[i][1]}
		ctr++
	}
	newts := &Series{}
	newts.Use(newdata)
	return newts
}
