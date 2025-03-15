package restdb

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func addUser(username string, password string, admin int, active int, conn *pgx.Conn) error {
	exists, err := getExists(username, conn)
	if err != nil {
		return err
	}
	if exists {
		return err
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO users (username, password, lastlogin, admin, active) VALUES ($1, $2, $3, $4, $5);",
		username, password, 1, admin, active)

	if err != nil {
		return err
	}
	return nil
}
func addNote(username string, article string, hidden bool, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "insert into Notes (username,article,hide) values ($1,$2,$3);", username, article, hidden)
	if err != nil {
		return err
	}
	return nil
}
func getPass(username string, conn *pgx.Conn) (string, error) {
	var password string
	err := conn.QueryRow(context.Background(), "select password from Users where username = $1;", username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func getExists(username string, conn *pgx.Conn) (bool, error) {
	var exists bool
	err := conn.QueryRow(context.Background(),
		"select exists(select 1 from users where username = $1);", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func getNotes(username string, conn *pgx.Conn) ([]string, error) {
	var notes []string
	rows, err := conn.Query(context.Background(), "SELECT article FROM notes WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var note string
		if err := rows.Scan(&note); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func getDB(host string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), host)
	if err != nil {
		return nil, err
	}
	return conn, err
}
func main() {
	host := "postgres://postgres:changeme@10.10.10.100:5432/restapi"
	// conn, err := pgx.Connect(context.Background(), host)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to DB %s: %v\n", host, err)
	// 	os.Exit(1)
	// }
	conn, err := getDB(host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Добавление пользователя для проверки
	err = addUser("testuser", "securepassword", 1, 1, conn)
	if err != nil {
		fmt.Println(err)
	}
	password, err := getPass("admin", conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(password)
	err = addNote("testuser", "testuser article", true, conn)
	if err != nil {
		fmt.Println(err)
	}
	notes, err := getNotes("testuser", conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(notes)
}
