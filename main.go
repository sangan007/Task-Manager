package main
import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)
type Task struct{
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
var db *gorm.DB
func init(){
    dsn := "host=localhost user=postgres password=sangan007 dbname=task_manager port=5432 sslmode=disable"
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    db.AutoMigrate(&Task{})
}
func CreateTask(w http.ResponseWriter, r *http.Request){
    var task Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := db.Create(&task).Error; err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}
func GetTasks(w http.ResponseWriter, r *http.Request){
    var tasks []Task
    db.Find(&tasks)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}
func GetTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    params := mux.Vars(r)
    if err := db.First(&task, params["id"]).Error; err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(task)
}
func UpdateTask(w http.ResponseWriter, r *http.Request){
    var task Task
    params := mux.Vars(r)
    if err := db.First(&task, params["id"]).Error; err != nil{
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil{
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    db.Save(&task)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(task)
}
func DeleteTask(w http.ResponseWriter, r *http.Request){
    var task Task
    params := mux.Vars(r)
    if err := db.First(&task, params["id"]).Error; err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    db.Delete(&task)
    w.WriteHeader(http.StatusNoContent)
}
func GetStatistics(w http.ResponseWriter, r *http.Request){
    var totalTasks, completedTasks int64
    db.Model(&Task{}).Count(&totalTasks)
    db.Model(&Task{}).Where("status = ?", "completed").Count(&completedTasks)
   stats := map[string]int64{
        "total_tasks":totalTasks,
        "completed_tasks": completedTasks,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(stats)
}
func main(){
    r := mux.NewRouter()
    r.HandleFunc("/tasks", CreateTask).Methods("POST")
    r.HandleFunc("/tasks", GetTasks).Methods("GET")
    r.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
    r.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
    r.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
    log.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
