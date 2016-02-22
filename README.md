# data-analytics
Data analytics and curve fitting library in go.  
The library provides a Series structure that can be initialized with a [][]float64 containing x/y data.  
Once initialized, functions can be applied to create new Series, or curve fitting can be applied to
interpolate or (depending on the algorithm) extrapolate data points.

##Data Manipulation Functions  
Smoother - Iterative noise removal algorithm.  
Pixelize - Quantization function.  
ITrend - John Ehlers instantaneous trend (iTrend) indicator  
MA - Moving average  
EMA - Exponential moving average  
LWMA - Linear weighted moving average  
TrendChanges - Apex for peaks and troughs for smoothed data.  
Move - Transpose a series  

##Data Splicing and Combining Functions  
RecentTrends - Splices smoothed data into multiple series, each describing a trend.  
Last - Extracts a copy of the last n points from the end of a series.  
From - Extracts a copy of data points starting from an arbitrary x value.  
Append - Joins two series together to forma new series.

##Misc Functions
ToValueArray - Extracts a 1D array of y-values  

##Curve fit types
Linear  
Logarithmic  
Exponential  
Power  
Polynomial (n-order)  
Gaussian  
Parabolic  

