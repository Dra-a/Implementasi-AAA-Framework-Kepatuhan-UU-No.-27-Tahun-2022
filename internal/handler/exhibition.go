package handler

import (
	"aaa_bekend/internal/db"
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"net/http"
	"strconv"

)

type SubmitProjectReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GradeProjectReq struct {
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

// POST /api/v1/projects

func SubmitProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	var req SubmitProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simpan proyek ke database menggunakan ID user yang sedang login
	projectID, err := db.CreateProject(r.Context(), req.Title, req.Description, userCtx.ID)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"project_id": projectID})
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := db.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []lmodel.Project{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&projects)
}
func GetAllProjectGrades(w http.ResponseWriter, r *http.Request) {
	grades, err := db.GetAllProjectGrades(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch project grades", http.StatusInternalServerError)
		return
	}

	if grades == nil {
		grades = []lmodel.ProjectGrade{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&grades)
}
func GradeProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	idStr := r.PathValue("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req GradeProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.CreateGrade(r.Context(), projectID, userCtx.ID, req.Score, req.Feedback)
	if err != nil {
		http.Error(w, "Failed to submit grade", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Grade Submitted"})
}

func GetProjectGrades(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())
	
	idStr := r.PathValue("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Ambil data proyek untuk memeriksa kepemilikan
	project, err := db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Jika role-nya student, pastikan dia hanya melihat proyek miliknya sendiri
	if userCtx.Role == "student" {
		if project.StudentID != userCtx.ID {
			http.Error(w, "Forbidden: you can only view grades for your own projects", http.StatusForbidden)
			return
		}
	}

	grades, err := db.GetGradesForProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
		return
	}

	if grades == nil {
		grades = []lmodel.Grade{}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}
