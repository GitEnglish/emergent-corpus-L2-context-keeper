# Gluttonous Ingestion Engine - Implementation Plan

## Sanity Project: git English Hub (rzug0rgk)
## Dataset: production

---

## PHASE 1: Sanity Schema Deployment (Using MCP)

### 1.1 Deploy `linguisticPhenomenon` Document Type

**Action:** Use `deploy_schema` MCP tool

**Schema Definition:**
```typescript
// Document: linguisticPhenomenon
// Purpose: Static taxonomy of ~100K possible phenomena
{
  name: 'linguisticPhenomenon',
  title: 'Linguistic Phenomenon',
  type: 'document',
  icon: DocumentIcon,
  fields: [
    {
      name: 'phenomenonId',
      title: 'Phenomenon ID',
      type: 'string',
      validation: (rule) => rule.required().regex(/^[a-z]+_[a-z-]+_\d+$/),
      description: 'Format: {strand}_{item}_{number} e.g., grammar_past-perfect_001'
    },
    {
      name: 'surfaceForm',
      title: 'Surface Form',
      type: 'string',
      description: 'What the student produces'
    },
    {
      name: 'targetForm',
      title: 'Target Form',
      type: 'string',
      description: 'Canonical/expected form'
    },
    {
      name: 'grammaticalCategory',
      title: 'Grammatical Category',
      type: 'string',
      options: {
        list: [
          {title: 'Past Perfect', value: 'past-perfect'},
          {title: 'Present Perfect', value: 'present-perfect'},
          {title: 'Modal Verb', value: 'modal-verb'},
          {title: 'Conditional', value: 'conditional'},
          {title: 'Passive Voice', value: 'passive-voice'},
          {title: 'Relative Clause', value: 'relative-clause'},
          {title: 'Lexical Item', value: 'lexical-item'},
          {title: 'Discourse Marker', value: 'discourse-marker'}
        ]
      }
    },
    {
      name: 'errantTriad',
      title: 'ERRAnT Classification',
      type: 'string',
      options: {
        list: [
          {title: 'Missing (M)', value: 'M'},
          {title: 'Replacement (R)', value: 'R'},
          {title: 'Addition (U)', value: 'U'}
        ],
        layout: 'radio'
      }
    },
    {
      name: 'heuristicBand',
      title: 'Heuristic Band (4-13)',
      type: 'number',
      description: 'Optional prioritization band when no other signal exists'
    },
    {
      name: 'patternSource',
      title: 'Pattern Source',
      type: 'string',
      initialValue: 'deduced',
      readOnly: true
    },
    {
      name: 'reasoningChain',
      title: 'Reasoning Chain',
      type: 'text',
      description: 'How the LLM identified this phenomenon'
    },
    {
      name: 'semanticRelations',
      title: 'Semantic Relations',
      type: 'array',
      of: [{
        type: 'object',
        fields: [
          {name: 'target', type: 'reference', to: [{type: 'linguisticPhenomenon'}]},
          {name: 'relationType', type: 'string', options: {list: ['metaphorical', 'taxonomic', 'collocational', 'derivational']}},
          {name: 'weight', type: 'number'},
          {name: 'validityConfidence', type: 'number'}
        ]
      }]
    }
  ]
}
```

**MCP Command:** `deploy_schema` with above declaration

---

### 1.2 Deploy `taxonomyProposal` Document Type

**Purpose:** Draft queue for new phenomena awaiting human approval

```typescript
{
  name: 'taxonomyProposal',
  title: 'Taxonomy Proposal',
  type: 'document',
  icon: HelpCircleIcon,
  fields: [
    {
      name: 'proposedId',
      title: 'Proposed ID',
      type: 'string',
      validation: (rule) => rule.required()
    },
    {
      name: 'phenomenonData',
      title: 'Phenomenon Data',
      type: 'object',
      fields: [
        {name: 'surfaceForm', type: 'string'},
        {name: 'targetForm', type: 'string'},
        {name: 'grammaticalCategory', type: 'string'},
        {name: 'reasoningChain', type: 'text'}
      ]
    },
    {
      name: 'status',
      title: 'Status',
      type: 'string',
      options: {
        list: [
          {title: 'Pending Review', value: 'pending'},
          {title: 'Approved', value: 'approved'},
          {title: 'Rejected', value: 'rejected'},
          {title: 'Merged (Duplicate)', value: 'merged'}
        ],
        layout: 'radio'
      },
      initialValue: 'pending'
    },
    {
      name: 'submittedAt',
      title: 'Submitted At',
      type: 'datetime',
      initialValue: (new Date()).toISOString()
    },
    {
      name: 'reviewedAt',
      title: 'Reviewed At',
      type: 'datetime'
    },
    {
      name: 'reviewerNotes',
      title: 'Reviewer Notes',
      type: 'text'
    }
  ]
}
```

---

## PHASE 2: Context-Keeper Implementation (Go)

### 2.1 ProjectEntity Model
**File:** `internal/models/project_entity.go`

Implements the simplified sighting-log model with Ghost State.

**Key Points:**
- No mastery scores
- Just sightings + owned flag
- Optional heuristic band (4-13)

### 2.2 Student Resolver
**File:** `internal/utils/student_resolver.go`

Directory-based isolation:
```
data/students/{studentID}/
├── sessions/          # Context-Keeper sessions
├── sightings/         # ProjectEntity JSON files
└── corpus.duckdb      # Analytics database
```

### 2.3 Batch Ingestion Engine
**File:** `internal/ingestion/batch_processor.go`

```go
type BatchProcessor struct {
    llmClient      *llm.Client
    studentResolver *utils.StudentResolver
    sanityClient   *sanity.Client  // optional
}

func (bp *BatchProcessor) ProcessTranscript(
    ctx context.Context,
    studentID string,
    transcript string,
) (*IngestionResult, error) {
    // 1. Call LLM with SightingLogPrompt
    // 2. Parse ProjectEntity[] from response
    // 3. For each entity:
    //    - Check if exists in Sanity taxonomy
    //    - If exists: save to sightings/
    //    - If new: create taxonomyProposal in Sanity
    // 4. Update DuckDB with token frequencies
}
```

### 2.4 LLM Prompt Integration
**File:** `internal/llm/prompts/sighting_log.go`

Neutral observation prompt (no judgment, just logging).

Target: 600+ observations per session.

---

## PHASE 3: DuckDB Analytics

### 3.1 Schema
**File:** `migrations/001_initial.sql`

```sql
-- Per-student token frequencies
CREATE TABLE token_frequencies (
    student_id TEXT,
    token TEXT,
    count INTEGER DEFAULT 0,
    distinct_contexts INTEGER DEFAULT 0,
    first_seen TIMESTAMP,
    last_seen TIMESTAMP,
    PRIMARY KEY (student_id, token)
);

-- Per-student phenomenon sightings
CREATE TABLE phenomenon_sightings (
    student_id TEXT,
    phenomenon_id TEXT,
    sighting_count INTEGER DEFAULT 0,
    contexts TEXT[],  -- array of lesson contexts
    first_sighted TIMESTAMP,
    last_sighted TIMESTAMP,
    owned BOOLEAN DEFAULT FALSE,
    status TEXT,
    PRIMARY KEY (student_id, phenomenon_id)
);

-- Lexical spread metrics (materialized view)
CREATE VIEW lexical_spread AS
SELECT 
    student_id,
    token,
    count,
    distinct_contexts,
    CAST(distinct_contexts AS FLOAT) / count as spread_ratio
FROM token_frequencies
WHERE count > 1;

-- Acquisition velocity (sightings per week)
CREATE VIEW acquisition_velocity AS
SELECT 
    student_id,
    DATE_TRUNC('week', last_seen) as week,
    COUNT(DISTINCT phenomenon_id) as new_phenomena,
    SUM(sighting_count) as total_sightings
FROM phenomenon_sightings
GROUP BY student_id, week;
```

### 3.2 Integration
**File:** `internal/analytics/duckdb.go`

```go
type DuckDBStore struct {
    db *sql.DB
}

func (d *DuckDBStore) UpdateFromSightings(
    studentID string,
    sightings []*models.ProjectEntity,
) error {
    // Update phenomenon_sightings table
    // Update token_frequencies from transcript tokens
}

func (d *DuckDBStore) GetLexicalSpread(
    studentID string,
) (*LexicalSpreadReport, error) {
    // Query lexical_spread view
}
```

---

## PHASE 4: Batch Worker

### 4.1 Hourly Processor
**File:** `cmd/worker/main.go`

```go
func main() {
    // Load config
    // Initialize: LLM client, student resolver, DuckDB, Sanity client
    
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        // Find unprocessed transcripts
        // Process each with BatchProcessor
        // Log completion
    }
}
```

### 4.2 Deployment
Docker Compose service or systemd timer.

---

## PHASE 5: Sanity ↔ Context-Keeper Integration

### 5.1 Taxonomy Sync
**File:** `internal/sanity/sync.go`

```go
type TaxonomySync struct {
    client *sanity.Client
}

// Load approved phenomena from Sanity
func (ts *TaxonomySync) LoadTaxonomy() (map[string]bool, error) {
    // GROQ: *[_type == "linguisticPhenomenon"]{phenomenonId}
}

// Propose new phenomenon to Sanity
func (ts *TaxonomySync) Propose(entity *models.ProjectEntity) error {
    // Create taxonomyProposal document
}

// Check if phenomenon exists
func (ts *TaxonomySync) Exists(phenomenonID string) (bool, error) {
    // GROQ: count(*[_type == "linguisticPhenomenon" && phenomenonId == $id])
}
```

### 5.2 GROQ Queries

**Get all approved phenomena:**
```groq
*[_type == "linguisticPhenomenon"]
  | order(heuristicBand asc)
  {phenomenonId, surfaceForm, targetForm, grammaticalCategory}
```

**Get pending proposals:**
```groq
*[_type == "taxonomyProposal" && status == "pending"]
  | order(submittedAt desc)
```

**Get phenomena by band:**
```groq
*[_type == "linguisticPhenomenon" && heuristicBand == $band]
  {phenomenonId, surfaceForm}
```

---

## Execution Checklist

### Week 1: Foundation
- [ ] Deploy `linguisticPhenomenon` schema to Sanity
- [ ] Deploy `taxonomyProposal` schema to Sanity
- [ ] Implement `ProjectEntity` model
- [ ] Implement `StudentResolver`

### Week 2: Ingestion
- [ ] Implement batch processor
- [ ] Wire LLM prompt
- [ ] Create taxonomy matcher (exists vs propose)
- [ ] Test: transcript → sightings

### Week 3: Analytics
- [ ] Set up DuckDB schema
- [ ] Implement tokenization
- [ ] Implement frequency counting
- [ ] Create lexical spread views

### Week 4: Integration & Polish
- [ ] Build batch worker (hourly)
- [ ] Wire Sanity sync
- [ ] End-to-end test
- [ ] Deploy to production

---

## MCP Commands to Run

```bash
# 1. Deploy schemas
# Use: deploy_schema MCP tool with linguisticPhenomenon schema
# Use: deploy_schema MCP tool with taxonomyProposal schema

# 2. Verify deployment
# Use: get_schema MCP tool to verify schemas exist

# 3. Create initial seed phenomena (500)
# Use: create_documents_from_json MCP tool with seed data
```

---

## Success Criteria

- [ ] Transcript submitted → sightings saved within 1 hour
- [ ] New phenomena create draft proposals in Sanity
- [ ] DuckDB answers: "What's student_alpha's lexical spread?"
- [ ] You can query: "Show me all Band 7 phenomena for student_beta"
- [ ] Batch worker runs without errors
