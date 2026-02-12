package workflow

import (
	"fmt"
	"sync"
)

type StepType string

const (
	StepTypeTransform StepType = "transform"
	StepTypeValidate  StepType = "validate"
	StepTypeFilter    StepType = "filter"
	StepTypeMap       StepType = "map"
)

type Step struct {
	ID     string                 `json:"id"`
	Type   StepType               `json:"type"`
	Config map[string]interface{} `json:"config"`
	Next   []string               `json:"next,omitempty"`
}

type Workflow struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Steps       []Step                 `json:"steps"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   string                 `json:"created_at"`
}

type Engine struct {
	workflows map[int]*Workflow
	nextID    int
	mu        sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		workflows: make(map[int]*Workflow),
		nextID:    1,
	}
}

func (e *Engine) Create(wf *Workflow) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	wf.ID = e.nextID
	e.nextID++
	e.workflows[wf.ID] = wf
	
	return wf.ID, nil
}

func (e *Engine) Get(id int) (*Workflow, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	wf, exists := e.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow %d not found", id)
	}
	
	return wf, nil
}

func (e *Engine) Execute(workflowID int, data map[string]interface{}) (map[string]interface{}, error) {
	wf, err := e.Get(workflowID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = v
	}

	for _, step := range wf.Steps {
		result, err = e.executeStep(step, result)
		if err != nil {
			return nil, fmt.Errorf("step %s failed: %w", step.ID, err)
		}
	}

	return result, nil
}

func (e *Engine) executeStep(step Step, data map[string]interface{}) (map[string]interface{}, error) {
	switch step.Type {
	case StepTypeTransform:
		return e.executeTransform(step, data)
	case StepTypeValidate:
		return e.executeValidate(step, data)
	case StepTypeFilter:
		return e.executeFilter(step, data)
	case StepTypeMap:
		return e.executeMap(step, data)
	default:
		return nil, fmt.Errorf("unknown step type: %s", step.Type)
	}
}

func (e *Engine) executeTransform(step Step, data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = v
	}
	
	if mapping, ok := step.Config["field_mapping"].(map[string]interface{}); ok {
		for oldKey, newKeyInterface := range mapping {
			if newKey, ok := newKeyInterface.(string); ok {
				if value, exists := result[oldKey]; exists {
					result[newKey] = value
					delete(result, oldKey)
				}
			}
		}
	}
	
	return result, nil
}

func (e *Engine) executeValidate(step Step, data map[string]interface{}) (map[string]interface{}, error) {
	return data, nil
}

func (e *Engine) executeFilter(step Step, data map[string]interface{}) (map[string]interface{}, error) {
	return data, nil
}

func (e *Engine) executeMap(step Step, data map[string]interface{}) (map[string]interface{}, error) {
	return data, nil
}