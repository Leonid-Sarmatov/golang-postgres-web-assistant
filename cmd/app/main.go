package main

import (
	"fmt"
	"log"
	"net/http"

	//"time"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	config "github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/config"

	web_ui "github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/handlers/web_ui"
	cors_headers "github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/middlewares/cors_headers"
	create_new_connection "github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/handlers/create_new_connection"
	send_sql_request "github.com/Leonid-Sarmatov/golang-postgres-web-assistant/internal/handlers/send_sql_request"
)

func main() {
	// Задаем путь до конфигов   /golang_yandex_last_figth/orchestrator_server
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// Инициализируем конфиги
	cfg := config.MustLoad()

	// Инициализируем логгер
	logger := setupLogger(cfg.EnvMode)
	logger.Debug("Successful read configurations.", slog.Any("cfg", cfg))

	// Инициализируем роутер
	router := chi.NewRouter()

	// Подключаем готовый middleware для логирования запросов
	router.Use(middleware.Logger)

	// Подключаем готовый middleware, который отлавливает возможные паники,
	// что бы избежать падение приложения
	router.Use(middleware.Recoverer)

	// Подключаем свой moddleware, который подключает CORS заголовки
	// что бы исключить возможные неполадки со стороны браузера
	router.Use(cors_headers.AddCorsHeaders())

	// Эндпоинт для веб панели управления
	router.Get("/webUI", web_ui.NewLoginSiteHandler(logger, cfg))

	router.Route("/api", func(r chi.Router) {
		// Эндпоинт принимающий запрос на создание нового подключения к базе данных
		r.Get("/createNewConnection", create_new_connection.NewCreateNewConnectorHandler(logger, cfg, ))

		// Эндпоинт принимающий запрос на удаление подключения к базе данных
		//r.Get("/deleteConnection")

		// Эндпоинт принимающий запрос на проверку жизни подключения
		//r.Get("/isConnectionAlive")

		// Эндпоинт принимающий запрос к базе данных
		r.Post("/sendSqlRequest", send_sql_request)
	})

	// Создаем сервер
	server := &http.Server{
		Addr:         cfg.HTTPServerConfig.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerConfig.RequestTimeout,
		WriteTimeout: cfg.HTTPServerConfig.RequestTimeout,
		IdleTimeout:  cfg.HTTPServerConfig.ConnectionTimeout,
	}

	fmt.Println("                Y.                      _             \n" +
		"                YiL                   .```.           \n" +
		"                Yii;                .; .;;`.          \n" +
		"                YY;ii._           .;`.;;;; :          \n" +
		"                iiYYYYYYiiiii;;;;i` ;;::;;;;          \n" +
		"            _.;YYYYYYiiiiiiYYYii  .;;.   ;;;          \n" +
		"         .YYYYYYYYYYiiYYYYYYYYYYYYii;`  ;;;;          \n" +
		"       .YYYYYYY$$YYiiYY$$$$iiiYYYYYY;.ii;`..          \n" +
		"      :YYY$!.  TYiiYY$$$$$YYYYYYYiiYYYYiYYii.         \n" +
		"      Y$MM$:   :YYYYYY$! `` 4YYYYYiiiYYYYiiYY.        \n" +
		"   `. :MM$$b.,dYY$$Yii  :'   :YYYYllYiiYYYiYY         \n" +
		"_.._ :`4MM$!YYYYYYYYYii,.__.diii$$YYYYYYYYYYY         \n" +
		".,._ $b`P`      4$$$$$iiiiiiii$$$$YY$$$$$$YiY;        \n" +
		"   `,.`$:       :$$$$$$$$$YYYYY$$$$$$$$$YYiiYYL       \n" +
		"     `;$$.    .;PPb$`.,.``T$$YY$$$$YYYYYYiiiYYU:      \n" +
		"    ;$P$;;: ;;;;i$y$ !Y$$$b;$$$Y$YY$$YYYiiiYYiYY      \n" +
		"    $Fi$$ .. ``:iii.`- :YYYYY$$YY$$$$$YYYiiYiYYY      \n" +
		"    :Y$$rb ````  `_..;;i;YYY$YY$$$$$$$YYYYYYYiYY:     \n" +
		"     :$$$$$i;;iiiiidYYYYYYYYYY$$$$$$YYYYYYYiiYYYY.    \n" +
		"      `$$$$$$$YYYYYYYYYYYYY$$$$$$YYYYYYYYiiiYYYYYY    \n" +
		"      .i!$$$$$$YYYYYYYYY$$$$$$YYY$$YYiiiiiiYYYYYYY    \n" +
		"     :YYiii$$$$$$$YYYYYYY$$$$YY$$$$YYiiiiiYYYYYYi'    ")

	// Запускаем сервер
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server was stoped")
	}
}

/*
setupLogger инициализирует логер
*/
func setupLogger(envMode string) *slog.Logger {
	var logger *slog.Logger

	switch envMode {
	case "local":
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prodaction":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
