package add_task

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
	"todolist/internal/lib/api/response"
)

type Request struct {
	Title       string
	Description string
	Deadline    time.Time
}

type Response struct {
	response.Response
}

type TaskAdder interface {
	AddTask(userLogin, title, description string, deadline time.Time) error
}

func New(log *slog.Logger, taskAdder TaskAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.add-task.New"
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

		err = taskAdder.AddTask(fmt.Sprint(login), req.Title, req.Description, req.Deadline)
		if err != nil {
			log.Error("failed to add task")
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, response.Error("failed to add task"))
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
