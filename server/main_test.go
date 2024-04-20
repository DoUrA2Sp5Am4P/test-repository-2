package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testSave() {
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
	saveAll(ctx, db)
}
func Test401Handler(t *testing.T) {
	f, _ := exists("store.db")
	if f == false {
		testSave()
		t.Log("Database not exist!")
		t.Log("Please wait while database is creating!")
	}
	currentUser = username_test
	ReadAll()
	req := httptest.NewRequest(http.MethodGet, "/expressions", nil)
	w := httptest.NewRecorder()
	AuthorizedUsers(http.HandlerFunc(handleExpressions)).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error("error in test401handler", err)
	}
	if w.Result().StatusCode != 303 {
		t.Error("test 401handler result error")
	}
}

var username_test = "admin"
var password_test = "admin"

func TestRegister(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"username": username_test,
		"password": password_test,
	}
	body, _ := json.Marshal(mcPostBody)
	req := httptest.NewRequest(http.MethodGet, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	Signin(w, req)
	res := w.Result()
	defer res.Body.Close()
	status_test := res.StatusCode
	if status_test != http.StatusSeeOther {
		t.Error("Error in TestRegister")
	}
	currentUser = username_test
}

func TestLogin(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"username": username_test,
		"password": password_test,
	}
	body, _ := json.Marshal(mcPostBody)
	req := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	LoginHandle(w, req)
	res := w.Result()
	defer res.Body.Close()
	status_test := res.StatusCode
	if status_test != http.StatusSeeOther {
		t.Error("Error in TestLogin")
	}
}

func TestAddition(t *testing.T) {
	req, err := http.NewRequest("POST", "/expressions?expression=2+2", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handler := AuthorizedUsers(http.HandlerFunc(handleExpressions))
	handler.ServeHTTP(rr, req)
	data, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Log(err)
	}
	if len(string(data)) != 19 {
		t.Error("Error in TestAddition")
	}
}

func TestSubstraction(t *testing.T) {
	req, err := http.NewRequest("POST", "/expressions?expression=4-2", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handler := AuthorizedUsers(http.HandlerFunc(handleExpressions))
	handler.ServeHTTP(rr, req)
	data, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Log(err)
	}
	if len(string(data)) != 19 {
		t.Error("Error in TestAddition")
	}
}

func TestMultiplication(t *testing.T) {
	req, err := http.NewRequest("POST", "/expressions?expression=4*2", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handler := AuthorizedUsers(http.HandlerFunc(handleExpressions))
	handler.ServeHTTP(rr, req)
	data, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Log(err)
	}
	if len(string(data)) != 19 {
		t.Error("Error in TestAddition")
	}
}

func TestDividion(t *testing.T) {
	req, err := http.NewRequest("POST", "/expressions?expression=4/2", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handler := AuthorizedUsers(http.HandlerFunc(handleExpressions))
	handler.ServeHTTP(rr, req)
	data, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Log(err)
	}
	if len(string(data)) != 19 {
		t.Error("Error in TestAddition")
	}
}

func TestGetExpressions(t *testing.T) {
	req, err := http.NewRequest("GET", "/expressions", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handleExpressions(rr, req)
	res := rr.Result()
	defer res.Body.Close()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error(err)
	}
	data, _ := io.ReadAll(res.Body)
	expr_test := []Expression{}
	json.Unmarshal(data, &expr_test)
	if len(expr_test) > 0 {
		req2, err := http.NewRequest("GET", "/expressions/"+expr_test[0].ID, nil)
		if err != nil {
			t.Error("Error getExpression by id")
		}
		AddCookie(req2)
		rr2 := httptest.NewRecorder()
		handleExpressionByID(rr2, req2)
		res2 := rr2.Result()
		defer res2.Body.Close()
		data2, _ := io.ReadAll(res2.Body)
		expr_test2 := Expression{}
		err = json.Unmarshal(data2, &expr_test2)
		if err != nil {
			t.Error(err)
		}
		if expr_test2.ID != expr_test[0].ID {
			t.Error("Error in TestGetExpressions")
		}
	}
}

func AddCookie(req *http.Request) {
	expirationTime := time.Now().Add(800 * time.Minute)
	claims := &Claims{
		Username: username_test,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	cookie := &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}
	req.AddCookie(cookie)
}
func TestOperations(t *testing.T) {
	currentUser = username_test
	operations = append(operations, Operation{Name: "Сложение", Username: currentUser, ExecutionTime: 1})
	operations = append(operations, Operation{Name: "Вычитание", Username: currentUser, ExecutionTime: 2})
	operations = append(operations, Operation{Name: "Умножение", Username: currentUser, ExecutionTime: 3})
	operations = append(operations, Operation{Name: "Деление", Username: currentUser, ExecutionTime: 4})
	req, err := http.NewRequest("POST", "/operations/addition?time=2", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("POST", "/operations/subtraction?time=3", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("POST", "/operations/multiplication?time=4", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("POST", "/operations/division?time=5", nil)
	if err != nil {
		t.Error(err)
	}
	req5, err := http.NewRequest("GET", "/operations/operations", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	AddCookie(req2)
	AddCookie(req3)
	AddCookie(req4)
	AddCookie(req5)
	rr := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()
	rr3 := httptest.NewRecorder()
	rr4 := httptest.NewRecorder()
	rr5 := httptest.NewRecorder()
	handleAddition(rr, req)
	handleSubtraction(rr2, req2)
	handleMultiplication(rr3, req3)
	handleDivision(rr4, req4)
	handleOperations(rr5, req5)
	defer rr.Result().Body.Close()
	defer rr2.Result().Body.Close()
	defer rr3.Result().Body.Close()
	defer rr4.Result().Body.Close()
	defer rr5.Result().Body.Close()
	data, _ := io.ReadAll(rr5.Result().Body)
	expr_test2 := []Operation{}
	err = json.Unmarshal(data, &expr_test2)
	if err != nil {
		t.Error(err)
	}
	if expr_test2[0].ExecutionTime != 2 || expr_test2[1].ExecutionTime != 3 || expr_test2[2].ExecutionTime != 4 || expr_test2[3].ExecutionTime != 5 {
		t.Error("Error in TestOperations")
	}
}

func TestGetResults(t *testing.T) {
	req, err := http.NewRequest("GET", "/results", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handleResults(rr, req)
	res := rr.Result()
	defer res.Body.Close()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error(err)
	}
	data, _ := io.ReadAll(res.Body)
	expr_test := []Expression_with_result{}
	json.Unmarshal(data, &expr_test)
	if len(expr_test) > 0 {
		req2, err := http.NewRequest("GET", "/results/"+expr_test[0].ID, nil)
		if err != nil {
			t.Error("Error getResults by id")
		}
		AddCookie(req2)
		rr2 := httptest.NewRecorder()
		handleResultsByID(rr2, req2)
		res2 := rr2.Result()
		defer res2.Body.Close()
		data2, _ := io.ReadAll(res2.Body)
		expr_test2 := Expression_with_result{}
		err = json.Unmarshal(data2, &expr_test2)
		if expr_test2.ID != expr_test[0].ID {
			t.Error("Error in TestGetResults")
		}
	}
}

func TestDatabase(t *testing.T) {
	req, err := http.NewRequest("POST", "/database?time=2", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req)
	rr := httptest.NewRecorder()
	handleDatabaseTime(rr, req)
	res := rr.Result()
	defer res.Body.Close()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "/database", nil)
	if err != nil {
		t.Error(err)
	}
	AddCookie(req2)
	rr2 := httptest.NewRecorder()
	handleDatabaseTime(rr2, req2)
	res2 := rr2.Result()
	defer res2.Body.Close()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error(err)
	}
	data2, _ := io.ReadAll(res2.Body)
	if string(data2) != "2\n" {
		t.Error("Error in TestDatabase")
	}
}
