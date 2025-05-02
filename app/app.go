package app

import (
	"context"
	"database/sql"
	"echo/tui"
	"echo/tui/styles"
	"errors"
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
	"github.com/muesli/termenv"
	_ "modernc.org/sqlite"
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
	db, err := sql.Open("sqlite", "file:echo.db")
	if err != nil {
		log.Fatalf("Failed to open sqlite db: %v", err)
	}

	echoDB.RunMigration(db)

	log.Printf("Pinging database...")
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Printf("✅ Database connection successful.")

	userRepo := repo.NewSQLiteUserRepository(db)
	echoDB.RunUserSeed(context.Background(), userRepo)

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

		m := tui.InitialRootModel(repo.NewSQLiteUserRepository(a.db))

		return tea.NewProgram(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
	}

	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
