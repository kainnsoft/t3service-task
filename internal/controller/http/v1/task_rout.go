package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	gp "team3-task/internal/controller/grpc"
	_ "team3-task/internal/docs"
	"team3-task/internal/entity"
	"team3-task/internal/errors"
	"team3-task/pkg/logging"

	httpSwagger "github.com/swaggo/http-swagger"
)

type taskRoutes struct {
	logger      *logging.ZeroLogger
	taskHandler TaskHandlerInterface
	grpcClient  *gp.GClient
}

func NewTaskRouter(mux *http.ServeMux, t TaskHandlerInterface, grpcClient *gp.GClient, log *logging.ZeroLogger) {
	rout := taskRoutes{log, t, grpcClient}

	// swagger
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.Handle("/ping", rout.Ping()) // GET
	mux.Handle("/task", rout.Task())
	// mux.Handle("/create", rout.Create())            task        POST
	// mux.Handle("/update", rout.Update())            task/{id}   PUT(id)
	// mux.Handle("/delete", rout.Delete())            task/{id}   DELETE(id)
	// mux.Handle("/get", rout.Get())                  task/{id}   GET(id)
	// mux.Handle("/list", rout.List())                task        GET

	mux.Handle("/approvetask", rout.ApproveTask()) // PATCH(id)
	mux.Handle("/rejecttask", rout.RejectTask())   // PATCH(id)
}

// Ping godoc
// @Summary accessailability task.service
// @Produce json
// @Success 200
// @Router /ping [get]
func (rout taskRoutes) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK, I'm ready to read and write from task-service"))

		if err != nil {
			rout.logger.Error("v1.rout Ping error: %v", err)
		}
	}
}

func (rout *taskRoutes) Task() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasID := r.URL.Query().Has("id")
		switch hasID {
		case true:
			curTaskID := r.URL.Query().Get("id")

			switch r.Method {
			case http.MethodGet:
				rout.Get(curTaskID)
			case http.MethodPut:
				rout.Update()
			case http.MethodDelete:
				rout.Delete()
			default:
				rout.MethodNotAllowed(w)
			}
		case false:
			switch r.Method {
			case http.MethodGet:
				rout.List(w, r)
			case http.MethodPost:
				rout.Create(w, r)
			default:
				rout.MethodNotAllowed(w)
			}
		}
	}
}

// CreateTask godoc
// @Summary create task
// @Description add (create) new task
// @Accept json
// @Produce json
// @Param article body entity.Task true "New Task"
// @Success 201 {string} entity.Task.Id
// @Failure 400 {string} string
// @Failure 401 {object} errors.customTaskError
// @Failure 404 {string} string
// @Failure 405 {object} errors.customTaskError
// @Failure 500 {string} string
// @Failure 503 {string} string
// @Router /task [post]
func (rout taskRoutes) Create(w http.ResponseWriter, r *http.Request) {
	// token validation from grpc auth.service // TODO test
	// validationAuthResponse, err := rout.checkValidation(r)
	// if err != nil {
	// 	rout.logger.Error("%v", err)
	// 	http.Error(w, "Validation error occures. Please try log in again.", http.StatusUnauthorized)
	// 	return
	// }
	validationAuthResponse := entity.AuthResponse{Username: "author@gmail.com"}

	// If we're right here, validation has successed and we got user email
	body, err := io.ReadAll(r.Body)
	if err != nil {
		rout.logger.Error("v1.rout Create: ioutil.ReadAll(Body) %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	defer r.Body.Close()

	idCreatedTask, err := rout.taskHandler.CreateTaskHandle(r.Context(), body, validationAuthResponse.Username)

	if err != nil {
		rout.logger.Error("v1.rout Create: rout.taskHandler.CreateTaskHandle error: %v", err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

		return
	}

	strAnswer := fmt.Sprintf("created with id: %d", idCreatedTask)
	rout.logger.Info("v1.rout Create: task %v", strAnswer)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Success! Task " + strAnswer))

	if err != nil {
		rout.logger.Error("v1.rout Create error: %v", err)
	}
}

func (rout *taskRoutes) Update() http.HandlerFunc { // TODO
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK, I'm ready to read and write from task-service"))

		if err != nil {
			rout.logger.Error("rout.Update() error: %v", err)
		}
	}
}

func (rout *taskRoutes) Delete() http.HandlerFunc { // TODO
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK, I'm ready to read and write from task-service"))

		if err != nil {
			rout.logger.Error("rout.Delete() error: %v", err)
		}
	}
}

func (rout *taskRoutes) Get(taskID string) http.HandlerFunc { // TODO
	return func(w http.ResponseWriter, r *http.Request) {
		// token validation from grpc auth.service
		validationAuthResponse, err := rout.checkValidation(r)
		if err != nil {
			rout.logger.Error("rout.Get() validation error: %v", err)
			http.Error(w, "Validation error occures. Please try log in again.", http.StatusUnauthorized)

			return
		}

		intTaskID, err := strconv.Atoi(taskID)

		if err != nil {
			rout.logger.Error("rout.Get() not valid id error: %v, id = %s", err, taskID)
			http.Error(w, "Not valid task id. Please try id parametr again.", http.StatusBadRequest)

			return
		}

		task, err := rout.taskHandler.GetTaskHandle(context.Background(), intTaskID)

		if err != nil {
			rout.logger.Error("rout.Get() rout.taskHandler.GetTaskHandle error: %v; intTaskId = %d", err, intTaskID)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
		// если запрашивающи информацию не автор задачи, то отбой
		if validationAuthResponse.Username != task.Author.Email {
			rout.logger.Error("rout.Get() requesting information is not task author: %v; intTaskId = %d; author = %s",
				validationAuthResponse.Username, intTaskID, task.Author.Email)
			http.Error(w, "Only author can get task information", http.StatusNotAcceptable)

			return
		}

		resp, err := json.MarshalIndent(task, "", "  ")

		if err != nil {
			rout.logger.Error("rout.Get() json.Marshal error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)

		if err != nil {
			rout.logger.Error("rout.Get() error: %v", err)
		}
	}
}

func (rout *taskRoutes) List(w http.ResponseWriter, r *http.Request) { // TODO
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK, I'm ready to read and write from task-service"))

	if err != nil {
		rout.logger.Error("rout.List() error: %v", err)
	}
}

func (rout *taskRoutes) ApproveTask() http.HandlerFunc { // TODO
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK, I'm ready to read and write from task-service"))
		if err != nil {
			rout.logger.Error("rout.ApproveTask() error: %v", err)
		}
	}
}

func (rout *taskRoutes) RejectTask() http.HandlerFunc { // TODO
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK, I'm ready to read and write from task-service"))
		if err != nil {
			rout.logger.Error("rout.RejectTask() error: %v", err)
		}
	}
}

func (rout *taskRoutes) checkValidation(r *http.Request) (entity.AuthResponse, error) {
	validationAuthResponse := entity.AuthResponse{}
	authRequest := entity.AuthRequest{}
	accessToken, err := r.Cookie("access")
	if err != nil {
		err := errors.Wrapf(err, " validation error occures while cookie getting")

		return validationAuthResponse, err
	} else {
		authRequest.AccessToken = accessToken.Value
	}

	validationAuthResponse, err = rout.grpcClient.CheckAccess(&authRequest)
	if (err != nil) || (validationAuthResponse.Error != "") {
		err := errors.Wrapf(err, " validation error occures")

		return validationAuthResponse, err
	}

	return validationAuthResponse, nil
}

func (rout taskRoutes) MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (rout *taskRoutes) HandleError(w http.ResponseWriter, err error) {
	var status int

	errorType := errors.GetType(err)

	switch errorType {
	case errors.BadRequest:
		status = http.StatusBadRequest
	case errors.NotFound:
		status = http.StatusNotFound
	case errors.MethodNotAllowed:
		status = http.StatusMethodNotAllowed
	case errors.NoType:
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)

	var respStr string

	errorContext := errors.GetErrorContext(err)
	if errorContext != nil {
		respStr = fmt.Sprintf("context %v", errorContext)
	} else {
		respStr = fmt.Sprintf("error %s", err.Error())
	}

	http.Error(w, respStr, status)
}
