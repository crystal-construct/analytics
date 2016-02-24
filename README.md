# analytics
Analytics and curve fitting library in go.  
The library provides a Series structure that can be initialized with a [][]float64 containing x/y data.  
Once initialized, functions can be applied to create new Series, or curve fitting can be applied to
interpolate or (depending on the algorithm) extrapolate data points.

Simple example usage:  
```go
	//Define a small dataset
	d := [][]float64{
		{1, 8}, {2, 6}, {3, 5}, {4, 5}, {5, 4}, {6, 3}, {7, 3}, {8, 2},
	}

	//Use the dataset as a series
	series1 := analytics.NewSeriesFrom(d)

	//Get a 3rd Order Polynomial fit for the series
	fit := series1.FitPolynomial(3)

	//Display two interpolated/extrapolated points
	fmt.Println(analytics.Extrapolate(fit, 4))
	fmt.Println(analytics.Extrapolate(fit, 9))

	//Create smoothed version of the dataset
	series2 := series1.Smoother(3)

	//Display the underlying smoothed values
	fmt.Println(series2.Data)
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
ToValueArray - Extracts a 1D slice of y-values
ToValues - Extracts a 1D slice of y-values with an offset from the last datapoint


##Curve fit types
Linear  
Logarithmic  
Exponential  
Power  
Polynomial (n-order)  
Gaussian  
Parabolic  

##Todo
Correct and clean up attribution    
Fix up Gaussian fit  
Fix up Parabolic fit  
Implement Cubic Spline  
Implement 4PL  
Cache sums for curve fit  
Error handling  
File routines to load and save / export in gnuplot format  
