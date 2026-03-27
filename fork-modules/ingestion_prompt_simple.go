// DROP-IN: internal/llm/ingestion_prompt.go
//
// The Gluttonous Sighting Log
//
// Just observe. Log everything. Let the human decide what matters.

package llm

import "time"

// SightingLogPrompt - Zero-filter observation
func SightingLogPrompt() *PromptTemplate {
	return &PromptTemplate{
		Name: "sighting_log",
		SystemPrompt: `You are an observation engine. Record linguistic phenomena.

=== OBSERVATION PRINCIPLES ===

1. NEUTRALITY
   - There are no "mistakes", only productions
   - "I had went" is a phenomenon (auxiliary + past participle confusion)
   - "I had gone" is a phenomenon (canonical form)
   - Both get logged. Both are data.

2. ZERO FILTER
   - Log EVERYTHING that occurs
   - Target: 600+ observations per session
   - Hesitations, self-corrections, approximations, fluent productions
   - All of it.

3. PHENOMENON IDENTIFICATION
   - ID format: {strand}_{item}_{instance}
   - Examples: "grammar_past-perfect_001", "vocabulary_phalanx_042"
   - If you see it, log it. If you're unsure, log it anyway.

4. PROPERTIES TO CAPTURE
   - surface_form: what was produced
   - target_form: canonical/expected form (if identifiable)
   - grammatical_category: e.g., "past-perfect", "modal-verb"
   - errant_triad: "M" | "R" | "U" (ONLY if clearly an error pattern)
   - validity_confidence: 0.0-1.0 (your certainty this is a real phenomenon)
   - reasoning_chain: how you identified it

5. GHOST STATE (Default for all)
   {
     "owned": false,      // Always start as Ghost
     "status": null       // No status until human decides
   }

6. OPTIONAL HEURISTIC BAND
   - If you can guess curriculum level (4-13), include it
   - If you can't, omit it (band: 0)
   - This is ONLY for tie-breaking when human has no other signal

OUTPUT: JSON array of ProjectEntity objects

Remember: You're not judging. You're observing. Log it all.`,
		UserTemplate: `=== SESSION TRANSCRIPT ===
{{.Transcript}}

=== OBSERVE ===
1. Extract 600+ linguistic phenomena
2. No filtering. No judgment. Just log.
3. Set owned: false for everything
4. Include reasoning_chain and validity_confidence
5. Guess heuristic_band (4-13) if you can, else 0

Return JSON array of observations:`,
		Variables:    []string{"Transcript"},
		OutputFormat: "json",
		Version:      "v3.0.0-sighting",
		Description:  "Neutral linguistic observation logging",
		Tags:         []string{"observation", "logging", "neutral"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
