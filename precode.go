package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func GetTasks(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func PostTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

}

func GetTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(res, "задача с таким id не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)

}

func DeleteTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(req, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(res, "задача с таким id не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", GetTasks)
	r.Get("/tasks/{id}", GetTask)
	r.Post("/tasks", PostTask)
	r.Delete("/tasks/{id}", DeleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

/*
Обработчик для получения всех задач
Обработчик должен вернуть все задачи, которые хранятся в мапе.
Конечная точка /tasks.
Метод GET.
При успешном запросе сервер должен вернуть статус 200 OK.
При ошибке сервер должен вернуть статус 500 Internal Server Error.

Обработчик для отправки задачи на сервер
Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
Конечная точка /tasks.
Метод POST.
При успешном запросе сервер должен вернуть статус 201 Created.
При ошибке сервер должен вернуть статус 400 Bad Request.

Обработчик для получения задачи по ID
Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.
В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.
Конечная точка /tasks/{id}.
Метод GET.
При успешном выполнении запроса сервер должен вернуть статус 200 OK.
В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.

Обработчик удаления задачи по ID
Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе, если нет вернуть соответствующий статус.
Конечная точка /tasks/{id}.
Метод DELETE.
При успешном выполнении запроса сервер должен вернуть статус 200 OK.
В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.


Во всех обработчиках тип контента Content-Type — application/json.

*/
