package utils

import "slices"

// gcd 计算两个数的最大公约数
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// findCoprimeNumbers 找到 m个 与 n 互质的数
func findCoprimeNumbers(n int64, m int) []int64 {
	coprimeList := make([]int64, 0)
	for i := int64(1); i < n; i++ {
		if gcd(i, n) == 1 {
			coprimeList = append(coprimeList, i)
			if m == len(coprimeList) {
				break
			}
		}
	}
	return coprimeList
}

// isCoprimeWithAll 检查一个数是否与给定的多个数都互质
func isCoprimeWithAll(x int64, numbers ...int64) bool {
	for _, num := range numbers {
		if gcd(x, num) != 1 {
			return false
		}
	}
	return true
}

// randSwap 置换多项式
func randSwap(a, b, x, n int64) int64 {
	return (a*x + b) % n
}

// randSwapGroup 置换多项式组
func randSwapGroup(a, b []int64, x, n int64) int64 {
	for i := 0; i < len(a)-1; i++ {
		if i%2 == 0 {
			if x%2 == 0 {
				x = randSwap(a[i], b[i], x/2, (n+1)/2) * 2
			} else {
				x = randSwap(a[i+1], b[i+1], x/2, n/2)*2 + 1
			}
		} else {
			if x >= n/2 {
				x = randSwap(a[i], b[i], x-n/2, n-n/2) + n/2
			} else {
				x = randSwap(a[i+1], b[i+1], x, n/2)
			}
		}
	}
	return x
}

// genFactorList 生成因数列表
func genFactorList(n int64) (factorList, constantList []int64) {
	coprimeList := findCoprimeNumbers(n/2, 10000)
	factorList = make([]int64, 0)
	if n%2 == 0 {
		factorList = coprimeList
	} else {
		for _, v := range coprimeList {
			if isCoprimeWithAll(v, n/2+1, n) {
				factorList = append(factorList, v)
			}
		}
	}
	constantList = slices.Clone(factorList)
	slices.Reverse(constantList)
	return
}

// GetPosition 获取位置
// n: 最大长度
// position: 原位置
// return: 置换后的位置
func GetPosition(n, position int64) int64 {
	//获取因数列表
	factorList, constantList := genFactorList(n)
	return randSwapGroup(factorList, constantList, position, n)
}
