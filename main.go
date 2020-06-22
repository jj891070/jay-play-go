package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
)

func main() {

	log.Println(fractionToDecimal(6, 3))
	log.Println(TurnToFraction(0))
	x, y := TurnToFraction(0)
	log.Println(Gcd(int(x), int(y)))

}

// Gcd 輾轉相除法，取最大公因數
func Gcd(x, y int) int {
	tmp := x % y
	if tmp > 0 {
		return Gcd(y, tmp)
	} else {
		return y
	}
}

// TurnToFraction 將小數轉成分子(numerator)／分母(numerator)
func TurnToFraction(target float64) (numerator float64, denominator float64) {
	numerator = target
	denominator = 1
	for {
		targetInteger := int(numerator)
		ans := numerator - float64(targetInteger)
		if ans == 0 {
			break
		}
		numerator *= 10
		denominator *= 10
	}
	return
}

// fractionToDecimal 擷取循環小數
func fractionToDecimal(numerator int, denominator int) string {
	if numerator == 0 {
		return "0"
	}
	if denominator == 0 {
		return "NaN"
	}

	// 判斷分子或分母小於0，加個負號
	var buffer bytes.Buffer
	if (numerator < 0 && denominator > 0) || (numerator > 0 && denominator < 0) {
		buffer.WriteString("-")
	}

	num := int(math.Abs(float64(numerator)))
	denom := int(math.Abs(float64(denominator)))

	buffer.WriteString(strconv.Itoa(num / denom))

	num %= denom

	if num == 0 {
		return buffer.String()
	}
	buffer.WriteString(".")
	log.Println(buffer.String())

	m := make(map[int]int, 10)
	repeatPos := -1
	for {
		num *= 10
		pos, ok := m[num]
		if !ok {
			m[num] = buffer.Len()
		} else {
			repeatPos = pos
			break
		}
		buffer.WriteString(strconv.Itoa(num / denom))
		//fmt.Println(buffer, len(buffer), num)
		num %= denom
		if num == 0 {
			break
		}
	}

	if repeatPos == -1 {
		return buffer.String()
	}
	res := buffer.String()
	return fmt.Sprintf("%s(%s)", res[0:repeatPos], res[repeatPos:])
}
