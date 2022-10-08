package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func (s *Student) CreateSQL() string {
	return fmt.Sprintf("INSERT INTO student VALUES (%s, %d)", s.Name, s.Age)
}

type SQL interface {
	CreateSQL() string
}

func createQuery(q interface{}) string {
	if reflect.TypeOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		v := reflect.ValueOf(q)
		query := fmt.Sprintf("INSERT INTO %s VALUES (", t)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query += fmt.Sprintf("%d", v.Field(i).Int())
				} else {
					query += fmt.Sprintf(", %d", v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query += fmt.Sprintf("'%s'", v.Field(i).String())
				} else {
					query += fmt.Sprintf(",'%s'", v.Field(i).String())
				}
			}
		}
		query += ")"
		return query
	}
	return ""
}

type Trade struct {
	id    int
	Price int
}

func main() {
	s := Student{Name: "John", Age: 20}
	fmt.Println(s.CreateSQL())

	t := Trade{id: 1, Price: 100}
	fmt.Println(createQuery(t))
}
