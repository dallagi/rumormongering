package main

type avg struct {
	c float64
	n int
}

func (a *avg) Add(v float64) {
	a.c += v
	a.n++
}

func (a *avg) Get() float64 {
	return a.c / float64(a.n)
}
