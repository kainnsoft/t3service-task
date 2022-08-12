package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	config "team3-task/config"
	gp "team3-task/internal/controller/grpc"
	v1 "team3-task/internal/controller/http/v1"
	repository "team3-task/internal/repository/db"
	kafkarepo "team3-task/internal/repository/kafka"
	"team3-task/internal/usecase"
	"team3-task/pkg/httpserver"
	kafkaPkg "team3-task/pkg/kafka"
	"team3-task/pkg/logging"
	"team3-task/pkg/pg"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
)

const (
	logForInfo  string = "forInfo"
	logForError string = "forError"
	logForDebug string = "forDebug"
)

func Run(cfg *config.Config) {
	log := logging.New(cfg.Log.Level) // ZeroLogger

	// db
	insPgDB := includePg(cfg, log)
	if insPgDB != nil {
		defer insPgDB.Close()
	}

	// kafka
	brokers := []string{cfg.KafkaAddress}
	kafkaProducerAboutTaskEvent, err := kafkaPkg.New(brokers, cfg.KafkaTopicTask, "")
	kafkaProducerToMailService, err := kafkaPkg.New(brokers, cfg.KafkaTopicMail, "")
	kafkaProducers := kafkarepo.KafkaProducers{
		KafProducerAboutTaskEvent: kafkaProducerAboutTaskEvent,
		KafProducerToMailService:  kafkaProducerToMailService,
	}

	if err != nil {
		log.Error("Can't create kafka Client: %v", err) // Fatal
	}

	// grpc
	var grpcClient *gp.GClient = &gp.GClient{}
	conn, err := grpc.Dial(cfg.GRPC.GRPCAddress, grpc.WithInsecure())
	if err != nil {
		log.Error("Can't create GRPC connection: %v", err) // Fatal
	} else {
		defer conn.Close()
		grpcClient = gp.NewGrpcClient(conn, log)
	}

	// Use case:
	inUseCase := getInUseCase(insPgDB, &kafkaProducers, log)

	// HTTP Server
	mux := http.NewServeMux()
	v1.NewRouter(mux, inUseCase, grpcClient, log)
	httpServer := httpserver.New(mux, cfg.HTTP, log)
	if httpServer != nil {
		log.Info("app - Run - httpServer has run on addr %v", httpServer.GetAddr())
	}

	//------------------------------------

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify: %w", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown: %w", err)
	}
}

func loggerWriter(path string, forErr string) *os.File {
	if path == "osStdOut" {
		return os.Stdout
	}

	if path == "osStdErr" {
		return os.Stderr
	}
	loggerFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o664)
	if err != nil {
		if forErr == logForError {
			return os.Stderr
		} else {
			return os.Stdout
		}
	}

	return loggerFile
}

func includePg(cfg *config.Config, log *logging.ZeroLogger) *pg.DB {
	strurl := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.PG.Username),
		url.QueryEscape(cfg.PG.Password),
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.DBName,
		cfg.PG.ConnTimeout)

	insPgDB, err := pg.NewInsPgDB(strurl, cfg.PG.PoolMax)
	if err != nil {
		log.Error("Can't create DB connection: %v", err) // Fatal

		return nil
	}

	migrationUp(strurl, log)

	return insPgDB
}

func migrationUp(strurl string, log *logging.ZeroLogger) {
	conn, err := sql.Open("postgres", strurl)
	if err != nil {
		log.Error("Can't sql.Open migrarion: %v", err) // Fatal
	}
	err = goose.Up(conn, "migrations")
	if err != nil {
		log.Error("Can't create migrarion: %v", err) // Fatal
	}
}

func getInUseCase(insPgDB *pg.DB, kafkaClient *kafkarepo.KafkaProducers, log *logging.ZeroLogger) *usecase.InUseCase {
	// first get proper repo
	var (
		currentTaskUseCase          usecase.TaskDBRepoInterface
		currentTaskApproversUseCase usecase.TaskApproversDBRepoInterface
		currentTaskEventsUseCase    usecase.TaskEventsDBRepoInterface
		currentUserUseCase          usecase.UserDBRepoInterface
		currentTxUseCase            usecase.TxDBRepoInterface
	)

	if insPgDB != nil {
		currentTaskUseCase = repository.NewTaskPGRepo(insPgDB)
		currentTaskApproversUseCase = repository.NewTaskApproversPGRepo(insPgDB)
		currentTaskEventsUseCase = repository.NewTaskEventsPGRepo(insPgDB)
		currentUserUseCase = repository.NewUserPGRepo(insPgDB)
		currentTxUseCase = repository.NewTxDBRepo(insPgDB)
	} else {
		currentUserUseCase, _ = repository.NewUserMockRepo(log)
	}
	currentInUseCase := usecase.NewInUseCase(
		currentTaskUseCase,
		currentTaskApproversUseCase,
		currentTaskEventsUseCase,
		currentUserUseCase,
		currentTxUseCase,
		kafkaClient,
		log,
	)

	return currentInUseCase
}
