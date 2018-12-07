package main

import "../leveldb"

func main() {
	options := leveldb.NewOptions()
	
	var db leveldb.DB 

	leveldb.Open(*options, "/tmp/testdb", &db)
}