package main

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestSaveAndRead(t *testing.T) {
	user_ := "save_test_save_read"
	currentUser = user_
	sl = append(sl, sleep{time: 1, Username: user_})
	expressions = append(expressions, Expression{Username: user_, ID: "123", Expression: "2+2", Status: "Не посчитано", CreatedAt: time.Now(), CompletedAt: time.Now()})
	tasks = append(tasks, Expression{Username: user_, ID: "123", Expression: "2+2", Status: "Не посчитано", CreatedAt: time.Now(), CompletedAt: time.Now()})
	results = append(results, Expression_with_result{Username: user_, Result: "4.00", ID: "123"})
	operations = append(operations, Operation{Name: "Сложение", Username: user_, ExecutionTime: 1})
	operations = append(operations, Operation{Name: "Вычитание", Username: user_, ExecutionTime: 2})
	operations = append(operations, Operation{Name: "Умножение", Username: user_, ExecutionTime: 3})
	operations = append(operations, Operation{Name: "Деление", Username: user_, ExecutionTime: 4})
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		t.Error(err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		t.Error(err)
	}
	if err = createTables(ctx, db); err != nil {
		t.Error(err)
	}
	saveAll(ctx, db)
	db.Close()
	sl = []sleep{}
	expressions = []Expression{}
	tasks = []Expression{}
	results = []Expression_with_result{}
	operations = []Operation{}
	ReadAll()
	var a = 0
	for _, expres := range expressions {
		if expres.ID == "123" {
			a = 1
		}
	}
	if a != 1 {
		t.Error("Expressions not saved")
	}
	b := 0
	for _, task := range tasks {
		if task.ID == "123" {
			b = 1
		}
	}
	if b != 1 {
		t.Error("Tasks not saved")
	}
	c := 0
	for _, res := range results {
		if res.ID == "123" {
			c = 1
		}
	}
	if c != 1 {
		t.Error("Results not saved")
	}
	d := 0
	for _, res := range operations {
		if res.Name == "Сложение" {
			d = 1
		}
	}
	if d != 1 {
		t.Error("Addiction not saved")
	}
	e := 0
	for _, res := range operations {
		if res.Name == "Вычитание" {
			e = 1
		}
	}
	if e != 1 {
		t.Error("Subtraction not saved")
	}
	j := 0
	for _, res := range operations {
		if res.Name == "Умножение" {
			j = 1
		}
	}
	if j != 1 {
		t.Error("Multiplication not saved")
	}
	z := 0
	for _, res := range operations {
		if res.Name == "Деление" {
			z = 1
		}
	}
	if z != 1 {
		t.Error("Division not saved")
	}
}
