package app

import (
	"context"
	"database/sql"
	"echo/tui"
	"echo/tui/styles"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	echoDB "echo/db"
	repo "echo/db/repository"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/muesli/termenv"
)

const (
	host = "0.0.0.0"
	port = "4242"
)

type App struct {
	*ssh.Server
	db *sql.DB
	// CentralHubReqChan chan clientHubReq
	// RoomHubNotifChan chan roomHubNotif

	// ClientRoomNotifChan chan clientRoomNotif
}

func NewApp() *App {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file, relying on environment variables", "error", err)
	}

	app := new(App)
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			app.echoMiddleware(),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	app.Server = s
	return app
}

func (a *App) dbSetup() {
	// Construct DSN from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Basic validation
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Database configuration environment variables are not fully set.")
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
		log.Warn("DB_SSLMODE not set, defaulting to 'disable'. Ensure this is secure for production.")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Or using the postgres:// URL format if preferred by pgx/stdlib
	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	// 	dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	log.Info("Connecting to PostgreSQL database...")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open PostgreSQL database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Minute)

	log.Info("Pinging database...")
	// Ping with context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ Failed to ping PostgreSQL database: %v", err)
	}
	log.Info("✅ PostgreSQL Database connection successful.")

	// Run migrations (pass the DSN for migrate to use if needed, or the db instance)
	echoDB.RunMigration(db) // Pass DSN if your RunMigration needs it for NewWithDatabaseInstance

	userRepo := repo.NewPostgresUserRepository(db) // Placeholder for now
	echoDB.RunUserSeed(context.Background(), userRepo) // Seeding will also need review

	a.db = db
}

// this function is responsible of starting and handling shutting down the server
// gracefully after recieving interuption signal or crashing, what ever
// the case it will help do the nassisay clean up.
func (a *App) Start() {

	a.dbSetup()

	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = a.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}

func (a *App) echoMiddleware() wish.Middleware {

	teaHandler := func(s ssh.Session) *tea.Program {
		// pty, _, active := s.Pty()
		// if !active {
		// 	wish.Fatalln(s, "no active terminal, skipping")
		// 	return nil
		// }

		styles.ClientRenderer = bubbletea.MakeRenderer(s)

		m := tui.InitialRootModel(repo.NewPostgresUserRepository(a.db))

		return tea.NewProgram(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
	}

	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
