package main

import "testing"

func BenchmarkNewPeople(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPeople()
	}
}

func BenchmarkNewPeopleByReflect(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPeopleByReflect()
	}
}
