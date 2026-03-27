# Gluttonous Ingestion Engine - Implementation Plan

## Current Status: 30% Complete

### ✅ DONE
- [x] ProjectEntity model (simplified, sighting-log based)
- [x] Ingestion prompt template
- [x] Student directory resolver
- [x] Sanity bridge (optional)

### ❌ NOT DONE (Blocking Shipping)

## Phase 1: Core Ingestion Pipeline (Week 1)

### 1.1 Batch Ingestion Engine
**File:** `internal/ingestion/batch_processor.go`

```go
type BatchProcessor struct {
    promptManager *llm.PromptManager
    resolver      *utils.StudentResolver
    outputDir     string
}

func (bp *BatchProcessor) ProcessTranscript(studentID string, transcript string) ([]*models.ProjectEntity, error) {
    // 1. Load LLM prompt
    // 2. Call LLM with transcript
    // 3. Parse JSON response into ProjectEntity slice
    // 4. Save each entity to student directory
    // 5. Return entities for further processing
}
```

**Deliverable:** Function that takes (studentID, transcript) → writes sightings to disk

### 1.2 LLM Pipeline Integration
**File:** `internal/llm/adapter.go` extension

- Wire `SightingLogPrompt()` into existing LLM client
- Handle response parsing & error recovery
- Retry logic for failed extractions

**Deliverable:** `ExtractPhenomena(transcript string) ([]ProjectEntity, error)`

### 1.3 Taxonomy Proposal System
**File:** `internal/ingestion/taxonomy_matcher.go`

```go
type TaxonomyMatcher struct {
    sanityClient *sanity.Client  // optional
    localCache   map[string]bool // known phenomenon IDs
}

func (tm *TaxonomyMatcher) MatchOrPropose(entity *ProjectEntity) (MatchType, error) {
    // 1. Check if entity.ID exists in Sanity taxonomy
    // 2. If yes: return MATCH_EXISTING
    // 3. If no: create proposal, return PROPOSAL_NEW
}
```

**Deliverable:** System that separates known phenomena from proposals

---

## Phase 2: Corpus Linguistics Engine (Week 2)

### 2.1 Tokenization Module
**File:** `internal/corpus/tokenizer.go`

```go
type Tokenizer struct {
    // Tokenize transcripts into lexical items
    // Extract n-grams
    // Calculate basic frequencies
}

func (t *Tokenizer) Tokenize(transcript string) []Token {
    // Return: [{word: "phalanx", pos: "noun", index: 42}, ...]
}
```

**Deliverable:** Text → structured tokens

### 2.2 Frequency Counter
**File:** `internal/corpus/frequency.go`

```go
type FrequencyCounter struct {
    duckDB *sql.DB
}

func (fc *FrequencyCounter) UpdateFrequency(studentID string, tokens []Token) error {
    // Increment per-student word frequencies
    // Update type-token ratios
}

func (fc *FrequencyCounter) GetLexicalSpread(studentID string) (LexicalSpreadReport, error) {
    // DISTINCT contexts per token
    // Dispersion metrics
}
```

**Deliverable:** Word frequency & dispersion calculation

### 2.3 DuckDB Schema
**File:** `migrations/001_initial.sql`

```sql
CREATE TABLE token_frequencies (
    student_id TEXT,
    token TEXT,
    count INTEGER,
    contexts INTEGER,  -- DISTINCT lesson_contexts
    first_seen TIMESTAMP,
    last_seen TIMESTAMP
);

CREATE TABLE phenomena_sightings (
    student_id TEXT,
    phenomenon_id TEXT,
    sighting_count INTEGER,
    contexts TEXT[],  -- array of lesson_contexts
    FOREIGN KEY (phenomenon_id) REFERENCES sanity_taxonomy(id)
);
```

**Deliverable:** Database schema for analytics

---

## Phase 3: Batch Scheduler (Week 3)

### 3.1 Hourly Worker
**File:** `cmd/worker/main.go`

```go
func main() {
    // 1. Find unprocessed transcripts
    // 2. For each: Run ingestion
    // 3. Update DuckDB analytics
    // 4. Log completion
}
```

**Deliverable:** Standalone worker process

### 3.2 Cron Integration
**File:** `docker-compose.yml` or systemd service

```yaml
worker:
  build: .
  command: ["/app/worker", "--interval=1h"]
  volumes:
    - ./data:/data
```

**Deliverable:** Automated batch processing

---

## Phase 4: Sanity Integration (Week 3-4)

### 4.1 Sanity Schema
**File:** `sanity/schemaTypes/linguisticPhenomenon.ts`

See `sanity_bridge.go` for schema definition.

**Deliverable:** Deployed Sanity schema

### 4.2 Taxonomy Sync
**File:** `internal/taxonomy/sync.go`

```go
func (s *TaxonomySync) ProposeToSanity(entity *ProjectEntity) error {
    // POST to Sanity as draft
    // Return proposal ID
}

func (s *TaxonomySync) LoadApprovedTaxonomy() (map[string]bool, error) {
    // Fetch all approved phenomena from Sanity
    // Build lookup map for matcher
}
```

**Deliverable:** Bidirectional Sanity sync

---

## Phase 5: Testing & Integration (Week 4)

### 5.1 End-to-End Test
**File:** `tests/e2e/ingestion_test.go`

```go
func TestFullPipeline(t *testing.T) {
    // 1. Submit transcript
    // 2. Wait for batch processing
    // 3. Verify sightings written
    // 4. Verify DuckDB updated
    // 5. Verify taxonomy proposals created
}
```

**Deliverable:** Working E2E test

### 5.2 Load Test
- Process 10 sessions × 2,000 words
- Verify performance: < 5 min per batch

**Deliverable:** Performance validation

---

## Dependencies

### External
- [ ] DuckDB Go driver: `github.com/marcboeker/go-duckdb`
- [ ] Sanity Go client: `github.com/sanity-io/client-go`
- [ ] Cron library: `github.com/robfig/cron/v3`

### Sanity
- [ ] Project setup
- [ ] API token
- [ ] Schema deployed

---

## Ship Checklist

- [ ] `ProcessTranscript()` works end-to-end
- [ ] Sightings saved to `data/students/{id}/sightings/`
- [ ] DuckDB analytics updating
- [ ] Batch worker running hourly
- [ ] Taxonomy proposals visible in Sanity
- [ ] You can query: "Show me lexical spread for student_alpha"

**Current blocking issues:**
1. No actual ingestion implementation
2. No DuckDB integration
3. No batch processor
4. No Sanity schema deployed

**Estimated time to ship: 3-4 weeks** (with focused work)
