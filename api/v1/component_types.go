package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ComponentSpec   `json:"spec,omitempty"`
	Status            ComponentStatus `json:"status,omitempty"`
}

type ComponentSpec struct {
	Name    string    `json:"name,omitempty"`
	Tooltip string    `json:"tooltip,omitempty"`
	Icon    string    `json:"icon,omitempty"`
	Owner   string    `json:"owner,omitempty"`
	Id      *Template `json:"id,omitempty"` //nolint
	Order   int       `json:"order,omitempty"`
	// The type of component, e.g. service, API, website, library, database, etc.
	Type string `json:"type,omitempty"`
	// The lifecycle state of the component e.g. production, staging, dev, etc.
	Lifecycle       string             `json:"lifecycle,omitempty"`
	Relationships   []RelationshipSpec `json:"relationships,omitempty"`
	Properties      []*Property        `json:"properties,omitempty"`
	Lookup          *CanarySpec        `json:"lookup,omitempty"`
	Components      []json.RawMessage  `json:"components,omitempty"`
	Selectors       ResourceSelectors  `json:"selectors,omitempty"`
	ComponentChecks ComponentChecks    `json:"checks,omitempty"`
	Configs         []Config           `json:"configs,omitempty"`
	Summary         *Summary           `json:"summary,omitempty"`
}

type Summary struct {
	Healthy   int `json:"healthy,omitempty"`
	Unhealthy int `json:"unhealthy,omitempty"`
	Warning   int `json:"warning,omitempty"`
	Info      int `json:"info,omitempty"`
}

func (s Summary) String() string {
	str := ""
	if s.Unhealthy > 0 {
		str += fmt.Sprintf("unhealthy=%d ", s.Unhealthy)
	}
	if s.Warning > 0 {
		str += fmt.Sprintf("warning=%d ", s.Warning)
	}
	if s.Healthy > 0 {
		str += fmt.Sprintf("healthy=%d ", s.Healthy)
	}
	return strings.TrimSpace(str)
}

func (s Summary) GetStatus() ComponentPropertyStatus {
	if s.Unhealthy > 0 {
		return ComponentPropertyStatusUnhealthy
	} else if s.Warning > 0 {
		return ComponentPropertyStatusWarning
	} else if s.Healthy > 0 {
		return ComponentPropertyStatusHealthy
	}
	return "unknown"
}

func (s Summary) Add(b Summary) Summary {
	if b.Healthy > 0 && b.Unhealthy > 0 {
		s.Warning += 1
	} else if b.Unhealthy > 0 {
		s.Unhealthy += 1
	} else if b.Healthy > 0 {
		s.Healthy += 1
	}
	if b.Warning > 0 {
		s.Warning += b.Warning
	}
	if b.Info > 0 {
		s.Info += b.Info
	}
	return s
}

type ComponentStatus struct {
	Status ComponentPropertyStatus `json:"status,omitempty"`
}

type RelationshipSpec struct {
	// The type of relationship, e.g. dependsOn, subcomponentOf, providesApis, consumesApis
	Type string `json:"type,omitempty"`
	Ref  string `json:"ref,omitempty"`
}

type ComponentPropertyStatus string

var (
	ComponentPropertyStatusHealthy   ComponentPropertyStatus = "healthy"
	ComponentPropertyStatusUnhealthy ComponentPropertyStatus = "unhealthy"
	ComponentPropertyStatusWarning   ComponentPropertyStatus = "warning"
	ComponentPropertyStatusError     ComponentPropertyStatus = "error"
	ComponentPropertyStatusInfo      ComponentPropertyStatus = "info"
)

type Owner string

type Text struct {
	Tooltip string `json:"tooltip,omitempty"`
	Icon    string `json:"icon,omitempty"`
	Text    string `json:"text,omitempty"`
	Label   string `json:"label,omitempty"`
}

type Link struct {
	// e.g. documentation, support, playbook
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
	Text `json:",inline"`
}

type Properties []Property

type Property struct {
	Label    string `json:"label,omitempty"`
	Name     string `json:"name,omitempty"`
	Tooltip  string `json:"tooltip,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Text     string `json:"text,omitempty"`
	Order    int    `json:"order,omitempty"`
	Headline bool   `json:"headline,omitempty"`
	Type     string `json:"type,omitempty"`
	Color    string `json:"color,omitempty"`
	// e.g. milliseconds, bytes, millicores, epoch etc.
	Unit           string        `json:"unit,omitempty"`
	Value          int64         `json:"value,omitempty"`
	Max            *int64        `json:"max,omitempty"`
	Min            int64         `json:"min,omitempty"`
	Status         string        `json:"status,omitempty"`
	LastTransition string        `json:"lastTransition,omitempty"`
	Links          []Link        `json:"links,omitempty"`
	Lookup         *CanarySpec   `json:"lookup,omitempty"`
	ConfigLookup   *ConfigLookup `json:"configLookup,omitempty"`
	Summary        *Template     `json:"summary,omitempty"`
}

type ConfigLookup struct {
	ID      string  `json:"id"`
	Config  Config  `json:"config,omitempty"`
	Field   string  `json:"field,omitempty"`
	Display Display `json:"display,omitempty"`
}

// +kubebuilder:object:root=true

// ComponentList contains a list of Canary
type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Component `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Component{}, &ComponentList{})
}
