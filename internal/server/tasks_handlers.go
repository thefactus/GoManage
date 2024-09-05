package server

import (
	"go-tasker/types"
	"go-tasker/utils"
	"net/http"
)

func (s *Server) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("listID")
	tasks, err := s.db.GetTasks(taskID)
	if err != nil {
		http.Error(w, "Error getting tasks", http.StatusInternalServerError)
		return
	}

	response := utils.PrepareJSONWithMessage("Tasks retrieved successfully", tasks)

	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Server) PostTasksHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("listID")

	var createTaskPayload types.CreateTaskPayload
	if err := utils.ParseAndValidateJSON(w, r, &createTaskPayload); err != nil {
		return
	}

	task, err := s.db.CreateTask(taskID, createTaskPayload)
	if err != nil {
		utils.WriteInternalServerError(w, err)
		return
	}

	response := utils.PrepareJSONWithMessage("Task created successfully", task)

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (s *Server) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")

	var updateTaskPayload types.UpdateTaskPayload
	if err := utils.ParseAndValidateJSON(w, r, &updateTaskPayload); err != nil {
		return
	}

	task, err := s.db.UpdateTask(taskID, updateTaskPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := utils.PrepareJSONWithMessage("Task updated successfully", task)

	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")

	err := s.db.DeleteTask(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := utils.PrepareJSONWithMessage("Task deleted successfully", nil)

	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Server) PatchTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")

	var updateTaskDonePayload types.UpdateTaskDonePayload
	updateTaskDonePayload.Done = true

	task, err := s.db.UpdateTaskDone(taskID, updateTaskDonePayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := utils.PrepareJSONWithMessage("Task marked as done successfully", task)

	utils.WriteJSON(w, http.StatusOK, response)
}

func (s *Server) PatchTaskUndoneHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")

	var updateTaskDonePayload types.UpdateTaskDonePayload
	updateTaskDonePayload.Done = false

	task, err := s.db.UpdateTaskDone(taskID, updateTaskDonePayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := utils.PrepareJSONWithMessage("Task marked as undone successfully", task)

	utils.WriteJSON(w, http.StatusOK, response)
}
