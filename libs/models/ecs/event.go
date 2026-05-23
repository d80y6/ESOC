package ecs

import "time"

// Base ECS Fields
type Base struct {
	Timestamp time.Time `json:"@timestamp"`
	Labels    map[string]string `json:"labels,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Message   string `json:"message,omitempty"`
}

type Event struct {
	ID          string    `json:"id,omitempty"`
	Action      string    `json:"action,omitempty"`
	Category    []string  `json:"category,omitempty"`
	Outcome     string    `json:"outcome,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Dataset     string    `json:"dataset,omitempty"`
	Module      string    `json:"module,omitempty"`
	Severity    int64     `json:"severity,omitempty"`
	RiskScore   float64   `json:"risk_score,omitempty"`
}

type Agent struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
}

type Host struct {
	Hostname string   `json:"hostname,omitempty"`
	IP       []string `json:"ip,omitempty"`
	ID       string   `json:"id,omitempty"`
	Type     string   `json:"type,omitempty"`
}

type Source struct {
	IP      string `json:"ip,omitempty"`
	Port    int64  `json:"port,omitempty"`
	Domain  string `json:"domain,omitempty"`
	Address string `json:"address,omitempty"`
}

type Destination struct {
	IP      string `json:"ip,omitempty"`
	Port    int64  `json:"port,omitempty"`
	Domain  string `json:"domain,omitempty"`
	Address string `json:"address,omitempty"`
}

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NormalizedEvent is the container for all ECS fields
type NormalizedEvent struct {
	Base
	Event       Event       `json:"event"`
	Agent       *Agent      `json:"agent,omitempty"`
	Host        *Host       `json:"host,omitempty"`
	Source      *Source     `json:"source,omitempty"`
	Destination *Destination `json:"destination,omitempty"`
	User        *User       `json:"user,omitempty"`

	// Custom fields or metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
