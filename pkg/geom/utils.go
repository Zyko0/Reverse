package geom

func DistanceSq(x0, y0, x1, y1 float64) float64 {
	return (x1-x0)*(x1-x0) + (y1-y0)*(y1-y0)
}
