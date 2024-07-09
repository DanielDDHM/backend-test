package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
)

func GetBreeds(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    species := params.Get("species")
    petSizes := params.Get("petSizes")

    var breeds []models.Breed

    db, err := gorm.Open(mysql.Open(MysqlDSN), &gorm.Config{})
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        return
    }

    query := db.Model(&models.Breed{})

    if species != "" {
        query = query.Where("species = ?", species)
    }

    if petSizes != "" {
        query = query.Where("pet_size = ?", petSizes)
    }

    if err := query.Find(&breeds).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(breeds)
}

func GetBreed(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
	var breed models.Breed

	db, err := gorm.Open(mysql.Open(MysqlDSN), &gorm.Config{})
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        return
    }

	if err := db.First(&breed, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breed)
}

func CreateBreed(w http.ResponseWriter, r *http.Request) {
    var breed models.Breed
    if err := json.NewDecoder(r.Body).Decode(&breed); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	db, err := gorm.Open(mysql.Open(MysqlDSN), &gorm.Config{})
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        return
    }
    
    if err := db.Create(&breed).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(breed)
}

func UpdateBreed(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

	db, err := gorm.Open(mysql.Open(MysqlDSN), &gorm.Config{})
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        return
    }
    
    var breed models.Breed
    if err := db.First(&breed, id).Error; err != nil {
        http.Error(w, "Breed not found", http.StatusNotFound)
        return
    }
    
    if err := json.NewDecoder(r.Body).Decode(&breed); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := db.Save(&breed).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(breed)
}

func DeleteBreed(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

	db, err := gorm.Open(mysql.Open(MysqlDSN), &gorm.Config{})
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        return
    }
    
    var breed models.Breed
    if err := db.First(&breed, id).Error; err != nil {
        http.Error(w, "Breed not found", http.StatusNotFound)
        return
    }
    
    if err := db.Delete(&breed).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}