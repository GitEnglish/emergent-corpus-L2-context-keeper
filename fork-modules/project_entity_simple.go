// DROP-IN: internal/models/project_entity.go
//
// The Reality: It's just a sighting log.
//
// Phenomena are static (~100K possible things humans do with language).
// We observe. We track. We occasionally guess (bands 4-13) when we have to.
//
// No mastery scores. No complex formulas. Just sightings and Ghost State.

package models

import (
	"fmt"
	"time"
)

// =============================================================================
// Sighting - A single observation
// =============================================================================

type Sighting struct {
	Timestamp     time.Time `json:"timestamp"`
	LessonContext string    `json:"lesson_context"`  // What was happening
	SessionID     string    `json:"session_id"`
	Utterance     string    `json:"utterance"`       // What they produced
	ModeledForm   string    `json:"modeled_form"`    // What you said (if relevant)
	
	// Optional: Student self-report (if they chose to interact)
	SelfReport *string `json:"self_report,omitempty"` // "mastered", "uncertain", "new", etc.
}

// =============================================================================
// Ghost State - The visibility toggle
// =============================================================================

type GhostState struct {
	Owned          bool       `json:"owned"`           // false = Ghost (invisible), true = Active (coaching)
	Status         *string    `json:"status"`          // null, "coaching", "monitoring", "archived"
	FirstSighted   time.Time  `json:"first_sighted"`
	LastSighted    time.Time  `json:"last_sighted"`
	SightingCount  int        `json:"sighting_count"`
	Sightings      []Sighting `json:"sightings"`
}

// =============================================================================
// ProjectEntity - The observation container
// =============================================================================

type ProjectEntity struct {
	ID   string `json:"id"`   // Format: {strand}_{item}_{instance}
	Type string `json:"type"` // "linguistic_phenomenon"

	// What it IS (static properties)
	Properties map[string]interface{} `json:"properties"`
	// Expected:
	// - "surface_form": what was produced
	// - "target_form": expected/canonical form (if applicable)
	// - "grammatical_category": e.g., "past-perfect", "modal-verb"
	// - "errant_triad": "M" | "R" | "U" (optional, for errors only)
	// - "pattern_source": "deduced"
	// - "reasoning_chain": how LLM identified it
	// - "validity_confidence": 0.0-1.0 (LLM certainty)

	// Observation state (per student, dynamic)
	Ghost GhostState `json:"ghost"`

	// Heuristic band (4-13) - ONLY used when no other signal exists
	// This is your shitty-but-better-than-nothing guess for prioritization
	HeuristicBand int `json:"heuristic_band,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// NewProjectEntity creates with defaults
// owned: false (Ghost State), no self-report
func NewProjectEntity(id string, properties map[string]interface{}) *ProjectEntity {
	now := time.Now()
	return &ProjectEntity{
		ID:         id,
		Type:       "linguistic_phenomenon",
		Properties: properties,
		Ghost: GhostState{
			Owned:         false, // Start as Ghost
			Status:        nil,
			FirstSighted:  now,
			LastSighted:   now,
			SightingCount: 0,
			Sightings:     []Sighting{},
		},
		HeuristicBand: 0, // 0 = unassigned, use only when needed
		CreatedAt:     now,
		UpdatedAt:     now,
		Version:       1,
	}
}

// AddSighting - Log an observation
func (e *ProjectEntity) AddSighting(sighting Sighting) {
	e.Ghost.Sightings = append(e.Ghost.Sightings, sighting)
	e.Ghost.SightingCount++
	e.Ghost.LastSighted = sighting.Timestamp
	e.Version++
	e.UpdatedAt = time.Now()
}

// PromoteToOwned - YOU decide when to coach this
func (e *ProjectEntity) PromoteToOwned(status string) {
	e.Ghost.Owned = true
	e.Ghost.Status = &status
	e.UpdatedAt = time.Now()
}

// DemoteToGhost - Stop coaching, keep tracking
func (e *ProjectEntity) DemoteToGhost() {
	e.Ghost.Owned = false
	e.Ghost.Status = nil
	e.UpdatedAt = time.Now()
}

// StoragePath - Where to save
func (e *ProjectEntity) StoragePath() string {
	return fmt.Sprintf("accomplishments/%s.json", e.ID)
}

// FullPath - With student directory
func (e *ProjectEntity) FullPath(studentID string) string {
	return fmt.Sprintf("data/students/%s/accomplishments/%s.json", studentID, e.ID)
}

// =============================================================================
// Sanity Bridge (Optional - only when YOU want it)
// =============================================================================

// ToSanityDocument - converts for Sanity import
// Call this when YOU decide to push an owned phenomenon to Sanity
func (e *ProjectEntity) ToSanityDocument() map[string]interface{} {
	return map[string]interface{}{
		"_type": "linguisticPhenomenon",
		"_id":   e.ID,
		"surfaceForm":       e.Properties["surface_form"],
		"targetForm":        e.Properties["target_form"],
		"grammaticalCategory": e.Properties["grammatical_category"],
		"errantTriad":       e.Properties["errant_triad"],
		"ghostState": map[string]interface{}{
			"owned":         e.Ghost.Owned,
			"status":        e.Ghost.Status,
			"sightingCount": e.Ghost.SightingCount,
		},
		"sightings":     e.Ghost.Sightings,
		"heuristicBand": e.HeuristicBand,
		"validityConfidence": e.Properties["validity_confidence"],
		"reasoningChain": e.Properties["reasoning_chain"],
	}
}
