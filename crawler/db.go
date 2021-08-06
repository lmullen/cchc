package main

// DBCreateSchema creates the database tables and sets up the database or returns
// an error.
func (app *App) DBCreateSchema() error {

	// _, err := app.DB.Exec(`DROP TABLE IF EXISTS collections CASCADE`)
	// if err != nil {
	// 	return err
	// }

	_, err := app.DB.Exec(`CREATE TABLE IF NOT EXISTS collections (
		id            text PRIMARY KEY,
		title         text,
		description   text,
		count         integer,
		url           text, 
		items_url     text,
		subjects      text[],
		subjects2     text[],
		topics        text[],
		api           jsonb
	);`)
	if err != nil {
		return err
	}

	// _, err = app.DB.Exec(`DROP TABLE IF EXISTS items CASCADE`)
	// if err != nil {
	// 	return err
	// }

	_, err = app.DB.Exec(`CREATE TABLE IF NOT EXISTS items (
		id                 text PRIMARY KEY,
		url                text,
		title              text,
		year               int,
		date               text,
		subjects           text[],
		fulltext           text,
		fulltext_service   text,
		fulltext_file      text,
		timestamp          bigint,
		api                jsonb
	);`)
	if err != nil {
		return err
	}

	// _, err = app.DB.Exec(`DROP TABLE IF EXISTS items_in_collections CASCADE`)
	// if err != nil {
	// 	return err
	// }

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
