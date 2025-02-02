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

// Ниже напишите обработчики для каждого эндпоинта

//обработчик для получения всех задач
//Конечная точка /tasks.
//Метод GET.
//При успешном запросе сервер должен вернуть статус 200 OK.
//При ошибке сервер должен вернуть статус 500 Internal Server Error.

func getTask(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

//Обработчик для отправки задачи на сервер
//Конечная точка /tasks.
//Метод POST.
//При успешном запросе сервер должен вернуть статус 201 Created.
//При ошибке сервер должен вернуть статус 400 Bad Request.

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	w.Header().Set("Contant_Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

//Обработчик для получения задачи по ID
//В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.
//Конечная точка /tasks/{id}.
//Метод GET.
//При успешном выполнении запроса сервер должен вернуть статус 200 OK.
//В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.

func getId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tasks, ok := tasks[id]

	if !ok {
		http.Error(w, "ID не найден", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

//Обработчик удаления задачи по ID
//Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе, если нет вернуть соответствующий статус.
//Конечная точка /tasks/{id}.
//Метод DELETE.
//При успешном выполнении запроса сервер должен вернуть статус 200 OK.
//В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
//Во всех обработчиках тип контента Content-Type — application/json.

func delId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tasks, ok := tasks[id]

	if !ok {
		http.Error(w, "ID не найден", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTask)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getId)
	r.Delete("/tasks/{id}", delId)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
