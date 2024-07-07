package createnewconnection

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/config"

	connector "github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/postgres/connector"
)

type CreaterNewConnector interface {
	CreaterConnector(config *connector.Config) (*connector.Connector, error)
}

type SaverNewConnector interface {
	SaveConnector(c *connector.Connector) error
}

type Response struct {
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
}

type Request struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
}

func NewCreateNewConnectorHandler(
		logger *slog.Logger, cfg *config.Config, 
		c CreaterNewConnector, s SaverNewConnector,
	) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Переменная для запроса
		var request Request
		// Декодируем запрос
		if err := render.DecodeJSON(r.Body, &request); err != nil {
			logger.Error("Decoding request body was failed", err.Error())
			render.JSON(w, r, Response{Status: "Error", Error: "Decoding request body was failed"})
			return
		}

		// Пробуем создать подключение
		conn, err := c.CreaterConnector(
			&connector.Config{
				Host: request.Host,
				Port: request.Port,
				User: request.User,
				Password: request.Password,
				DBname: request.DBname,
			},
		)
		if err != nil {
			logger.Error("Create connector was failed", err.Error())
			render.JSON(w, r, Response{Status: "Error", Error: "Decoding request body was failed"})
		}

		// Сохраняем подключение 
		err = s.SaveConnector(conn)
		if err != nil {
			logger.Error("Create connector was failed", err.Error())
			render.JSON(w, r, Response{Status: "Error", Error: "Decoding request body was failed"})
		}

		// Отправляем отчем об успешности создания подключения
		render.JSON(w, r, Response{
			Status:   "OK",
			Message:  fmt.Sprintf(`Connection to database "%v" was successfully`, request.DBname),
		})
	}
}