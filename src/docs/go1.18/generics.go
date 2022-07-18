/*
 * @Date: 2022.05.02 14:42
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.05.02 14:42
 */

package main

import "golang.org/x/exp/constraints"

func GMin[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}
