// DROP-IN: scripts/sanity_bridge.go (or internal/utils/sanity_bridge.go)
//
// THE OPTIONAL BRIDGE
// 
// This file is ONLY used when YOU decide to push Ghost items to Sanity.
// Context-Keeper doesn't know about it. Sanity doesn't know about it.
// They're not married. They just hook up when you feel like it.
//
// Usage:
//   bridge := NewSanityBridge(projectID, dataset, token)
//   entity := loadEntityFromDisk("student_alpha", "grammar_past-perfect_001")
//   entity.PromoteToOwned("coaching")  // You decided this is worth coaching
//   err := bridge.PublishEntity(entity)  // Now it's in Sanity too

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/contextkeeper/service/internal/models"
)

// SanityBridge - Optional connector to Sanity CMS
type SanityBridge struct {
	ProjectID   string
	Dataset     string
	APIToken    string
	APIVersion  string
	HTTPClient  *http.Client
}

// NewSanityBridge creates bridge (call this when YOU want Sanity integration)
func NewSanityBridge(projectID, dataset, apiToken string) *SanityBridge {
	return &SanityBridge{
		ProjectID:  projectID,
		Dataset:    dataset,
		APIToken:   apiToken,
		APIVersion: "v2024-04-26",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// PublishEntity pushes a ProjectEntity to Sanity
// Only call this when YOU decide to promote a Ghost item
func (sb *SanityBridge) PublishEntity(entity *models.ProjectEntity) error {
	if !entity.Ghost.Owned {
		return fmt.Errorf("refusing to publish Ghost item: promote it first with PromoteToOwned()")
	}
	
	sanityDoc := entity.ToSanityDocument()
	
	mutation := map[string]interface{}{
		"mutations": []map[string]interface{}{
			{
				"createOrReplace": sanityDoc,
			},
		},
	}
	
	jsonData, err := json.Marshal(mutation)
	if err != nil {
		return fmt.Errorf("marshal mutation: %w", err)
	}
	
	url := fmt.Sprintf("https://%s.api.sanity.io/%s/data/mutate/%s", 
		sb.ProjectID, sb.APIVersion, sb.Dataset)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sb.APIToken)
	
	resp, err := sb.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sanity API error: %s", resp.Status)
	}
	
	return nil
}

// GenerateSanitySchema outputs a suggested Sanity schema
// Paste this into your Sanity Studio schema folder
// Then use Sanity's AI to refine it if you want
func GenerateSanitySchema() string {
	return `// sanity/schemaTypes/linguisticPhenomenon.js
// Paste this into your Sanity Studio
// Then tweak it however you want

export default {
  name: 'linguisticPhenomenon',
  title: 'Linguistic Phenomenon',
  type: 'document',
  fields: [
    {
      name: 'surfaceForm',
      title: 'Surface Form',
      type: 'string',
      description: 'What the student actually said/wrote'
    },
    {
      name: 'targetForm',
      title: 'Target Form', 
      type: 'string',
      description: 'The correct/expected form'
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
        ]
      }
    },
    {
      name: 'masteryCategory',
      title: 'Mastery Category',
      type: 'number',
      description: 'Curriculum level 4-13'
    },
    {
      name: 'ghostState',
      title: 'Ghost State',
      type: 'object',
      fields: [
        {
          name: 'owned',
          title: 'Owned (Active)',
          type: 'boolean',
          description: 'true = Active Toolkit, false = Ghost State'
        },
        {
          name: 'status',
          title: 'Status',
          type: 'string',
          options: {
            list: ['coaching', 'mastered', 'archived']
          }
        },
        {
          name: 'sightingCount',
          title: 'Times Sighted',
          type: 'number'
        }
      ]
    },
    {
      name: 'sightings',
      title: 'Sightings',
      type: 'array',
      of: [{
        type: 'object',
        fields: [
          {name: 'timestamp', type: 'datetime'},
          {name: 'lessonContext', type: 'string'},
          {name: 'utterance', type: 'text'},
          {name: 'correction', type: 'text'}
        ]
      }]
    },
    {
      name: 'masteryScore',
      title: 'Mastery Score',
      type: 'number'
    },
    {
      name: 'weightedMastery',
      title: 'Weighted Mastery (W-final)',
      type: 'number',
      description: 'Calculated: Σ(Mastery·Confidence)/N'
    },
    {
      name: 'validityConfidence',
      title: 'Validity Confidence',
      type: 'number',
      description: 'LLM certainty 0.0-1.0'
    },
    {
      name: 'reasoningChain',
      title: 'Reasoning Chain',
      type: 'text',
      description: 'LLM deduction path'
    }
  ],
  preview: {
    select: {
      title: 'surfaceForm',
      subtitle: 'grammaticalCategory'
    }
  }
}
`
}

// Example main() showing usage
/*
func main() {
	// Load entity from Context-Keeper storage
	entity := loadEntity("student_alpha", "grammar_past-perfect_001")
	
	// You decided this is worth coaching now
	entity.PromoteToOwned("coaching")
	
	// Optional: Push to Sanity
	bridge := NewSanityBridge(
		os.Getenv("SANITY_PROJECT_ID"),
		os.Getenv("SANITY_DATASET"),
		os.Getenv("SANITY_API_TOKEN"),
	)
	
	if err := bridge.PublishEntity(entity); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Published to Sanity! Now go coach it.")
}
*/
