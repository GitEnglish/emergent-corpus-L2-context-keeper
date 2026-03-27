// DROP-IN FILE: Add to internal/config/config.go
// 
// Add these fields to the Config struct, then update Load() to populate them

package config

// GluttonousConfig holds the boutique tutoring fork configuration
// Add this as a field to the main Config struct:
// 
// type Config struct {
//     ... existing fields ...
//     Gluttonous GluttonousConfig
// }
//
type GluttonousConfig struct {
	// Student isolation: each student has separate context.json file
	// ContextFilePath is parameterized via STUDENT_ID env var
	ContextFilePath string // e.g., "./data/students/${STUDENT_ID}/context.json"
	
	// Ingestion settings
	ZeroFilterMode       bool    // Extract EVERY phenomenon, no filtering
	TargetItemsPerSession int    // Target: 600+ items per session
	DefaultOwned         bool    // Default ghost.owned value (false)
	
	// ERRAnT Triad classification
	ErrantEnabled        bool     // Enable M/R/U error classification
	ErrantTriadValues    []string // ["M", "R", "U"]
	
	// Mastery tracking (Categories 4-13)
	MasteryCategories        []int   // [4, 5, 6, 7, 8, 9, 10, 11, 12, 13]
	DefaultMasteryScore      float64 // 1.0 for Cat 4-13 successes
	
	// Ghost State thresholds
	ConfidenceThreshold    float64 // Hide if validity_confidence < 0.7
	VisibilityThreshold    float64 // Show if weighted_mastery > 0.5
	
	// Weighted Mastery Formula
	// W_final = Σ(Mastery · Confidence) / N_sightings
	WeightingFormula string
}

// LoadGluttonousConfig creates default gluttonous configuration
// Call this in Load() function:
//
// config.Gluttonous = LoadGluttonousConfig()
//
func LoadGluttonousConfig() GluttonousConfig {
	return GluttonousConfig{
		// Student isolation via parameterized context file path
		// Set STUDENT_ID env var to switch between students
		ContextFilePath: getEnv("STUDENT_CONTEXT_PATH", "./data/students/${STUDENT_ID}/context.json"),
		
		// Ingestion
		ZeroFilterMode:        getEnvAsBool("ZERO_FILTER_MODE", true),
		TargetItemsPerSession: getEnvAsInt("TARGET_ITEMS_PER_SESSION", 600),
		DefaultOwned:          false, // Ghost State default
		
		// ERRAnT
		ErrantEnabled:     true,
		ErrantTriadValues: []string{"M", "R", "U"},
		
		// Mastery
		MasteryCategories:   []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		DefaultMasteryScore: 1.0,
		
		// Thresholds
		ConfidenceThreshold: 0.7,  // Noise filter
		VisibilityThreshold: 0.5,  // W_final threshold
		
		// Formula
		WeightingFormula: "weighted_average", // W_final = Σ(M·C)/N
	}
}

// GetStudentContextPath returns the context file path for a given student
// Usage: config.Gluttonous.GetStudentContextPath("student_alpha")
func (g *GluttonousConfig) GetStudentContextPath(studentID string) string {
	// Simple variable substitution
	// In production, use proper template replacement
	return replaceVariable(g.ContextFilePath, "${STUDENT_ID}", studentID)
}

// Helper for variable substitution
func replaceVariable(template, variable, value string) string {
	// Simple string replacement
	// Consider using os.Expand or template package for complex cases
	result := template
	for {
		newResult := replaceOnce(result, variable, value)
		if newResult == result {
			break
		}
		result = newResult
	}
	return result
}

func replaceOnce(s, old, new string) string {
	idx := 0
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			idx = i
			return s[:idx] + new + s[idx+len(old):]
		}
	}
	return s
}

// NOTE: Add to your Load() function:
//
// func Load() *Config {
//     // ... existing code ...
//     
//     config := &Config{
//         // ... existing fields ...
//     }
//     
//     // Add this line:
//     config.Gluttonous = LoadGluttonousConfig()
//     
//     return config
// }
