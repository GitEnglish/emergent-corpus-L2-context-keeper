// DROP-IN: internal/models/project_entity.go
// 
// The Extensibility Engine - Fluid schema via map[string]interface{}
// No recompilation needed for v3.0+ schema evolution
// 
// This is YOUR data structure. Sanity can read it if it wants. Or not.

package models

import (
	"fmt"
	"time"
)

// =============================================================================
// ERRAnT Triad - Error Classification System
// =============================================================================

type ErrantTriad string

const (
	ErrantMissing     ErrantTriad = "M" // Missing - should have been used but wasn't
	ErrantReplacement ErrantTriad = "R" // Replacement - wrong form used instead
	ErrantAddition    ErrantTriad = "U" // Addition - unnecessary element added
)

// =============================================================================
// Sighting - Every time we see a phenomenon
// =============================================================================

type Sighting struct {
	Timestamp     time.Time `json:"timestamp"`
	LessonContext string    `json:"lesson_context"`  // "Travel narrative exercise"
	SessionID     string    `json:"session_id"`
	Utterance     string    `json:"utterance"`       // What they said
	Correction    string    `json:"correction"`      // Target form
}

// =============================================================================
// Semantic Relations - LLM-deduced connections
// =============================================================================

type SemanticRelation struct {
	TargetID           string  `json:"target_id"`
	RelationType       string  `json:"relation_type"`       // metaphorical, taxonomic, collocational
	Weight             float64 `json:"weight"`
	DeductionReasoning string  `json:"deduction_reasoning"`
	ValidityConfidence float64 `json:"validity_confidence"` // 0.0-1.0
}

type EntityConnection struct {
	Incoming []SemanticRelation `json:"incoming,omitempty"`
	Outgoing []SemanticRelation `json:"outgoing,omitempty"`
}

// =============================================================================
// Ghost State - The Magic Sauce
// 
// owned: false = "Ghost State" (invisible, tracked, calculated)
// owned: true  = "Active Toolkit" (appears in your Sanity if you want it)
// =============================================================================

type GhostState struct {
	Owned         bool       `json:"owned"`           // false = Ghost, true = Active
	Status        *string    `json:"status"`          // null, "coaching", "mastered"
	FirstSighted  time.Time  `json:"first_sighted"`
	LastSighted   time.Time  `json:"last_sighted"`
	SightingCount int        `json:"sighting_count"`
	Sightings     []Sighting `json:"sightings"`
}

// =============================================================================
// Mastery Metrics - Weighted calculation
// W_final = Σ(Mastery · Confidence) / N_sightings
// =============================================================================

type MasteryMetrics struct {
	MasteryScore       float64 `json:"mastery_score"`       // 1.0 for Cat 4-13 successes
	ConfidenceSum      float64 `json:"confidence_sum"`
	WeightedMastery    float64 `json:"weighted_mastery"`    // The calculated W_final
	SightingsCount     int     `json:"sightings_count"`
}

// =============================================================================
// ProjectEntity - THE DUMB PIPE
// 
// ID format: {strand}_{item}_{number} 
// Example: "grammar_past-perfect_001", "vocabulary_phalanx_042"
//
// Properties holds ANYTHING. It's your fluid v3.0 schema:
//   - "surface_form": what they said
//   - "target_form": correct version
//   - "grammatical_category": "past-perfect", "modal-verb", etc.
//   - "semantic_field": "temporal", "spatial", "causative"
//   - "errant_triad": "M" | "R" | "U"
//   - "mastery_category": 4-13
//   - "pattern_source": "deduced"
//   - "reasoning_chain": "Student used past participle..."
//   - "validity_confidence": 0.95
//
// Sanity can read this. Or not. Your call.
// =============================================================================

type ProjectEntity struct {
	ID   string `json:"id"`
	Type string `json:"type"` // e.g., "linguistic_phenomenon"

	// THE FORK POINT: Fluid schema
	Properties map[string]interface{} `json:"properties"`

	// Connections between phenomena
	Connections EntityConnection `json:"connections"`

	// Ghost State (decoupled retention from visibility)
	Ghost GhostState `json:"ghost"`

	// Weighted mastery tracking
	Mastery MasteryMetrics `json:"mastery"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// NewProjectEntity creates with Ghost State defaults
// owned: false (invisible), status: null
func NewProjectEntity(id string, entityType string, properties map[string]interface{}) *ProjectEntity {
	now := time.Now()
	return &ProjectEntity{
		ID:         id,
		Type:       entityType,
		Properties: properties,
		Ghost: GhostState{
			Owned:         false, // Start as Ghost
			Status:        nil,
			FirstSighted:  now,
			LastSighted:   now,
			SightingCount: 0,
			Sightings:     []Sighting{},
		},
		Mastery: MasteryMetrics{
			MasteryScore:    0.0,
			ConfidenceSum:   0.0,
			WeightedMastery: 0.0,
			SightingsCount:  0,
		},
		Connections: EntityConnection{
			Incoming: []SemanticRelation{},
			Outgoing: []SemanticRelation{},
		},
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}
}

// AddSighting - "Blind Update" 
// Tracks EVERY occurrence regardless of ownership
func (e *ProjectEntity) AddSighting(sighting Sighting) {
	e.Ghost.Sightings = append(e.Ghost.Sightings, sighting)
	e.Ghost.SightingCount++
	e.Ghost.LastSighted = sighting.Timestamp
	e.Mastery.SightingsCount++
	e.Version++
	e.UpdatedAt = time.Now()
}

// CalculateWeightedMastery - W_final = Σ(Mastery · Confidence) / N
func (e *ProjectEntity) CalculateWeightedMastery() float64 {
	if e.Mastery.SightingsCount == 0 {
		e.Mastery.WeightedMastery = 0.0
		return 0.0
	}

	confidence := 0.7 // default
	if c, ok := e.Properties["validity_confidence"].(float64); ok {
		confidence = c
	}

	numerator := e.Mastery.MasteryScore * confidence
	e.Mastery.WeightedMastery = numerator / float64(e.Mastery.SightingsCount)
	return e.Mastery.WeightedMastery
}

// IsVisible - Should this appear in Active Toolkit?
// Rules: owned=true AND validity_confidence >= 0.7 AND weighted_mastery > 0.5
func (e *ProjectEntity) IsVisible() bool {
	var confidence float64 = 1.0
	if c, ok := e.Properties["validity_confidence"].(float64); ok {
		confidence = c
	}
	if confidence < 0.7 {
		return false // Noise filter
	}
	if !e.Ghost.Owned {
		return false // Still a Ghost
	}
	if e.Mastery.WeightedMastery <= 0.5 {
		return false // Not mature enough
	}
	return true
}

// PromoteToOwned - YOU decide when a Ghost becomes Active
// This is your human-in-the-loop intervention point
func (e *ProjectEntity) PromoteToOwned(status string) {
	e.Ghost.Owned = true
	e.Ghost.Status = &status
	e.UpdatedAt = time.Now()
}

// StoragePath - Where to save this (relative to student directory)
func (e *ProjectEntity) StoragePath() string {
	return fmt.Sprintf("accomplishments/%s.json", e.ID)
}

// FullPath - Complete path given student ID
func (e *ProjectEntity) FullPath(studentID string) string {
	return fmt.Sprintf("data/students/%s/accomplishments/%s.json", studentID, e.ID)
}

// ToSanityDocument - Converts to Sanity-compatible structure
// Call this when YOU decide to push to Sanity
func (e *ProjectEntity) ToSanityDocument() map[string]interface{} {
	return map[string]interface{}{
		"_type": "linguisticPhenomenon",
		"_id":   e.ID,
		"surfaceForm":       e.Properties["surface_form"],
		"targetForm":        e.Properties["target_form"],
		"grammaticalCategory": e.Properties["grammatical_category"],
		"errantTriad":       e.Properties["errant_triad"],
		"masteryCategory":   e.Properties["mastery_category"],
		"ghostState": map[string]interface{}{
			"owned":         e.Ghost.Owned,
			"status":        e.Ghost.Status,
			"sightingCount": e.Ghost.SightingCount,
		},
		"sightings": e.Ghost.Sightings,
		"masteryScore": e.Mastery.MasteryScore,
		"weightedMastery": e.Mastery.WeightedMastery,
		"validityConfidence": e.Properties["validity_confidence"],
		"reasoningChain": e.Properties["reasoning_chain"],
	}
}
