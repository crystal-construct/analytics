# analytics
Analytics and curve fitting library in go.  
The library provides a Series structure that can be initialized with a [][]float64 containing x/y data.  
Once initialized, functions can be applied to create new Series, or curve fitting can be applied to
interpolate or (depending on the algorithm) extrapolate data points.

Simple example usage:  
```go
	//Define a small dataset
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{8, 6, 5, 5, 4, 3, 3, 2}
	
	//Use the dataset as a series
	series1 := analytics.NewSeriesFrom(x, y)

	//Get a 3rd Order Polynomial fit for the series
	fit := series1.FitPolynomial(3)

	//Display two interpolated/extrapolated points
	fmt.Println(analytics.Extrapolate(fit, 4))
	fmt.Println(analytics.Extrapolate(fit, 9))

	//Create smoothed version of the dataset
	series2 := series1.Smoother(3)

	//Display the underlying smoothed values
	x1, y1 := series2.ToArrays()
	fmt.Println(series2.x1)
	fmt.Println(series2.y1)
```

##General Data Manipulation Functions  
Smoother - Iterative noise removal algorithm.  
Pixelize - Quantization function.  
MA - Moving average  
EMA - Exponential moving average  
LWMA - Linear weighted moving average  
TrendChanges - Apex for peaks and troughs for smoothed data.  
ApplyOffset - Move a series on x/y axes

##Data Splicing and Combining Functions  
RecentTrends - Splices smoothed data into multiple series, each describing a trend.  
Last - Extracts a copy of the last n points from the end of a series.  
From - Extracts a copy of data points starting from an arbitrary x value.  
Append - Joins two series together to form a new series.  

##Financial Analysis Based Functions  
ITrend - John Ehlers instantaneous trend (iTrend) indicator  
CCI - Commodity Channel Index  

##Misc Functions
ToArrays - Extracts two 1D slices of values, one for x and one for y
ToValues - As ToArrays, but takes an offset from the last datapoint


##Curve fit types
Linear  
Logarithmic  
Exponential  
Power  
Polynomial (n-order)  
Gaussian  
Parabolic  

##Todo
Tests!
Correct and clean up attribution    
Fix up Gaussian fit  
Fix up Parabolic fit  
Implement Cubic Spline  
Implement 4PL  
Cache sums for curve fit  
Error handling