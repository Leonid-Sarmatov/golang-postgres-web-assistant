package connector

import (
	"fmt"
	"log"
	"database/sql"

	_ "github.com/lib/pq"
)

/*
Соединение к базе данных
*/
type PostgresConnector struct {
	Name string
	Status string
	Connection *sql.DB
	*Config
}

type Config struct {
	Host string
	Port string
	User string
	Password string
	DBname string
}

func NewPostgresConnector(config *Config) (*PostgresConnector, error) {
	var pc PostgresConnector

	// Формируем строку для подключения
	connectString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBname,
	)

	// Пробуем создать соединение с базой данных
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Printf("Spawn connection to database was failed %v", err)
		return nil, err
	}

	pc.Config = config
	pc.Connection = db
	return &pc, nil
}

func (pc *PostgresConnector) RequestWithResponse(request string) (*sql.Rows, error) {
	rows, err := pc.Connection.Query(request)
	if err != nil {
		return nil, err
	}

	rows.Columns()
	x, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	x[0].ScanType()
	return pc.Connection.Query(request)
}

func (pc *PostgresConnector) RequestWithoutResponse(request string) (sql.Result, error) {
	return pc.Connection.Exec(request)
}

func (pc *PostgresConnector) IsAlive() bool {
	err := pc.Connection.Ping()
	if err != nil {
		pc.Status = "Dead"
		return false
	}
	return true
}

func (pc *PostgresConnector) CloseConnection() {
	pc.Connection.Close()
	pc.Status = "Was closed"
}