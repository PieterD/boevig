package primegen

import "context"

// Primes generates prime numbers, and publishes them on drain.
// Memory will be allocated to create a prime wheel of size wheelSize.
// Be careful, the size of the wheel is exponential!
// A wheelSize of 5 has a memory footprint of 2310 bytes.
// A wheelSize of 7 has a memory footprint of about 0.5 Megabytes.
// A wheelSize of 10 has a memory footprint of about 6.5 Gigabytes.
func Primes(ctx context.Context, wheelSize int, drain chan<- uint64) error {
	// TODO: This is very old code lifted from past me and needs a rewrite before production use.
	// It was written to be fast, using a prime wheel to skip small divisors and a heap to store the next divisor.
	// The context and channel stuff was added by future me (I guess past me again now.)
	q, b, initialPrimes, wheel := buildSieve(wheelSize)
	var D = [][2]uint64{{b * b, b}}

	for _, initialPrime := range initialPrimes {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case drain <- initialPrime:
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if D[0][0] != q {
			// OLD WAY: list = append(list, q)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case drain <- q:
			}
			D = pushHeap(D, q*q, q)
		} else {
			for D[0][0] == q {
				updateHeap(D, wheel, q)
			}
		}
		q += lookupWheel(wheel, q)
	}
}

// Build a wheel sieve to avoid all numbers divisible by the first num primes.
// The size of the wheel is exponential. Anything more than 10 primes
// will create an enormous wheel; try 7.
func buildSieve(num int) (uint64, uint64, []uint64, []byte) {
	var startprimes = []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	n1 := startprimes[num]
	n2 := startprimes[num+1]
	sl := startprimes[:num+1]
	startprimes = startprimes[:num]

	// Multiply our starting primes.
	var product int = 1
	for i := range startprimes {
		product *= int(startprimes[i])
	}

	// We start with an array with one element for every number less than
	// product, initialized to 0.
	var wheel = make([]byte, product+1)
	// Set the element at every prime power for every starting prime to 1.
	for _, p := range startprimes {
		for i := int(p); i <= product; i += int(p) {
			wheel[i] = 1
		}
		/*
			for i:=int(2*p); i<=product; i+=int(p) {
				wheel[i] = 1
			}
			wheel[p] = 1
		*/
	}

	// The increment to the next possible-prime from the current maybe-prime
	// is equal to the number of ones that follow our current maybe-prime,
	// plus one. To get this, traverse the array in reverse, adding up ones
	// until we find a zero.
	// Halve the size of the array by ignoring even numbers.
	var c int = 1
	var wheel2 = make([]byte, product/2)
	for i := product; i > 0; i-- {
		if wheel[i] == 1 {
			c++
		} else {
			if c > 255 {
				panic("OH GOD c > 255")
			}
			wheel2[i/2] = byte(c)
			c = 1
		}
	}

	return n2, n1, sl, wheel2
}

// Look up the increment to reach the next possible-prime from the given
// maybe-prime.
func lookupWheel(wheel []byte, num uint64) uint64 {
	return uint64(wheel[(num%uint64(len(wheel)*2))/2])
}

// Move some horribleness from the inner loop.
func updateHeap(D [][2]uint64, wheel []byte, q uint64) {
	p := D[0][1]
	nv := D[0][0]
	m := nv / p
	// As before, we add p the same number of times as we increment q.
	// That way we don't have to clean up all the prime powers less than
	// q every time we cycle.
	nv += p * lookupWheel(wheel, m)
	D[0][0] = nv
	downHeap(D, 0)
}

// Nicked from the Golang heap library because I needed downHeap.
func upHeap(list [][2]uint64, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || (list[i][0] < list[j][0]) {
			break
		}
		list[i], list[j] = list[j], list[i]
		j = i
	}
}

func downHeap(list [][2]uint64, i int) {
	n := len(list)
	for {
		j1 := 2*i + 1
		if j1 >= n {
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !(list[j1][0] < list[j2][0]) {
			j = j2 // right child
		}
		if list[i][0] < list[j][0] {
			break
		}
		list[i], list[j] = list[j], list[i]
		i = j
	}
}

func pushHeap(list [][2]uint64, a, b uint64) [][2]uint64 {
	list = append(list, [2]uint64{a, b})
	upHeap(list, len(list)-1)
	return list
}
