package storage

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/Courtcircuits/mitter-server/types"
	"github.com/Courtcircuits/mitter-server/util"
	_ "github.com/lib/pq"
)

type Database struct {
	user     string
	password string
	host     string
	port     string
	database string
}

func NewDatabase() *Database {
	return &Database{
		util.Get("DB_USER"),
		util.Get("DB_PASSWORD"),
		util.Get("DB_HOST"),
		util.Get("DB_PORT"),
		util.Get("DB_DATABASE"),
	}
}

func (db *Database) connect() (*sql.DB, error) {
	connStr := "user=" + db.user + " password=" + db.password + " host=" + db.host + " port=" + db.port + " dbname=" + db.database
	client, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// create a user in the database and intance it, hash the password
func (db *Database) CreateUser(name string, password string) (types.User, error) {
	password = util.Hash(password)

	client, err := db.connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `INSERT INTO users(name, password) VALUES($1, $2) RETURNING id, name, password;`
	user, err := types.ScanUser(client.QueryRow(query, name, password))

	return user, err
}

func (db *Database) FindUser(name string, password string) (types.User, error) {
	password = util.Hash(password)

	client, err := db.connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `SELECT id, name, password FROM users WHERE name=$1 AND password=$2`
	user, err := types.ScanUser(client.QueryRow(query, name, password))

	return user, err
}

func (db *Database) GetUser(id int) (types.User, error) {
	client, err := db.connect()

	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `--sql SELECT id, name, password FROM users WHERE id=$1`
	user, err := types.ScanUser(client.QueryRow(query, id))

	return user, err
}

func (db *Database) CreateMessage(content string, id_owner int, name_owner string) (types.Message, error) {
	log.Printf("User : %d, %q sends message %q\n", id_owner, name_owner, content)

	timestamp := time.Now()

	client, err := db.connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `INSERT INTO message(content, timestamp, id_owner) VALUES($1, $2, $3) RETURNING id, content, timestamp, id_owner;`
	message, err := types.ScanMessage(client.QueryRow(query, content, timestamp, id_owner), name_owner)
	return message, err
}

func (db *Database) GetMessages() ([]types.Message, error) {
	client, err := db.connect()

	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `SELECT m.id, m.content, m.timestamp, m.id_owner, o.name FROM message m JOIN users o ON m.id_owner=o.id order by m.timestamp desc limit 50`

	type message_db struct {
		id         int
		content    sql.NullString
		timestamp  sql.NullTime
		id_owner   int
		name_owner sql.NullString
	}

	rows, err := client.Query(query)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var messages []types.Message

	for rows.Next() {
		var message message_db
		if err := rows.Scan(&message.id, &message.content, &message.timestamp, &message.id_owner, &message.name_owner); err != nil {
			return messages, err
		}
		messages = append(messages, types.Message{
			ID:         message.id,
			Content:    message.content.String,
			Timestamp:  strconv.FormatInt(message.timestamp.Time.Unix(), 10),
			Name_owner: message.name_owner.String,
		})
	}

	return messages, nil
}

func (db *Database) GetMessagesSince(date time.Time) ([]types.Message, error) {
	client, err := db.connect()

	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `SELECT m.id, m.content, m.timestamp, m.id_owner, o.name FROM message m JOIN users o ON m.id_owner=o.id WHERE EXTRACT(EPOCH FROM m.timestamp) > $1 order by m.timestamp desc limit 50`

	type message_db struct {
		id         int
		content    sql.NullString
		timestamp  sql.NullTime
		id_owner   int
		name_owner sql.NullString
	}

	rows, err := client.Query(query, date.UnixMilli())

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var messages []types.Message

	for rows.Next() {
		var message message_db
		if err := rows.Scan(&message.id, &message.content, &message.timestamp, &message.id_owner, &message.name_owner); err != nil {
			return messages, err
		}
		messages = append(messages, types.Message{
			ID:         message.id,
			Content:    message.content.String,
			Timestamp:  strconv.FormatInt(message.timestamp.Time.Unix(), 10),
			Name_owner: message.name_owner.String,
		})
	}

	return messages, nil
}
