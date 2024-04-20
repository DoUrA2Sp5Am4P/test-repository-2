package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Name     string
	Password string
}

func auth(mux *http.ServeMux) {
	//mux.Handle("/welcome", AuthorizedUsers(http.HandlerFunc(testHandle)))
	mux.Handle("/signin", http.HandlerFunc(Signin))
	mux.Handle("/login", http.HandlerFunc(LoginHandle))
	mux.Handle("/logout", http.HandlerFunc(logoutHandle))
}

func logoutHandle(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
		Path:    "/",
	})
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	body, err := io.ReadAll(r.Body)
	//err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	str_cred := strings.Split(string(body), "&")
	if len(str_cred) > 1 {
		str_login := strings.Split(str_cred[0], "=")[1]
		str_password := strings.Split(str_cred[1], "=")[1]
		creds.Username = str_login
		creds.Password = str_password
	}
	state, current := login(creds.Username, creds.Password)
	if state == false {
		//w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/goback.html", http.StatusSeeOther)
		return
	}
	currentUser = current.Name
	expirationTime := time.Now().Add(800 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
	//fmt.Fprintf(w, "OK Login complete")
}

var currentUser string

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	body, err := io.ReadAll(r.Body)
	//err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	str_cred := strings.Split(string(body), "&")
	if len(str_cred) > 1 {
		str_login := strings.Split(str_cred[0], "=")[1]
		str_password := strings.Split(str_cred[1], "=")[1]
		creds.Username = str_login
		creds.Password = str_password
	}
	currentUser_, err := register(creds.Username, creds.Password)
	if err != nil {
		http.Redirect(w, r, "/usernotexist.html", http.StatusSeeOther)
		return
	}
	currentUser = currentUser_.Name
	operations = append(operations, Operation{Name: "Сложение", Username: currentUser, ExecutionTime: 1})
	operations = append(operations, Operation{Name: "Вычитание", Username: currentUser, ExecutionTime: 2})
	operations = append(operations, Operation{Name: "Умножение", Username: currentUser, ExecutionTime: 3})
	operations = append(operations, Operation{Name: "Деление", Username: currentUser, ExecutionTime: 4})
	sl = append(sl, sleep{Username: currentUser, time: 1})
	expirationTime := time.Now().Add(800 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
	//fmt.Fprintf(w, "OK Register complete")
}
func AuthorizedUsers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var Servehttpvar int = 0
		if r.URL.Path == "/login.html" || r.URL.Path == "/register.html" || r.URL.Path == "/favicon.ico" || r.URL.Path == "/goback.html" || r.URL.Path == "/jquery-3.6.0.min.js" || r.URL.Path == "/usernotexist.html" {
			Servehttpvar = 1
			next.ServeHTTP(w, r)
		}
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/register.html", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/register.html", http.StatusSeeOther)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "/register.html", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/register.html", http.StatusSeeOther)
			return
		}
		if !tkn.Valid {
			http.Redirect(w, r, "/register.html", http.StatusSeeOther)
			return
		}
		currentUser = claims.Username
		if Servehttpvar != 1 {
			next.ServeHTTP(w, r)
		}
	})
}
func createTable(ctx context.Context, db *sql.DB) error {
	const usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT UNIQUE,
		password TEXT
	);`

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	return nil
}

func insertUser(ctx context.Context, db *sql.DB, user *User) (int64, error) {
	var q = `
	INSERT INTO users (name, password) values ($1, $2)
	`
	result, err := db.ExecContext(ctx, q, user.Name, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func selectUser(ctx context.Context, db *sql.DB, name string) (User, error) {
	var (
		user User
		err  error
	)

	var q = "SELECT id, name, password FROM users WHERE name=$1"
	err = db.QueryRowContext(ctx, q, name).Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func login(userInput string, passwordInput string) (bool, User) {
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
	userFromDB, err := selectUser(ctx, db, userInput)
	if err != nil {
		return false, User{ID: 0, Name: "", Password: ""}
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(passwordInput))
	if err == nil {
		return true, userFromDB
	} else {
		return false, User{ID: 0, Name: "", Password: ""}
	}
}
func register(userInput string, passwordInput string) (User, error) {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	err = db.PingContext(ctx)
	if err != nil {
		return User{}, err
	}
	if err = createTable(ctx, db); err != nil {
		return User{}, err
	}

	password, err := generate(passwordInput)
	if err != nil {
		return User{}, err
	}

	user := &User{
		Name:     userInput,
		Password: password,
	}
	userID, err := insertUser(ctx, db, user)
	if err != nil {
		log.Println("user already exists")
		return User{}, fmt.Errorf("user already exists")
	} else {
		user.ID = userID
	}
	return *user, nil
}
