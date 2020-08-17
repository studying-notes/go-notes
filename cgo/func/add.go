package main

/*
#include <stdio.h>

static int add(int a, int b) {
    printf("%d", a+b);
}
*/
import "C"

func main() {
	C.add(1, 1)
}
