package get_tasks

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
	TaskStatus string `json:"task_status" validate:"required,oneof=done new all"`
}

type Response struct {
	response.Response
}

type TaskGetter interface {
	GetAllTasks(login string, status string) ([]storage.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get-tasks.New"
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

		if req.TaskStatus == "all" {
			req.TaskStatus = ""
		}
		tasks, err := taskGetter.GetAllTasks(fmt.Sprint(login), req.TaskStatus)
		if err != nil {
			log.Error("failed get tasks", err)
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, response.Error("failed to get tasks"))
			return
		}

		render.JSON(w, r, tasks)
	}
}

func responseOk() Response {
	return Response{
		response.OK(),
	}
}
