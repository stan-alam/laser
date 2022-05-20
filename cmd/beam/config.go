package main

import (
	"flag"
	"os"
)

type config struct {
	Address          string
	ConnectionString string
}

func configure() *config {
	c := &config{
		Address:          os.Getenv("ADDRESS"),
		ConnectionString: os.Getenv("CONNECTION_STRING"),
	}

	if c.Address == "" {
		flag.StringVar(&c.Address, "address", ":3000", "Service address[:port]")
	}
	if c.ConnectionString == "" {
		flag.StringVar(&c.ConnectionString, "connString", "postgresql://postgres@localhost:5432/postgres", "Postgresql connection string")
	}
	flag.Parse()

	return c
}
