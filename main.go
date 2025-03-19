package main //m@8246.ru

import (
	"embed"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

//go:embed static
var staticFiles embed.FS

func main() {
	r := mux.NewRouter()

	// Обслуживание статических файлов
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	// Маршруты
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/qr", generateQR).Methods("GET")

	// Добавляем обработчик 404
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// Проверяем, существует ли файл
	_, err := staticFiles.Open("static/index.html")
	if err != nil {
		log.Printf("Ошибка: файл index.html не найден в embedded FS: %v", err)
		http.Error(w, "Файл index.html не найден", http.StatusInternalServerError)
		return
	}
	http.ServeFileFS(w, r, staticFiles, "static/index.html")
}

func generateQR(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "Параметр 'text' обязателен", http.StatusBadRequest)
		return
	}

	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		http.Error(w, "Ошибка при генерации QR-кода", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=\"qrcode.png\"")

	err = qr.Write(256, w)
	if err != nil {
		http.Error(w, "Ошибка при записи QR-кода", http.StatusInternalServerError)
		return
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Страница не найдена", http.StatusNotFound)
}
