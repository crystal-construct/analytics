package analytics

import (
	"fmt"
	"math"
)

type FitParameters struct {
	fittype int
	values  []float64
	xoffset float64
	yoffset float64
}

const (
	FitTypeLinear = iota
	FitTypeLinearThroughOrigin
	FitTypeLogarithmic
	FitTypePower
	FitTypeExponential
	FitTypePolynomial
	FitTypeGaussian
	FitTypeParabolic
)

func (ts *Series) FitExponential() (params FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	sum := []float64{0, 0, 0, 0, 0, 0}

	for n := range xdata {
		time := xdata[n] - xoffset
		value := ydata[n] - yoffset
		sum[0] += time                           // X
		sum[1] += value                          // Y
		sum[2] += time * time * value            // XXY
		sum[3] += value * math.Log(value)        // Y Log Y
		sum[4] += time * value * math.Log(value) //YY Log Y
		sum[5] += time * value                   //XY
	}

	denominator := (sum[1]*sum[2] - sum[5]*sum[5])
	A := math.Pow(math.E, (sum[2]*sum[3]-sum[5]*sum[4])/denominator)
	B := (sum[1]*sum[4] - sum[5]*sum[3]) / denominator

	params = FitParameters{
		fittype: FitTypeExponential,
		values:  []float64{A, B},
	}
	return
}

/**
 *              N * Σ(XY) - Σ(X)
 * intercept = ---------------------
 *              N * Σ(X^2) - Σ(X)^2
 *
 * correlation = N * Σ(XY) - Σ(X) * Σ (Y) / √ (  N * Σ(X^2) - Σ(X) ) * ( N * Σ(Y^2) - Σ(Y)^2 ) ) )
 *
 */
func (ts *Series) FitLinear() (params FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	sum := []float64{0, 0, 0, 0, 0}
	N := float64(ts.Len)

	for n := range xdata {
		x := xdata[n] - xoffset
		y := ydata[n] - yoffset
		sum[0] += x     //Σ(X)
		sum[1] += y     //Σ(Y)
		sum[2] += x * x //Σ(X^2)
		sum[3] += x * y //Σ(XY)
		sum[4] += y * y //Σ(Y^2)
	}

	var gradient = (N*sum[3] - sum[0]*sum[1]) / (N*sum[2] - sum[0]*sum[0])
	var intercept = (sum[1] / N) - (gradient*sum[0])/N
	var correlation = (N*sum[3] - sum[0]*sum[1]) / math.Sqrt((N*sum[2]-sum[0]*sum[0])*(N*sum[4]-sum[1]*sum[1]))

	params = FitParameters{
		fittype: FitTypeLinear,
		values:  []float64{gradient, intercept, correlation},
	}
	return
}

func (ts *Series) FitLinearThroughOrigin() (params FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	sum := []float64{0, 0, 0, 0, 0}

	for n := 0; n < ts.Len; n++ {
		x := xdata[n] - xoffset
		y := ydata[n] - yoffset
		sum[0] += x * x //sumSqX
		sum[1] += x * y //sumXY
	}

	var gradient = sum[1] / sum[0]

	params = FitParameters{
		fittype: FitTypeLinearThroughOrigin,
		values:  []float64{gradient},
	}
	return params
}

func (ts *Series) FitLogarithmic() (params FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	var sum = []float64{0, 0, 0, 0}
	N := float64(ts.Len)

	for n := range ts.x {
		x := xdata[n] - xoffset
		y := ydata[n] - yoffset
		sum[0] += math.Log(x)
		sum[1] += y * math.Log(x)
		sum[2] += y
		sum[3] += math.Pow(math.Log(x), 2)
	}

	var B = (N*sum[1] - sum[2]*sum[0]) / (N*sum[3] - sum[0]*sum[0])
	var A = (sum[2] - B*sum[0]) / N

	params = FitParameters{
		fittype: FitTypeLogarithmic,
		values:  []float64{A, B},
	}
	return
}

func (ts *Series) FitPower() (params FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	var sum = []float64{0, 0, 0, 0}
	N := float64(ts.Len)

	for n := range xdata {
		x := xdata[n] - xoffset
		y := ydata[n] - yoffset
		sum[0] += math.Log(x)
		sum[1] += math.Log(y) * math.Log(x)
		sum[2] += math.Log(y)
		sum[3] += math.Pow(math.Log(x), 2)
	}

	var B = (N*sum[1] - sum[2]*sum[0]) / (N*sum[3] - sum[0]*sum[0])
	var A = math.Pow(math.E, (sum[2]-B*sum[0])/N)

	params = FitParameters{
		fittype: FitTypePower,
		values:  []float64{A, B},
	}
	return
}

func (ts *Series) FitPolynomial(order int) (params FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	rhs := [][]float64{}
	lhs := []float64{}
	k := order + 1
	a := float64(0)
	var b float64 = 0
	for i := 0; i < k; i++ {
		for l := range ts.x {
			a += math.Pow(xdata[l]-xoffset, float64(i)) * (ydata[l] - yoffset)
		}
		lhs = append(lhs, a)
		a = 0
		var c = []float64{}
		for j := 0; j < k; j++ {
			for l := range xdata {
				b += math.Pow(xdata[l]-xoffset, float64(i+j))
			}
			c = append(c, b)
			b = 0
		}
		rhs = append(rhs, c)
	}
	rhs = append(rhs, lhs)

	equation := gaussianElimination(rhs, k)

	params = FitParameters{
		fittype: FitTypePolynomial,
		values:  equation,
		xoffset: xoffset,
		yoffset: yoffset,
	}
	return
}

func Extrapolate(params FitParameters, x float64) float64 {
	switch params.fittype {
	case FitTypePolynomial:
		var answer float64 = 0
		for w := 0; w < len(params.values); w++ {
			answer += params.values[w] * math.Pow((x-params.xoffset), float64(w))
		}
		return answer + params.yoffset
	case FitTypeExponential:
		A := params.values[0]
		B := params.values[1]
		return A*math.Pow(math.E, B*(x-params.xoffset)) + params.yoffset
	case FitTypeLinear:
		gradient := params.values[0]
		intercept := params.values[1]
		return (x-params.xoffset)*gradient + intercept
	case FitTypeLinearThroughOrigin:
		Gradient := params.values[0]
		return x * Gradient
	case FitTypeLogarithmic:
		A := params.values[0]
		B := params.values[1]
		return A + B*math.Log((x-params.xoffset)) + params.yoffset
	case FitTypePower:
		A := params.values[0]
		B := params.values[1]
		return A*math.Pow((x-params.xoffset), B) + params.yoffset
	case FitTypeGaussian:
		A := params.values[0]
		B := params.values[1]
		C := params.values[2]
		return (A*math.Pow((x-params.xoffset), 2) + B*(x-params.xoffset) + C) + params.yoffset
	case FitTypeParabolic:
		Height := params.values[0]
		Position := params.values[1]
		Width := params.values[2]
		return (Height * math.Exp(-1*math.Pow(((x-params.xoffset)-Position)/(0.6006*Width), 2))) + params.yoffset
	}

	panic(fmt.Errorf("No Fit Available"))
}

func (ts *Series) FitGaussianParabolic() (params []FitParameters) {
	xoffset := ts.x[0] - 1
	yoffset := ts.Min - 1
	xdata := ts.x
	ydata := ts.y
	var n float64 = float64(ts.Len)
	var sumx, sumy, sumxy, sumx2, sumx3, sumx4, sumx2y float64
	for i := range ts.x {
		x := xdata[i] - xoffset
		y := ydata[i] - yoffset
		lny := math.Log(y)
		sumx += x
		sumy += lny
		sumxy += x * lny
		x2 := x * x
		sumx2 += x2
		sumx3 += math.Pow(x, 3)
		sumx4 += math.Pow(x, 4)
		sumx2y += lny * math.Pow(x, 2)
	}

	D := n*sumx2*sumx4 + 2*sumx*sumx2*sumx3 - math.Pow(sumx2, 3) - math.Pow(sumx, 2)*sumx4 - n*math.Pow(sumx3, 2)
	a := (n*sumx2*sumx2y + sumx*sumx3*sumy + sumx*sumx2*sumxy - math.Pow(sumx2, 2)*sumy - math.Pow(sumx, 2)*sumx2y - n*sumx3*sumxy) / D
	b := (n*sumx4*sumxy + sumx*sumx2*sumx2y + sumx2*sumx3*sumy - math.Pow(sumx2, 2)*sumxy - sumx*sumx4*sumy - n*sumx3*sumx2y) / D
	c := (sumx2*sumx4*sumy + sumx2*sumx3*sumxy + sumx*sumx3*sumx2y - math.Pow(sumx2, 2)*sumx2y - sumx*sumx4*sumxy - math.Pow(sumx3, 2)*sumy) / D
	height := math.Exp(c - a*math.Pow(b/(2*a), 2))
	position := -b / (2 * a)
	width := 2.35703 / (math.Sqrt(2) * math.Sqrt(-a))
	params = []FitParameters{
		FitParameters{
			fittype: FitTypeGaussian,
			values:  []float64{a, b, c, D},
			xoffset: xoffset,
			yoffset: yoffset,
		},
		FitParameters{
			fittype: FitTypeParabolic,
			values:  []float64{height, position, width},
			xoffset: xoffset,
			yoffset: yoffset,
		},
	}
	return
}

func (ts *Series) FitLoess(bandwidth float64) (points *Series) {
	xdata := ts.x
	ydata := ts.y
	xval := make([]float64, ts.Len)
	yval := make([]float64, ts.Len)
	for i := range xdata {
		xval[i] = xdata[i]
		yval[i] = ydata[i]
	}
	var distinctX = array_unique(xval)

	if 2/float64(len(distinctX)) > bandwidth {
		bandwidth = math.Min(float64(len(distinctX)), 1)
	}

	res := []float64{}

	var left = 0
	var right int = int(math.Floor(bandwidth*float64(len(xval))) - 1)

	for i := range xval {
		var x = xval[i]

		if i > 0 {
			if right < len(xval)-1 &&
				xval[int(right)+1]-xval[i] < xval[i]-xval[left] {
				left++
				right++
			}
		}

		var edge int
		if xval[i]-xval[left] > xval[right]-xval[i] {
			edge = left
		} else {
			edge = right
		}
		denom := math.Abs(1.0 / (xval[edge] - x))
		var sumWeights float64 = 0
		var sumX float64 = 0
		var sumXSquared float64 = 0
		var sumY float64 = 0
		var sumXY float64 = 0

		var k = left
		for k <= right {
			var xk = xval[k]
			var yk = yval[k]
			var dist float64
			if k < i {
				dist = (x - xk)
			} else {
				dist = (xk - x)
			}
			var w = tricube(dist * denom)
			var xkw = xk * w
			sumWeights += w
			sumX += xkw
			sumXSquared += xk * xkw
			sumY += yk * w
			sumXY += yk * xkw
			k++
		}

		var meanX = sumX / sumWeights

		var meanY = sumY / sumWeights
		var meanXY = sumXY / sumWeights
		var meanXSquared = sumXSquared / sumWeights

		var beta float64
		if meanXSquared == meanX*meanX {
			beta = 0
		} else {
			beta = (meanXY - meanX*meanY) / (meanXSquared - meanX*meanX)
		}
		alpha := meanY - beta*meanX
		res = append(res, beta*x+alpha)
	}
	newts := NewSeriesFrom(xval, res)
	points = newts
	return
}

func (ts *Series) CoefficientOfDetermination(pred *Series) float64 {
	xdata := ts.x
	ydata := ts.y
	var sse, ssyy float64
	for i := range xdata {
		y := ydata[i]
		ssyy += math.Pow(y-pred.y[i], 2)
		sse += math.Pow(y-ts.Mean, 2)
	}
	return 1 - (ssyy / sse)
}

func (ts *Series) StandardError(pred [][]float64) float64 {
	ydata := ts.y
	var SE float64 = 0
	for i := range ydata {
		SE += math.Pow(ydata[i]-pred[i][1], 2)
	}
	SE = math.Sqrt(SE / (float64(ts.Len) - 2))

	return SE
}

func array_unique(values []float64) []float64 {
	o := make(map[float64]float64)
	r := []float64{}
	for i := range values {
		o[values[i]] = values[i]
	}
	for i := range o {
		r = append(r, o[i])
	}
	return r
}

func tricube(x float64) float64 {
	var tmp = 1 - math.Pow(x, 3)
	return math.Pow(tmp, 3)
}

/**
 * Code extracted from https://github.com/Tom-Alexander/regression-js/
 */
func gaussianElimination(a [][]float64, o int) []float64 {
	x := make([]float64, o)
	n := len(a) - 1
	maxrow := 0
	var tmp float64 = 0
	for i := 0; i < n; i++ {
		maxrow = i
		for j := i + 1; j < n; j++ {
			if math.Abs(a[i][j]) > math.Abs(a[i][maxrow]) {
				maxrow = j
			}
		}
		for k := i; k < n+1; k++ {
			tmp = a[k][i]
			a[k][i] = a[k][maxrow]
			a[k][maxrow] = tmp
		}
		for j := i + 1; j < n; j++ {
			for k := n; k >= i; k-- {
				a[k][j] -= a[k][i] * a[i][j] / a[i][i]
			}
		}
	}
	for j := n - 1; j >= 0; j-- {
		tmp = 0
		for k := j + 1; k < n; k++ {
			tmp += a[k][j] * x[k]
		}
		x[j] = (a[n][j] - tmp) / a[j][j]
	}
	return x
}
