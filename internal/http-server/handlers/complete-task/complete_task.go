package complete_task

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"todolist/internal/lib/api/response"
	"todolist/internal/storage"
)

type Request struct {
	TaskId int `json:"task-id" validate:"required"`
}

type Response struct {
	response.Response
}

type TaskCompleter interface {
	CompleteTask(id int, login string, newTaskStatus storage.TaskStatus) error
}

func New(log *slog.Logger, taskCompleter TaskCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.complete-task.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		login := r.Context().Value("login")
		if login == nil {
			log.Error("empty login from context")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty login from context"))
			return
		}
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request")
			validateError := err.(validator.ValidationErrors)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateError))
			return
		}

		err = taskCompleter.CompleteTask(req.TaskId, fmt.Sprint(login), storage.StatusDone)
		if err != nil {
			log.Error("failed change task status", err)
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, response.Error("failed change task status"))
			return
		}

		render.JSON(w, r, responseOk())
	}
}

func responseOk() Response {
	return Response{
		response.OK(),
	}
}
