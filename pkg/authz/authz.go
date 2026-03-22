// Package authz provides authorization primitives for generated applications.
// It defines the Policy interface, a DefaultPolicy with role-based and
// ownership-based access control, and a global policy registry.
package authz

import "sync"

// Standard actions for CRUD operations.
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"
)

// Standard roles.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User is the minimal interface for authorization checks.
// The generated sqlc models.User struct can be adapted via UserFrom().
type User interface {
	GetID() string
	GetRole() string
}

// Ownable is implemented by resources that track ownership via a created_by column.
type Ownable interface {
	GetCreatedBy() string
}

// Policy defines authorization rules for a resource type.
type Policy interface {
	Can(user User, action string, resource any) bool
}

// DefaultPolicy implements reasonable defaults:
//   - Admin can do everything
//   - Any authenticated user can create and list
//   - Owner can read, update, delete their own resources
//   - Non-owner can read but not update/delete
type DefaultPolicy struct{}

func (p *DefaultPolicy) Can(user User, action string, resource any) bool {
	if user == nil {
		return false
	}
	if user.GetID() == "" {
		return false
	}
	if user.GetRole() == RoleAdmin {
		return true
	}
	switch action {
	case ActionCreate, ActionList:
		return true
	case ActionRead:
		return true
	case ActionUpdate, ActionDelete:
		if ownable, ok := resource.(Ownable); ok {
			return ownable.GetCreatedBy() == user.GetID()
		}
		return false
	}
	return false
}

// registry stores policies keyed by resource type name.
var (
	registry   = make(map[string]Policy)
	registryMu sync.RWMutex
)

// Register registers a policy for a resource type name.
// If no policy is registered for a type, DefaultPolicy is used.
func Register(resourceType string, p Policy) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[resourceType] = p
}

// Can checks if a user can perform an action on a resource.
// It looks up the policy by resourceType; falls back to DefaultPolicy.
func Can(user User, action string, resourceType string, resource any) bool {
	registryMu.RLock()
	p, ok := registry[resourceType]
	registryMu.RUnlock()
	if !ok {
		p = &DefaultPolicy{}
	}
	return p.Can(user, action, resource)
}

// IsAdmin returns true if the user has the admin role.
func IsAdmin(user User) bool {
	return user != nil && user.GetRole() == RoleAdmin
}

// userAdapter adapts arbitrary user data to the User interface.
type userAdapter struct {
	id   string
	role string
}

func (u *userAdapter) GetID() string   { return u.id }
func (u *userAdapter) GetRole() string { return u.role }

// UserFrom creates a User from an ID and role string.
// Useful for adapting sqlc-generated models: authz.UserFrom(user.ID, user.Role)
func UserFrom(id, role string) User {
	return &userAdapter{id: id, role: role}
}

// ownedResource wraps a createdBy string to implement Ownable.
type ownedResource struct{ createdBy string }

func (o *ownedResource) GetCreatedBy() string { return o.createdBy }

// OwnedBy returns an Ownable resource for a given creator ID.
// Use when the sqlc model struct can't implement Ownable directly:
//
//	authz.Can(user, authz.ActionUpdate, "posts", authz.OwnedBy(post.CreatedBy))
func OwnedBy(createdBy string) Ownable {
	return &ownedResource{createdBy: createdBy}
}
