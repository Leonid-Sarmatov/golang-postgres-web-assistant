package connector

import (
	"fmt"
	//"reflect"
	"log"
	"unicode"

	//"time"
	//"reflect"
	"database/sql"
	//"os"
	//"os/exec"

	//"plugin"
	//"text/template"

	_ "github.com/lib/pq"
)

/*
Соединение к базе данных
*/
type PostgresConnector struct {
	Name       string
	Status     string
	Connection *sql.DB
	*Config
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
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

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func SqlRowsToSliceOfMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	// Получаем имя столбцов
	columnsNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := make([]map[string]interface{}, 0)
    for rows.Next() {
		row := make([]interface{}, len(columnsNames))
        if err := rows.Scan(row...); err != nil {
            return nil, err
        }
		rowMap := make(map[string]interface{})
		for i, val := range row {
			rowMap[columnsNames[i]] = val
		} 
		res = append(res, rowMap)
    }
	return res, nil
}
