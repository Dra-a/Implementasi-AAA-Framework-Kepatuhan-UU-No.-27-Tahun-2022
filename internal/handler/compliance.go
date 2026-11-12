package handler

import (
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"net/http"
	"time"
)

func GetMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	profile, err := lmodel.GetUserProfile(user.ID)

	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func UpdateMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	ar req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := lmodel.UpdateUserProfile(user.ID, req.Username, req.Email); err != nil {
		http.Error(w, "Failed to update data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data updated successfully",
		"username": req.Username,
	})
}

func GetMyAuditLog(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	logs, err := lmodel.GetAuditLogsByUserID(user.ID)

	if err != nil {
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user.Username,
		"total": len(logs),
		"logs":  logs,
	})
}

func ExportMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	profile, err := lmodel.GetUserProfile(user.ID)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}
	projects, err := lmodel.GetProjectsByUserID(user.ID)
	if err != nil {
		http.Error(w, "Failed to get projects", http.StatusInternalServerError)
		return
	}

	logs, err := lmodel.GetAuditLogsByUserID(user.ID)
	if err != nil {
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	export := map[string]interface{}{
		"exported_at": time.Now(),
		"profile":     profile,
		"projects":    projects,
		"audit_logs":  logs,
	}

	w.Header().Set("Content-Disposition", "attachment; filename=my_data_export.json")
	json.NewEncoder(w).Encode(export)
}
func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	if err := lmodel.DeleteUserAndData(user.ID); err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your account and all data have been permanently deleted",
	})
}
