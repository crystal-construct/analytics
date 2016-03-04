package analytics

import (
	"testing"
)

func TestNewSeries(t *testing.T) {
	s := NewSeries()
	if s.Len != 0 {
		t.Error("New series length was not 0")
	}
}

func TestNewSeriesFrom(t *testing.T) {
	x, y := testdata[0][0], testdata[0][1]

	s := NewSeriesFrom(x, y)

	if s.Len != len(x) {
		t.Error("New series length was not", len(x))
	}
}

func TestSeriesFunctions(t *testing.T) {
	x, y := testdata[0][0], testdata[0][1]

	t.Log("Testing create.")
	s := NewSeriesFrom(x, y)
	originalLength := len(x)
	StatsCheck(t, s, 9, 1, 5, originalLength)

	t.Log("Testing add new maximum.")
	s.Add(5, 12)
	StatsCheck(t, s, 12, 1, 6.166666666666667, originalLength+1)

	t.Log("Testing add new minimum.")
	s.Add(6, 0)
	StatsCheck(t, s, 12, 0, 5.285714285714286, originalLength+2)

	t.Log("Testing slice.")
	s2 := s.Slice(2, 5)
	StatsCheck(t, s2, 9, 5, 7, 3)

	t.Log("Testing append.")
	moredata := NewSeries()
	moredata.Add(6, 9)
	moredata.Add(7, 13)
	s3 := s2.Append(moredata)
	StatsCheck(t, s3, 13, 5, 8.6, 5)

	t.Log("Testing search.")
	ord := s3.SearchX(6)
	x1, y1 := s3.Point(ord)

	if ord != 4 {
		t.Error("Search function returned", ord, ", should be 4")
	}

	if x1 != 7 {
		t.Error("X is", x1, ", should be 7")
	}

	if y1 != 13 {
		t.Error("Y is", y1, ", should be 13")
	}

}

func StatsCheck(t *testing.T, s *Series, max float64, min float64, mean float64, length int) {
	if s.Len != length {
		t.Error("Add length is incorrect. Was", s.Len, ", should be", length)
	}

	if s.Max != max {
		t.Error("Maximum recalculated incorrectly. Was", s.Max, ", should be", max)
	}

	if s.Min != min {
		t.Error("Minimum recalculated incorrectly. Was", s.Min, ", should be", min)
	}

	if s.Mean != mean {
		t.Error("Mean recalculated incorrectly. Was", s.Mean, ", should be", mean)
	}
}
