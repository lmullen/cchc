package main

import "fmt"

// DBInit creates the database schema
func (app *App) DBInit() error {

	_, err := app.DB.Exec(`DROP TABLE IF EXISTS items`)
	if err != nil {
		return err
	}
	temp, err := app.DB.Exec(`CREATE TABLE items (
		id    text PRIMARY KEY,
		lccn  text,
		date  integer,
		title text,
		api   jsonb
	)`)
	fmt.Println(temp)
	if err != nil {
		return err
	}
	return nil

}
