package main

// Database
type Database interface {
	Transaction() (interface{}, error)
	Session() (interface{}, error)
}

type Storage interface {
}
