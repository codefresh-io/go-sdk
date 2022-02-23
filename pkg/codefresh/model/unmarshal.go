package model

import "encoding/json"

// This file includes unmarshal overrides for types that can
// not be simply unmarshal otherwise.
//
// For example, if one of the fields on the type you are trying
// to unmarshal to is an interface, you need to provide custom
// behavior that decides to which concrete type it should be
// umarshaled into.

// Application entity
type ApplicationJSON struct {
	// Object metadata
	Metadata *ObjectMeta `json:"metadata"`
	// Errors
	Errors []SyncError `json:"errors"`
	// Entities referencing this entity
	ReferencedBy []BaseEntity `json:"referencedBy"`
	// Entities referenced by this enitity
	References []BaseEntity `json:"references"`
	// Relations between parents and child applications in tree
	AppsRelations *AppsRelations `json:"appsRelations"`
	// History of the application
	History *GitOpsSlice `json:"history"`
	// Version of the entity (generation)
	Version *int `json:"version"`
	// Is this the latest version of this entity
	Latest *bool `json:"latest"`
	// Entity source
	Source *GitopsEntitySource `json:"source"`
	// Sync status
	SyncStatus SyncStatus `json:"syncStatus"`
	// Health status
	HealthStatus *HealthStatus `json:"healthStatus"`
	// Health message
	HealthMessage *string `json:"healthMessage"`
	// Desired manifest
	DesiredManifest *string `json:"desiredManifest"`
	// Actual manifest
	ActualManifest *string `json:"actualManifest"`
	// Projects
	Projects []string `json:"projects"`
	// Updated At
	UpdatedAt *string `json:"updatedAt"`
	// Path
	Path *string `json:"path"`
	// RepoURL
	RepoURL *string `json:"repoURL"`
	// Number of resources
	Size *int `json:"size"`
	// Revision
	Revision *string `json:"revision"`
	// Status
	Status *ArgoCDApplicationStatus `json:"status"`
	// Favorites
	Favorites []string `json:"favorites"`
	// Argo CD application destination config
	Destination *ArgoCDApplicationDestination `json:"destination"`
}

func (a *Application) UnmarshalJSON(data []byte) error {
	aj := ApplicationJSON{}
	if err := json.Unmarshal(data, &aj); err != nil {
		return err
	}

	errs := make([]Error, len(aj.Errors))
	for i := range aj.Errors {
		errs[i] = aj.Errors[i]
	}

	a.Metadata = aj.Metadata
	a.Errors = errs
	a.ReferencedBy = aj.ReferencedBy
	a.References = aj.References
	a.AppsRelations = aj.AppsRelations
	a.History = aj.History
	a.Version = aj.Version
	a.Latest = aj.Latest
	a.Source = aj.Source
	a.SyncStatus = aj.SyncStatus
	a.HealthStatus = aj.HealthStatus
	a.HealthMessage = aj.HealthMessage
	a.DesiredManifest = aj.DesiredManifest
	a.ActualManifest = aj.ActualManifest
	a.Projects = aj.Projects
	a.UpdatedAt = aj.UpdatedAt
	a.Path = aj.Path
	a.RepoURL = aj.RepoURL
	a.Size = aj.Size
	a.Revision = aj.Revision
	a.Status = aj.Status
	a.Favorites = aj.Favorites
	a.Destination = aj.Destination

	return nil
}
