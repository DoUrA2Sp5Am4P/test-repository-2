package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func createTables(ctx context.Context, db *sql.DB) error {
	const (
		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		username TEXT,
		id INTEGER UNIQUE, 
		expression TEXT,
		status TEXT,
		createdat TEXT,
		completedat TEXT
	);`
		tasksTable = `
	CREATE TABLE IF NOT EXISTS tasks(
		username TEXT,
		id INTEGER UNIQUE, 
		expression TEXT,
		status TEXT,
		createdat TEXT,
		completedat TEXT
	);`
		expressions_with_resultTable = `
	CREATE TABLE IF NOT EXISTS expressionswithresult(
		username TEXT,
		id TEXT UNIQUE,
		result TEXT
	);`
		operationsTable = `
	CREATE TABLE IF NOT EXISTS operations(
		username TEXT,
		name TEXT,
		executiontime INTEGER
	);`
		sleepTable = `
	CREATE TABLE IF NOT EXISTS sleep(
		username TEXT UNIQUE,
		sl INTEGER
	);`
	)
	if _, err := db.Exec(expressionsTable); err != nil {
		return err
	}
	if _, err := db.Exec(tasksTable); err != nil {
		return err
	}
	if _, err := db.Exec(expressions_with_resultTable); err != nil {
		return err
	}
	if _, err := db.Exec(operationsTable); err != nil {
		return err
	}
	if _, err := db.Exec(sleepTable); err != nil {
		return err
	}
	return nil
}

func insertExpressions(ctx context.Context, db *sql.DB) error {
	var q = `
	INSERT INTO expressions (username,id,expression,status,createdat,completedat) values ($1, $2,$3,$4,$5,$6)
	`
	var q1 = `
	REPLACE INTO expressions (username,id,expression,status,createdat,completedat) values ($1, $2,$3,$4,$5,$6)
	`
	for _, expression := range expressions {
		_, err := db.Exec(q, expression.Username, expression.ID, expression.Expression, expression.Status, expression.CreatedAt.Format(time.RFC3339), expression.CompletedAt.Format(time.RFC3339))
		if err != nil {
			_, err := db.Exec(q1, expression.Username, expression.ID, expression.Expression, expression.Status, expression.CreatedAt.Format(time.RFC3339), expression.CompletedAt.Format(time.RFC3339))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func insertTasks(ctx context.Context, db *sql.DB) error {
	var q = `
	INSERT INTO tasks (username,id,expression,status,createdat,completedat) values ($1, $2,$3,$4,$5,$6)
	`
	var q1 = `
	REPLACE INTO tasks (username,id,expression,status,createdat,completedat) values ($1, $2,$3,$4,$5,$6)
	`
	for _, expression := range tasks {
		_, err := db.Exec(q, expression.Username, expression.ID, expression.Expression, expression.Status, expression.CreatedAt.Format(time.RFC3339), expression.CompletedAt.Format(time.RFC3339))
		if err != nil {
			_, err := db.Exec(q1, expression.Username, expression.ID, expression.Expression, expression.Status, expression.CreatedAt.Format(time.RFC3339), expression.CompletedAt.Format(time.RFC3339))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func InsertExpression_with_result(ctx context.Context, db *sql.DB) error {
	var q = `
	INSERT INTO expressionswithresult (username,id,result) values ($1, $2,$3)
	`
	var q1 = `
	REPLACE INTO expressionswithresult (username,id,result) values ($1, $2,$3)
	`
	for _, expression := range results {
		_, err := db.Exec(q, expression.Username, expression.ID, expression.Result)
		if err != nil {
			_, err := db.Exec(q1, expression.Username, expression.ID, expression.Result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func InsertOperations(ctx context.Context, db *sql.DB) error {
	var q = `
	INSERT INTO operations (username,name,executiontime) values ($1, $2,$3)
	`
	var q1 = `
	UPDATE operations SET executiontime = $1 WHERE username = $2 AND name = $3;
	`
	var q3 = `
	SELECT username,name,executiontime FROM operations
	`
	for _, operation := range operations {
		rows, err := db.Query(q3)
		if err != nil {
			return err
		}
		defer rows.Close()
		v := 0
		for rows.Next() {
			e := Operation{}
			err := rows.Scan(&e.Username, &e.Name, &e.ExecutionTime)
			if err != nil {
				return err
			}
			if e.Username == currentUser && e.Name == operation.Name {
				v = 1
			}
		}
		if v == 1 {
			_, err = db.Exec(q1, operation.ExecutionTime, currentUser, operation.Name)
		} else {
			_, err = db.Exec(q, currentUser, operation.Name, operation.ExecutionTime)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func insertSleep(ctx context.Context, db *sql.DB) error {
	var q = `
	INSERT INTO sleep (username,sl) values ($1,$2)
	`
	var q1 = `
	REPLACE INTO sleep (username,sl) values ($1,$2)
	`
	for _, s := range sl {
		_, err := db.Exec(q, currentUser, s.time)
		if err != nil {
			_, err := db.Exec(q1, currentUser, s.time)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func saveAll(ctx context.Context, db *sql.DB) {
	err := insertExpressions(ctx, db)
	if err != nil {
		ErrorLogger.Println("While saving expressions occured error:")
		ErrorLogger.Println(err)
	}
	err = insertTasks(ctx, db)
	if err != nil {
		ErrorLogger.Println("While saving tasks occured error:")
		ErrorLogger.Println(err)
	}
	err = InsertExpression_with_result(ctx, db)
	if err != nil {
		ErrorLogger.Println("While saving expression_with_result occured error:")
		ErrorLogger.Println(err)
	}
	err = InsertOperations(ctx, db)
	if err != nil {
		ErrorLogger.Println("While saving operations occured error:")
		ErrorLogger.Println(err)
	}
	err = insertSleep(ctx, db)
	if err != nil {
		ErrorLogger.Println("While saving sleep occured error:")
		ErrorLogger.Println(err)
	}
}

func saveAll_r() {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	if err = createTables(ctx, db); err != nil {
		panic(err)
	}
	for {
		if Exit == 0 {
			saveAll(ctx, db)
		} else {
			db.Close()
		}
		InfoLogger.Println("saved")
		if (len(sl) > 0) && !(currentUser == "") {
			for i := 0; i < len(sl); i++ {
				if sl[i].Username == currentUser {
					sl_1 := time.Duration(sl[i].time)
					time.Sleep(sl_1 * time.Second)
				}
			}
		} else {
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}
