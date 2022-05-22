package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

var (
	ctx context.Context
	db  *sql.DB
)

func Query() (string, error) {
	var name string
	age := 27
	rows, err := db.QueryContext(ctx, "SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	err = rows.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return name, nil
}

func main() {
	name, err := Query()
	if err != nil {
		println(err.Error())
		return

	}
	fmt.Printf("res:%s", name)
}
