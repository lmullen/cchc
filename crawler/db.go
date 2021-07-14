package main

// DBCreateSchema creates the database tables and sets up the database or returns
// an error.
func (app *App) DBCreateSchema() error {

	// TODO: Delete the DROP TABLE statements, which are there for dev purposes
	_, err := app.DB.Exec(`DROP TABLE IF EXISTS collections CASCADE`)
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

	_, err = app.DB.Exec(`DROP TABLE IF EXISTS items CASCADE`)
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

	_, err = app.DB.Exec(`DROP TABLE IF EXISTS items_in_collections CASCADE`)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(`CREATE TABLE IF NOT EXISTS items_in_collections (
		item_id  text REFERENCES items(id),
		collection_id text REFERENCES collections(id),
		PRIMARY KEY (item_id, collection_id)
	);`)
	if err != nil {
		return err
	}

	return nil

}
