// DROP-IN: internal/llm/gluttonous_prompt.go
//
// The Gluttonous Meta-Prompt - Ingests 600+ items per session
// Zero filtering. Everything goes in. YOU decide later what's worth seeing.
//
// Usage:
//   pm := llm.NewPromptManager(config)
//   pm.RegisterTemplate(llm.GluttonousIngestionTemplate())
//   prompt, _ := pm.BuildPrompt("gluttonous_ingestion", &llm.PromptContext{
//       Transcript: sessionText,
//   })

package llm

import "time"

// GluttonousIngestionTemplate - The vacuum cleaner prompt
func GluttonousIngestionTemplate() *PromptTemplate {
	return &PromptTemplate{
		Name: "gluttonous_ingestion",
		SystemPrompt: `You are a linguistic vacuum cleaner. Extract EVERYTHING.

=== INGESTION LOGIC (v3.0) ===

ZERO-FILTER CAPTURE - Get it ALL:
- Successes (Category 4-13 correct usage)
- Errors (slips, developmental, interference)
- Idiosyncrasies (personal phrases, quirky patterns)
- Micro-structures (single words, collocations)
- Macro-structures (discourse, rhetoric, cohesion)
- Hesitations, self-corrections, approximations
- Everything. No exceptions. Target: 600+ items.

GHOST STATE DEFAULTS (Critical):
{
  "ghost": {
    "owned": false,      // ALWAYS false for new items
    "status": null       // ALWAYS null until human intervenes
  }
}

THE ERRAnT TRIAD (Label every error):
- "M" = Missing: Target not produced when needed
  Example: "I *going" (missing "am")
- "R" = Replacement: Wrong form substituted for target  
  Example: "I am go" ("go" replaces "going")
- "U" = Addition: Unnecessary element inserted
  Example: "I am am going" (extra "am")

MASTERY CATEGORIES 4-13:
- Assign category 4-13 based on curriculum stage
- mastery_score: 1.0 for successful Cat 4-13 usage
- mastery_score: 0.0 for errors (any category)

AUDIT TRAIL (Every item MUST have):
- pattern_source: "deduced"
- reasoning_chain: Your step-by-step deduction
- validity_confidence: 0.0-1.0 (your certainty)

SEMANTIC RELATIONS (Connect phenomena):
- Create weighted relations between items
- Types: metaphorical, taxonomic, collocational, derivational
- Include deduction_reasoning and validity_confidence

OUTPUT: JSON array of ProjectEntity objects

Example structure:
{
  "id": "grammar_past-perfect_001",
  "type": "linguistic_phenomenon", 
  "properties": {
    "surface_form": "I had went",
    "target_form": "I had gone",
    "grammatical_category": "past-perfect",
    "semantic_field": "temporal-sequence",
    "errant_triad": "R",
    "mastery_category": 7,
    "mastery_score": 0.0,
    "pattern_source": "deduced",
    "reasoning_chain": "Student used simple past 'went' instead of past participle 'gone' in past perfect auxiliary construction",
    "validity_confidence": 0.95
  },
  "connections": {
    "outgoing": [{
      "target_id": "morphology_irregular-participle_003",
      "relation_type": "taxonomic",
      "weight": 0.85,
      "deduction_reasoning": "Irregular participle error suggests related morphology pattern",
      "validity_confidence": 0.8
    }]
  },
  "ghost": {
    "owned": false,
    "status": null,
    "first_sighted": "2026-03-22T10:30:00Z",
    "last_sighted": "2026-03-22T10:30:00Z",
    "sighting_count": 1,
    "sightings": [{
      "timestamp": "2026-03-22T10:30:00Z",
      "lesson_context": "Travel narrative: Paris trip",
      "session_id": "sess_alpha_001",
      "utterance": "I had went to Paris before",
      "correction": "I had gone to Paris before"
    }]
  },
  "mastery": {
    "mastery_score": 0.0,
    "confidence_sum": 0.95,
    "weighted_mastery": 0.0,
    "sightings_count": 1
  },
  "created_at": "2026-03-22T10:30:00Z",
  "updated_at": "2026-03-22T10:30:00Z",
  "version": 1
}`,
		UserTemplate: `=== SESSION TRANSCRIPT ===
{{.Transcript}}

=== EXTRACTION TARGETS ===
1. Extract 600+ linguistic phenomena (count them!)
2. Apply ERRAnT triad to ALL errors
3. Set owned: false, status: null for EVERY item
4. Include reasoning_chain and validity_confidence
5. Create semantic relations between related items

Vacuum it all up. Return JSON array:`,
		Variables:    []string{"Transcript"},
		OutputFormat: "json",
		Version:      "v3.0.0-gluttonous",
		Description:  "Zero-filter linguistic extraction with Ghost State",
		Tags:         []string{"ingestion", "linguistic", "errant", "ghost"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Category4To13List - Optional helper to inject into prompt context
// Shows the LLM what categories 4-13 mean in your curriculum
func Category4To13List() string {
	return `
Category 4: Basic sentence patterns, present/past simple
Category 5: Progressive aspects, future forms
Category 6: Perfect aspects, time markers
Category 7: Modal verbs (can, could, may, might)
Category 8: Conditionals (zero, first, second)
Category 9: Passive voice, reported speech
Category 10: Relative clauses, noun phrases
Category 11: Complex sentences, subordination
Category 12: Academic register, hedging
Category 13: Nuanced modality, metadiscourse
`
}
