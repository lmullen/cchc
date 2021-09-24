package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/jackc/pgx/v4"
// 	"github.com/lmullen/cchc/qftext/models"
// )

// func GetUnqueued(ctx context.Context) error {

// 	tx, err := app.DB.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	_, err = tx.Exec(ctx, `DECLARE unqcurs NO SCROLL CURSOR FOR SELECT id FROM jobs.fulltext_unqueued;`)
// 	if err != nil {
// 		return err
// 	}

// 	var item models.JobFulltextPredict
// 	for {
// 		err = tx.QueryRow(ctx, `FETCH NEXT FROM unqcurs;`).Scan(&item.ItemID)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				break
// 			}
// 			return err
// 		}
// 		fmt.Println(item)
// 	}

// 	_, err = tx.Exec(ctx, "CLOSE unqcurs;")
// 	if err != nil {
// 		return err
// 	}

// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
