package analytics

/*
Package analytics implements a library for the manipulation of x/y data series.


*/

type Series struct {
	Data [][]float64
	Max  float64
	Min  float64
	Mean float64
	Len  int
	sum  float64
}

//Creates a new series, and initializes it with a blank backing store
func NewSeries() *Series {
	ts := &Series{}
	ts.Clear()
	return ts
}

//Creates a new series from a slice of float64 slices
func NewSeriesFrom(data [][]float64) *Series {
	ts := &Series{}
	ts.Use(data)
	return ts
}

//Clears the series, and initializes it with a blank backing store
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

//Sets a value at the ordinal position in the series.
//Altering values that are outside the max/min of the existing data
//will cause the statistics for min, mean and max to be recalculated.
func (ts *Series) Set(ordinal int, value float64) {
	oldValue := ts.Data[ordinal][1]
	ts.Data[ordinal][1] = value
	ts.minmax(value)
	ts.sum -= oldValue
	ts.sum += value

	if ordinal == 0 && ts.Len == 1 {
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

//Creates a new series containing the last n values.
func (ts *Series) Last(n int) *Series {
	data := make([][]float64, n)
	for i, j := range ts.Data[ts.Len-n:] {
		data[i] = []float64{j[0], j[1]}
	}
	return NewSeriesFrom(data)
}

//Extracts the values from the series as a 1 dimensional slice
func (ts *Series) ToValues(length int, offset int) []float64 {
	slice := ts.Data[ts.Len-length-offset : ts.Len-offset]
	ret := make([]float64, length)
	for i := range ret {
		ret[i] = slice[i][1]
	}
	return ret
}

//Shifts a dataset on the x and y axes
func (ts *Series) ApplyOffset(x float64, y float64) *Series {
	newdata := make([][]float64, ts.Len)
	for i, j := range ts.Data {
		newdata[i] = []float64{j[0] + x, j[1] + y}
	}
	return NewSeriesFrom(newdata)
}

//Creates a new series starting from the earliset point a particular time
func (ts *Series) From(time float64) *Series {
	newts := &Series{}
	pos := ts.SearchX(time)
	if pos == -1 {
		newts.Clear()
	} else {
		newts.Use(ts.Data[pos:])
	}
	return newts
}

//Appends one series to another
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
	return NewSeriesFrom(newdata)
}

//Applies two functions.  The map function recieves a series representing a period,
//and returns a []float64.  The reduce function takes the aggregated results and
//translates them into a series.
func (ts *Series) MapReduce(mapFunction func(*Series) []float64, reduceFunction func([][]float64) *Series, periodLength float64, numberOfPeriods int) *Series {
	p := 0
	maxx := ts.Data[ts.Len-1][0]
	start := maxx - (periodLength * float64(numberOfPeriods))
	mapped := make([][]float64, 0, numberOfPeriods)
	for start <= maxx && p < numberOfPeriods {
		pos := ts.SearchX(start)
		end := ts.SearchX(start + periodLength)
		start += periodLength
		p++
		if pos == -1 {
			continue
		}
		mapped = append(mapped, mapFunction(ts.Slice(pos, end)))
	}

	return reduceFunction(mapped)
}

//Slices a series - this is equivalent to go's slice
func (ts *Series) Slice(start int, end int) *Series {
	return NewSeriesFrom(ts.Data[start:end])
}

//Uses binary search to find the earliest ordinal occurance of a x value.
func (ts *Series) SearchX(value float64) int {
	data := ts.Data
	if data[0][0] > value {
		return -1
	}
	if data[ts.Len-1][0] < value {
		return ts.Len
	}

	i, j := 0, ts.Len
	for i < j {
		h := i + (j-i)/2
		if !(data[h][0] > value) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
