package main

var Tables = []interface{}{&PublicKey{}, &Service{}, &URL{}}

func Migrate() []interface{} {
	return Tables
}