package frac

import "fmt"

type Fraction struct {
	Overflow   bool
	Normalized bool
	Num        int64
	Den        int64
}

func New(num, den int64) Fraction {
	return Fraction{
		Num: num,
		Den: den,
	}
}

func (f Fraction) Normalize() Fraction {
	if f.Normalized || f.Overflow {
		return f
	}
	gcd := GCD(f.Num, f.Den)
	for gcd > 1 {
		f.Num /= gcd
		f.Den /= gcd
	}
	f.Normalized = true
	return f
}

func (f Fraction) Multiply(f2 Fraction) Fraction {
	if f.Overflow {
		return f
	}
	num, numOflow := Mul(f.Num, f2.Num)
	den, denOflow := Mul(f.Den, f2.Den)
	if numOflow || denOflow {
		f = f.Normalize()
		f2 = f2.Normalize()
		num, numOflow = Mul(f.Num, f2.Num)
		den, denOflow = Mul(f.Den, f2.Den)
	}
	if numOflow || denOflow {
		return Fraction{Overflow: true}
	}
	f.Num = num
	f.Den = den
	f.Normalized = false
	return f
}

func (f Fraction) ScalarMultiply(s int64) Fraction {
	if f.Overflow {
		return f
	}
	m, oflow := Mul(f.Num, s)
	if oflow {
		f = f.Normalize()
		m, oflow = Mul(f.Num, s)
	}
	if oflow {
		return Fraction{Overflow: true}
	}
	f.Num = m
	f.Normalized = false
	return f
}

func (f Fraction) Less(f2 Fraction) bool {
	// Given a/b < c/d
	// a/b - c/d < 0
	// (ad - bc)/(bd) < 0
	// bd does not affect sign, so:
	// ad - bc < 0
	ad, adOflow := Mul(f.Num, f2.Den)
	bc, bcOflow := Mul(f2.Num, f.Den)
	if adOflow || bcOflow {
		ad, adOflow = Mul(f.Num, f2.Den)
		bc, bcOflow = Mul(f2.Num, f.Den)
	}
	if adOflow || bcOflow {
		panic(fmt.Errorf("fraction overflow"))
	}
	return ad-bc < 0
}

func Mul(l, r int64) (i int64, overflow bool) {
	i = l * r
	overflow = i/l == r
	return i, overflow
}

func GCD(l, r int64) int64 {
	for l != r {
		if l > r {
			l -= r
		} else {
			r -= l
		}
	}
	return l
}
