// DROP-IN: internal/utils/student_resolver.go
//
// Simple path resolver for student isolation
// No complex multi-tenancy. Just directories.
//
// Usage:
//   resolver := utils.NewStudentResolver(cfg.StoragePath)
//   studentPath := resolver.GetStudentPath("student_alpha")
//   // Returns: {StoragePath}/students/student_alpha/

package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// StudentResolver - Simple path generator for student directories
type StudentResolver struct {
	BasePath string // e.g., ~/Library/Application Support/context-keeper
}

// NewStudentResolver creates resolver with base storage path
func NewStudentResolver(basePath string) *StudentResolver {
	return &StudentResolver{BasePath: basePath}
}

// GetStudentPath returns full path to student directory
func (r *StudentResolver) GetStudentPath(studentID string) string {
	return filepath.Join(r.BasePath, "students", studentID)
}

// GetSessionsPath returns path to student's sessions
func (r *StudentResolver) GetSessionsPath(studentID string) string {
	return filepath.Join(r.GetStudentPath(studentID), "sessions")
}

// GetHistoriesPath returns path to student's histories
func (r *StudentResolver) GetHistoriesPath(studentID string) string {
	return filepath.Join(r.GetStudentPath(studentID), "histories")
}

// GetAccomplishmentsPath returns path to student's Gluttonous Ingestion data
func (r *StudentResolver) GetAccomplishmentsPath(studentID string) string {
	return filepath.Join(r.GetStudentPath(studentID), "accomplishments")
}

// GetEntityPath returns full path for a specific entity
func (r *StudentResolver) GetEntityPath(studentID string, entityID string) string {
	return filepath.Join(r.GetAccomplishmentsPath(studentID), entityID+".json")
}

// InitializeStudent creates all directories for a new student
func (r *StudentResolver) InitializeStudent(studentID string) error {
	dirs := []string{
		r.GetSessionsPath(studentID),
		r.GetHistoriesPath(studentID),
		r.GetAccomplishmentsPath(studentID),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}
	return nil
}

// ListStudents returns all student IDs in the system
func (r *StudentResolver) ListStudents() ([]string, error) {
	studentsDir := filepath.Join(r.BasePath, "students")
	
	entries, err := os.ReadDir(studentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	
	var students []string
	for _, entry := range entries {
		if entry.IsDir() {
			students = append(students, entry.Name())
		}
	}
	return students, nil
}

// StudentExists checks if student has initialized directories
func (r *StudentResolver) StudentExists(studentID string) bool {
	_, err := os.Stat(r.GetStudentPath(studentID))
	return !os.IsNotExist(err)
}
