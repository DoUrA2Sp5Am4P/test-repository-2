package main

import (
	"context"
	"database/sql"
	"os"
	"time"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadAll() {
	b, err := exists("store.db")
	if b == false {
		InfoLogger.Println("Не найдено предыдущих сессий")
		return
	}
	os.Remove("store.db-journal")
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
	err = selectExpressions(ctx, db)
	if err != nil {
		ErrorLogger.Println("error while reading expressins")
		ErrorLogger.Println(err)
	}
	err = selectTasks(ctx, db)
	if err != nil {
		ErrorLogger.Println("error while reading tasks")
		ErrorLogger.Println(err)
	}
	err = selectResults(ctx, db)
	if err != nil {
		ErrorLogger.Println("error while reading results")
		ErrorLogger.Println(err)
	}
	err = selectOperations(ctx, db)
	if err != nil {
		ErrorLogger.Println("error while reading operations")
		ErrorLogger.Println(err)
	}
	err = selectSleep(ctx, db)
	if err != nil {
		ErrorLogger.Println("error while reading sleep")
		ErrorLogger.Println(err)
	}

}

func selectExpressions(ctx context.Context, db *sql.DB) error {
	var q = `
	SELECT username,id,expression,status,createdat,completedat FROM expressions
	`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		e := Expression{}
		c := ""
		c2 := ""
		err := rows.Scan(&e.Username, &e.ID, &e.Expression, &e.Status, &c, &c2)
		e.CreatedAt, _ = time.Parse(time.RFC3339, c)
		e.CompletedAt, _ = time.Parse(time.RFC3339, c2)
		if err != nil {
			return err
		}
		expressions = append(expressions, e)
	}
	return nil
}

func selectTasks(ctx context.Context, db *sql.DB) error {
	var q = `
	SELECT username,id,expression,status,createdat,completedat FROM tasks
	`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		e := Expression{}
		c := ""
		c2 := ""
		err := rows.Scan(&e.Username, &e.ID, &e.Expression, &e.Status, &c, &c2)
		e.CreatedAt, _ = time.Parse(time.RFC3339, c)
		e.CompletedAt, _ = time.Parse(time.RFC3339, c2)
		if err != nil {
			return err
		}
		tasks = append(tasks, e)
	}
	return nil
}

func selectResults(ctx context.Context, db *sql.DB) error {
	var q = `
	SELECT username,id,result FROM expressionswithresult
	`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		e := Expression_with_result{}
		err := rows.Scan(&e.Username, &e.ID, &e.Result)
		if err != nil {
			return err
		}
		results = append(results, e)
	}
	return nil
}

func selectOperations(ctx context.Context, db *sql.DB) error {
	var q = `
	SELECT username,name,executiontime FROM operations
	`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		e := Operation{}
		err := rows.Scan(&e.Username, &e.Name, &e.ExecutionTime)
		if err != nil {
			return err
		}
		operations = append(operations, e)
	}
	return nil
}

func selectSleep(ctx context.Context, db *sql.DB) error {
	var q = `
	SELECT username,sl FROM sleep
	`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		e := sleep{}
		err := rows.Scan(&e.Username, &e.time)
		if err != nil {
			return err
		}
		sl = append(sl, e)
	}
	return nil
}
