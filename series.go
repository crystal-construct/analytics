/*
Package analytics implements a library for the manipulation of x/y data series.
*/
package analytics

import "fmt"

type Series struct {
	x         []float64
	y         []float64
	Max       float64
	Min       float64
	Mean      float64
	Len       int
	sum       float64
	seriesCap int
}

//Create a new series, and initialize it with a blank backing store
func NewSeries() *Series {
	ts := &Series{}
	ts.Clear()
	return ts
}

//Create a new series from a slice of float64 slices
func NewSeriesFrom(x []float64, y []float64) *Series {
	ts := &Series{}
	ts.x = x
	ts.y = y
	ts.applyCap()
	ts.UpdateStats()
	return ts
}

//Clears the series, and initializes it with a blank backing store
func (ts *Series) Clear() {
	ts.x = []float64{}
	ts.y = []float64{}
	ts.UpdateStats()
}

//Convert the data to a 1D array
func (ts *Series) ToArrays() (x []float64, y []float64) {
	x = ts.x
	y = ts.y
	return
}

//Update stats: Min, Max, Mean, Sum
func (ts *Series) UpdateStats() {
	ts.Len = len(ts.y)
	if ts.Len == 0 {
		ts.Min = 0
		ts.Mean = 0
		ts.Max = 0
		ts.sum = 0
		return
	}

	ts.Min = ts.y[0]
	ts.Max = ts.y[0]
	ts.sum = 0
	for _, j := range ts.y {
		ts.minmax(j)
		ts.sum += j
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

//Assign a new slice to the series, and initialize it.
func (ts *Series) Use(x []float64, y []float64) {
	ts.x = x
	ts.y = y
	ts.applyCap()
	ts.UpdateStats()
}

//Add a new value to the end of the series
func (ts *Series) Add(x float64, y float64) {

	ts.x = append(ts.x, x)
	ts.y = append(ts.y, y)
	ts.sum += y
	ts.Len++
	if ts.Len == 1 {
		ts.Min = y
		ts.Max = y
	} else {
		ts.minmax(y)
	}
	ts.applyCap()
	ts.calculatemean()
}

//Set a value at the ordinal position in the series.
//Altering values that are outside the max/min of the existing data
//will cause the statistics for min, mean and max to be recalculated.
func (ts *Series) Set(ordinal int, value float64) {
	oldValue := ts.y[ordinal]
	ts.y[ordinal] = value
	ts.minmax(value)
	ts.sum -= oldValue
	ts.sum += value

	if ordinal == 0 && ts.Len == 1 {
		ts.Min = value
		ts.Max = value
		ts.Mean = value
		return
	}
	ts.minmax(value)
	ts.calculatemean()
}

//Creates a new series containing the last n values.
func (ts *Series) Last(n int) *Series {
	x := make([]float64, n)
	y := make([]float64, n)
	copy(x, ts.x[ts.Len-n:])
	copy(y, ts.y[ts.Len-n:])
	return NewSeriesFrom(x, y)
}

//Extracts the values from the series as a 1 dimensional slice
func (ts *Series) ToValues(length int, offset int) (x []float64, y []float64) {
	x = ts.x[ts.Len-length-offset : ts.Len-offset]
	y = ts.y[ts.Len-length-offset : ts.Len-offset]
	return
}

//Shifts a dataset on the x and y axes
func (ts *Series) ApplyOffset(x float64, y float64) *Series {
	newx := make([]float64, ts.Len)
	newy := make([]float64, ts.Len)
	for i := range ts.x {
		newx[i] = ts.x[i] + x
		newy[i] = ts.y[i] + y
	}
	return NewSeriesFrom(newx, newy)
}

//Creates a new series starting from the earliset point a particular time
func (ts *Series) From(time float64) *Series {
	newts := &Series{}
	pos := ts.SearchX(time)
	if pos == -1 {
		newts.Clear()
	} else {
		newts.Use(ts.x[pos:], ts.y[pos:])
	}
	return newts
}

//Appends one series to another
func (ts *Series) Append(toAdd *Series) *Series {
	newx := make([]float64, 0, ts.Len+toAdd.Len)
	newy := make([]float64, 0, ts.Len+toAdd.Len)
	newx = append(ts.x, toAdd.x...)
	newy = append(ts.y, toAdd.y...)
	return NewSeriesFrom(newx, newy)
}

//Applies two functions.  The map function recieves a series representing a period,
//and returns a []float64.  The reduce function takes the aggregated results and
//translates them into a series.
func (ts *Series) MapReduce(mapFunction func(*Series) (float64, float64), reduceFunction func([]float64, []float64) *Series, periodLength float64, numberOfPeriods int) *Series {
	p := 0
	maxx := ts.x[ts.Len-1]
	start := maxx - (periodLength * float64(numberOfPeriods))
	mappedx := make([]float64, numberOfPeriods)
	mappedy := make([]float64, numberOfPeriods)
	for start <= maxx && p < numberOfPeriods {
		pos := ts.SearchX(start)
		end := ts.SearchX(start + periodLength)
		start += periodLength
		if pos == -1 {
			continue
		}
		slice := ts.Slice(pos, end)
		if slice.Len == 0 {
			continue
		}
		mappedx[p], mappedy[p] = mapFunction(slice)
		p++
	}
	return reduceFunction(mappedx[:p], mappedy[:p])
}

//Slices a series - this is equivalent to go's slice
func (ts *Series) Slice(start int, end int) *Series {
	return NewSeriesFrom(ts.x[start:end], ts.y[start:end])
}

//Uses binary search to find the earliest ordinal occurance of a x value.
func (ts *Series) SearchX(value float64) int {
	xdata := ts.x
	if xdata[0] > value {
		return -1
	}
	if xdata[ts.Len-1] < value {
		return ts.Len
	}

	i, j := 0, ts.Len
	for i < j {
		h := i + (j-i)/2
		if !(xdata[h] > value) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func (ts *Series) applyCap() {
	if ts.seriesCap == 0 {
		return
	}
	if ts.Len > ts.seriesCap {
		updateRequired := false
		overflow := ts.Len - ts.seriesCap
		var sum float64
		for i := 0; i < overflow; i++ {
			sum -= ts.y[i]
			if ts.y[i] == ts.Min || ts.y[i] == ts.Max {
				updateRequired = true
			}
		}

		ts.x = ts.x[overflow:]
		ts.y = ts.y[overflow:]
		if !updateRequired {
			ts.sum -= sum
		} else {
			ts.UpdateStats()
		}
		ts.Len = len(ts.x)
		//Garbage collect the slice if backing store %10 > len
		if cap(ts.x) > ts.Len+ts.Len/10 {
			newtsx := make([]float64, ts.Len)
			newtsy := make([]float64, ts.Len)
			copy(newtsx, ts.x)
			copy(newtsy, ts.y)
			ts.x = newtsx
			ts.y = newtsy
		}
	}
}

func (ts *Series) SetCap(n int) {
	if ts.Len > 0 {
		panic(fmt.Errorf("Capacity cannot be set on a series with length > 0"))
	}
	ts.seriesCap = n
	ts.x = make([]float64, 0, n)
	ts.y = make([]float64, 0, n)
}

func (ts *Series) Point(ordinal int) (x float64, y float64) {
	return ts.x[ordinal], ts.y[ordinal]
}
