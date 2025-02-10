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

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

//Обработчик для отправки задачи на сервер
//Конечная точка /tasks.
//Метод POST.
//При успешном запросе сервер должен вернуть статус 201 Created.
//При ошибке сервер должен вернуть статус 400 Bad Request.

func addTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contant_Type", "application/json")
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	w.WriteHeader(http.StatusCreated)
}

//Обработчик для получения задачи по ID
//В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.
//Конечная точка /tasks/{id}.
//Метод GET.
//При успешном выполнении запроса сервер должен вернуть статус 200 OK.
//В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	tasks, ok := tasks[id]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(resp)

	delete(tasks, id)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)
	//
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
