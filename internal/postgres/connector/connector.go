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
type Connector struct {
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

func NewConnector(config *Config) (*Connector, error) {
	var c Connector

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

	c.Config = config
	c.Connection = db
	return &c, nil
}

func (c *Connector) RequestWithResponse(request string) (*sql.Rows, error) {	
	return c.Connection.Query(request)
}

func (c *Connector) RequestWithoutResponse(request string) (sql.Result, error) {
	return c.Connection.Exec(request)
}

func (c *Connector) IsAlive() bool {
	err := c.Connection.Ping()
	if err != nil {
		c.Status = "Dead"
		return false
	}
	return true
}

func (c *Connector) CloseConnection() {
	c.Connection.Close()
	c.Status = "Was closed"
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
