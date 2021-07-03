package main

// DBInit creates the database schema
func (app *App) DBInit() error {

	_, err := app.DB.Exec(`DROP TABLE IF EXISTS collections`)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(`CREATE TABLE IF NOT EXISTS collections (
		id          text PRIMARY KEY,
		title       text,
		count       integer,
		url         text, 
		items_url   text,
		subjects    text[],
		api         jsonb
	);`)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(`DROP TABLE IF EXISTS items`)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(`CREATE TABLE IF NOT EXISTS items (
		id       text PRIMARY KEY,
		lccn     text,
		url      text,
		date     integer,
		subjects text[],
		title    text,
		api      jsonb
	);`)
	if err != nil {
		return err
	}

	return nil

}
