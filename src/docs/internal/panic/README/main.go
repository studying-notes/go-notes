package main

func a() {
	for i := 0; i < 3; i++ {
		defer fa()
	}
	b()
}

func fa() {
	println("fa")
}

func b() {
	defer fb1()
	if true {
		defer fb2()
	}
	c()
}

func fb1() {
	println("fb1")
}

func fb2() {
	println("fb2")
}

func c() {
	for i := 0; i < 3; i++ {
		defer fc()
	}

	panic("c")
}

func fc() {
	println("fc")
}
