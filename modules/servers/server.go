package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/MarkTBSS/go-CORS/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func NewServer(cfg config.IConfig) IServer {
	return &server{
		cfg: cfg,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

type IServer interface {
	Start()
	GetServer() *server
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {
	// Middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Cors())

	// Modules
	v1 := s.app.Group("v1")
	modules := InitModule(v1, s, middlewares)
	modules.MonitorModule()

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	log.Printf("server is starting on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
