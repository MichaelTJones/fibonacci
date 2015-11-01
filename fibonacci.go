package fibonacci

import (
	"math/big"
)

// static cache of first 92 Fibonacci numbers--these are precisely the
// values representable in 63 bits.
var fibonacciTable = []int64{
	0, // value for index 0
	1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597,
	2584, 4181, 6765, 10946, 17711, 28657, 46368, 75025, 121393, 196418,
	317811, 514229, 832040, 1346269, 2178309, 3524578, 5702887, 9227465,
	14930352, 24157817, 39088169, 63245986, 102334155, 165580141,
	267914296, 433494437, 701408733, 1134903170, 1836311903, 2971215073,
	4807526976, 7778742049, 12586269025, 20365011074, 32951280099,
	53316291173, 86267571272, 139583862445, 225851433717, 365435296162,
	591286729879, 956722026041, 1548008755920, 2504730781961, 4052739537881,
	6557470319842, 10610209857723, 17167680177565, 27777890035288,
	44945570212853, 72723460248141, 117669030460994, 190392490709135,
	308061521170129, 498454011879264, 806515533049393, 1304969544928657,
	2111485077978050, 3416454622906707, 5527939700884757, 8944394323791464,
	14472334024676221, 23416728348467685, 37889062373143906,
	61305790721611591, 99194853094755497, 160500643816367088,
	259695496911122585, 420196140727489673, 679891637638612258,
	1100087778366101931, 1779979416004714189, 2880067194370816120,
	4660046610375530309, 7540113804746346429,
}

// compute the Nth Fibonacci number using big integer arithmentic.  Efficient algorithms
// are used so the 10 millionth value requires about a second to create the result, which
// has 2,089,877 digits when formatted in decimal. Time measured and algorithm breakpoint
// determined on 2013 MacBook Pro test system (2.7 GHz Intel Core i7, MacBookPro10,1)
func Fibonacci(n int) (f *big.Int) {
	switch {
	case n < 1:
		f = big.NewInt(0)

	// static table for small cases (optional, but always faster and less storage debris)
	case n < len(fibonacciTable):
		f = big.NewInt(fibonacciTable[n])

	// big integer evaluation using algorithims in their most efficient ranges
	case n <= 100: // Direct series evaluation is fast for small values
		f = fibSeries(n)
	case n <= 5504: // Blenkinsop algorithm is faster for values in 100..5504 on test system
		f = fibBlenkinsop(n)
	default: // Takahashi algorithm is faster for values > 5504 on test system
		f = fibTakahashi(n)
	}
	return
}

func log2(n int) (bits int) {
	for n>>uint(bits+1) != 0 {
		bits++
	}
	return
}

func fibSeries(n int) *big.Int {
	if n < 1 {
		return big.NewInt(0)
	}

	a := big.NewInt(0)
	b := big.NewInt(1)

	for i := 0; i < n; i++ {
		a, b = b, a.Add(a, b)
	}

	return a
}

func fibBlenkinsop(n int) *big.Int {
	if n < 1 {
		return big.NewInt(0)
	}

	h := uint(log2(n))

	f1 := big.NewInt(0)
	f2 := big.NewInt(1)
	f3 := new(big.Int)

	for ; h > 0; h-- {
		f3.Add(f1, f2)
		if (n>>(h-1))&1 == 1 {
			f1.Add(f1, f3).Mul(f1, f2)
			f2.Mul(f2, f2).Add(f2, f3.Mul(f3, f3))
		} else {
			f3.Add(f1, f3).Mul(f3, f2)
			f1.Mul(f1, f1).Add(f1, f2.Mul(f2, f2))
			f2, f3 = f3, f2
		}
	}
	return f2
}

var c2 = big.NewInt(2)
var c3 = big.NewInt(3)
var c5 = big.NewInt(5)

func fibTakahashi(n int) *big.Int {
	if n <= 0 {
		return big.NewInt(0)
	}
	if n <= 2 {
		return big.NewInt(1)
	}

	f := big.NewInt(1)
	l := big.NewInt(1)
	sign := big.NewInt(-1)

	bits := log2(n)
	mask := 1 << uint(bits-1)

	t1 := big.NewInt(0)
	t2 := big.NewInt(0)

	for i := 1; i < bits; i++ {
		t1.Mul(f, f)          // t1 := f * f
		f.Add(f, l).Rsh(f, 1) // f = (f + l) >> 1
		f.Mul(f, f).Lsh(f, 1) // f = (f*f)<<1 - 3*t1 - 2*sign
		f.Sub(f, t2.Mul(t1, c3))
		f.Sub(f, t2.Mul(sign, c2))
		l.Mul(t1, c5) // l = 5*t1 + 2*sign
		l.Add(l, t2.Mul(sign, c2))

		sign.SetInt64(1) // sign = 1

		if n&mask != 0 {
			t1.Set(f)             //t1 = f
			f.Add(f, l).Rsh(f, 1) //f = (f + l) >> 1
			t1.Lsh(t1, 1)         //l = f + 2*t1
			l.Add(f, t1)
			sign.SetInt64(-1) //sign = -1
		}
		mask >>= 1
	}

	if n&mask == 0 {
		f.Mul(f, l) //f = f * l
	} else {
		f.Add(f, l).Rsh(f, 1)    //f = (f + l) >> 1
		f.Mul(f, l).Sub(f, sign) //f = f*l - sign
	}

	return f
}

// fib(k) returns the kth fibonacci number
func fibDouble(k int) *big.Int {
	// http://www.nayuki.io/page/fast-fibonacci-algorithms
	var a, b, c = big.NewInt(0), big.NewInt(1), new(big.Int)
	var bit uint64
	for bit = 1 << 63; bit > 0; bit >>= 1 {
		// a, b = a*b + a*b - a*a, b*b + a*a
		c.Mul(a, b).Add(c, c).Sub(c, a.Mul(a, a))
		b.Add(b.Mul(b, b), a)
		a.Set(c)
		if uint64(k)&bit != 0 {
			c.Add(a, b)
			a.Set(b)
			b.Set(c)
		}
	}
	return a
}
