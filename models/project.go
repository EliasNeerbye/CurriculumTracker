package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
}

type Project struct {
	ID                 int         `json:"id"`
	CurriculumID       int         `json:"curriculum_id"`
	Identifier         string      `json:"identifier"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	LearningObjectives StringArray `json:"learning_objectives"`
	EstimatedTime      string      `json:"estimated_time"`
	Prerequisites      StringArray `json:"prerequisites"`
	ProjectType        string      `json:"project_type"`
	PositionOrder      int         `json:"position_order"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Progress           *Progress   `json:"progress,omitempty"`
}

type CreateProjectRequest struct {
	Identifier         string      `json:"identifier"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	LearningObjectives StringArray `json:"learning_objectives"`
	EstimatedTime      string      `json:"estimated_time"`
	Prerequisites      StringArray `json:"prerequisites"`
	ProjectType        string      `json:"project_type"`
	PositionOrder      int         `json:"position_order"`
}

type UpdateProjectRequest struct {
	Identifier         string      `json:"identifier"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	LearningObjectives StringArray `json:"learning_objectives"`
	EstimatedTime      string      `json:"estimated_time"`
	Prerequisites      StringArray `json:"prerequisites"`
	ProjectType        string      `json:"project_type"`
	PositionOrder      int         `json:"position_order"`
}

const (
	ProjectTypeRoot            = "root"
	ProjectTypeRootTest        = "rootTest"
	ProjectTypeBase            = "base"
	ProjectTypeBaseTest        = "baseTest"
	ProjectTypeLowerBranch     = "lowerBranch"
	ProjectTypeMiddleBranch    = "middleBranch"
	ProjectTypeUpperBranch     = "upperBranch"
	ProjectTypeFlowerMilestone = "flowerMilestone"
)
