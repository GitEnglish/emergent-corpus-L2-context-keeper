# Final Architecture Clarifications

## 1. Component Scope

**Decision:** Gluttonous Ingestion Engine = ingestion pipeline only (transcript → sightings)

**Why:** You said "that injection is ONLY 1 small part" - the original spec described a wholesale rewrite, but you only wanted the ingestion module added to existing Context-Keeper.

**How:** Create `internal/ingestion/batch_processor.go` that takes transcripts and outputs `ProjectEntity` sightings. Context-Keeper's session store, API handlers, and existing functionality remain untouched.

---

## 2. Student Isolation

**Decision:** Directory-based: `data/students/{studentID}/`

**Why:** You said "a different context.json file for every single student" - file-based isolation, not environment variables or complex multi-tenancy.

**How:** `StudentResolver` generates paths:
- `data/students/{id}/sightings/{phenomenonId}.json`
- `data/students/{id}/corpus.duckdb`
- Each student operates in complete isolation. No cross-contamination possible.

---

## 3. Sanity Relationship

**Decision:** Modular, not coupled. Context-Keeper produces JSON; Sanity can read it optionally.

**Why:** "I want everything to be modular" + "I don't want to marry them" + "They just hook up when you feel like it."

**How:** 
- Context-Keeper saves sightings to local JSON files
- Optional `sanity_bridge.go` script YOU run when YOU want to push to Sanity
- Context-Keeper never queries Sanity during normal operation
- No webhooks, no real-time sync, no automatic pushes

---

## 4. Phenomena Nature

**Decision:** Static taxonomy (~100K possible), start with small seed set, grow organically via human approval.

**Why:** "There's just a hundred thousand of them. They are static. Because you can only say things in so many possible ways." + "I should give it like 500 things and just let it wild on transcripts."

**How:**
- Seed file: 500 initial phenomenon definitions
- LLM ingestion matches transcripts to known phenomena OR proposes new ones
- New proposals saved as Sanity drafts (taxonomyProposal documents)
- YOU review in Sanity Studio: approve (adds to taxonomy) or reject (maps to existing)
- Only approved phenomena get official IDs

---

## 5. Data Model

**Decision:** Sighting log only. No mastery scores. No formulas.

**Why:** "It's just a sighting log, exactly" + "There is no mastery level, exactly. It's observed versus not observed."

**How:**
```go
type ProjectEntity struct {
    ID          string                 // static phenomenon ID
    Properties  map[string]interface{} // surface_form, target_form, etc.
    Ghost       GhostState             // owned flag + sightings array
    CreatedAt   time.Time
}

type GhostState struct {
    Owned         bool       // false = tracking only, true = coaching
    Status        *string    // null, "coaching", "mastered", "archived"
    SightingCount int
    Sightings     []Sighting // every observation
}
```

---

## 6. Ghost State

**Decision:** Simple boolean visibility. You control promotion.

**Why:** "Ghost State... invisible majority" + decoupling "data retention (the vacuum) from user visibility (the toolkit)."

**How:**
- Default: `owned: false` - phenomenon is tracked but invisible in coaching UI
- You call `entity.PromoteToOwned("coaching")` when YOU decide to address it
- `entity.DemoteToGhost()` when you stop actively coaching but keep tracking
- Sightings accumulate regardless of owned status ("blind update")

---

## 7. Self-Reports

**Decision:** Optional free-form string. Student describes their own state.

**Why:** "They may or may not have chosen to interact with it and say, like, I mastered this, or I have no idea."

**How:**
```go
type Sighting struct {
    Timestamp     time.Time
    LessonContext string
    Utterance     string
    SelfReport    *string  // "mastered", "uncertain", "new", or null
}
```

---

## 8. Processing Model

**Decision:** Batch processing, hourly acceptable.

**Why:** "It could do one section every hour and that's fine. It doesn't need to be instant."

**How:**
- Cron job or systemd timer runs `cmd/worker/main.go` every hour
- Worker scans for unprocessed transcripts
- Runs ingestion pipeline
- Updates DuckDB analytics
- No real-time requirements, no webhooks, no streaming

---

## 9. Analytics

**Decision:** DuckDB reads Context-Keeper JSON directly. One query interface.

**Why:** "DuckDB can handle that for me and just be one point of contact, one point of complete confusion instead of two."

**How:**
```sql
-- DuckDB reads JSON files directly
CREATE VIEW all_sightings AS
SELECT * FROM read_json_auto('data/students/*/sightings/*.json');

-- Lexical spread: distinct contexts per token
SELECT token, COUNT(DISTINCT lesson_context) as spread 
FROM all_sightings 
GROUP BY token;
```

---

## 10. COCA 5000

**Decision:** Your data seeds taxonomy. Not used during ingestion.

**Why:** "The COCA 5000 as much as I would like it to be more exotic" - frequency data informs taxonomy, doesn't drive classification.

**How:**
- COCA ranks stored in phenomenon properties (`coca_rank`, `coca_freq_pm`)
- Analytics can query: "Show me high-frequency phenomena this student hasn't mastered"
- Ingestion doesn't care about frequency - it logs what it sees

---

## 11. Schema Deployment

**Decision:** No automatic deployments to live Sanity. You control all schema changes.

**Why:** "It's contaminating shit" + "It needs to be the way that I left it, exactly the way I left it."

**How:**
- All schemas prepared as local files only (`fork-modules/`)
- You review and deploy manually when ready
- MCP tools blocked for deployments
- Your existing `errorPattern` (91k Master List), `grammarError`, `vocabularyItem` etc. preserved exactly as-is
