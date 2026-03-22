package authz

import (
	"testing"
)

// testOwnable implements Ownable for testing.
type testOwnable struct {
	createdBy string
}

func (o *testOwnable) GetCreatedBy() string { return o.createdBy }

func TestDefaultPolicy_AdminCanDoEverything(t *testing.T) {
	p := &DefaultPolicy{}
	admin := UserFrom("admin-1", RoleAdmin)

	for _, action := range []string{ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionList} {
		if !p.Can(admin, action, &testOwnable{createdBy: "other-user"}) {
			t.Errorf("admin should be able to %s", action)
		}
	}
}

func TestDefaultPolicy_UserCanCreateAndList(t *testing.T) {
	p := &DefaultPolicy{}
	user := UserFrom("user-1", RoleUser)

	if !p.Can(user, ActionCreate, nil) {
		t.Error("user should be able to create")
	}
	if !p.Can(user, ActionList, nil) {
		t.Error("user should be able to list")
	}
}

func TestDefaultPolicy_OwnerCanEditDelete(t *testing.T) {
	p := &DefaultPolicy{}
	user := UserFrom("user-1", RoleUser)
	ownResource := &testOwnable{createdBy: "user-1"}

	if !p.Can(user, ActionUpdate, ownResource) {
		t.Error("owner should be able to update own resource")
	}
	if !p.Can(user, ActionDelete, ownResource) {
		t.Error("owner should be able to delete own resource")
	}
}

func TestDefaultPolicy_NonOwnerCannotEditDelete(t *testing.T) {
	p := &DefaultPolicy{}
	user := UserFrom("user-1", RoleUser)
	otherResource := &testOwnable{createdBy: "user-2"}

	if p.Can(user, ActionUpdate, otherResource) {
		t.Error("non-owner should NOT be able to update")
	}
	if p.Can(user, ActionDelete, otherResource) {
		t.Error("non-owner should NOT be able to delete")
	}
}

func TestDefaultPolicy_NonOwnerCanRead(t *testing.T) {
	p := &DefaultPolicy{}
	user := UserFrom("user-1", RoleUser)
	otherResource := &testOwnable{createdBy: "user-2"}

	if !p.Can(user, ActionRead, otherResource) {
		t.Error("non-owner should be able to read")
	}
}

func TestDefaultPolicy_NilUserDenied(t *testing.T) {
	p := &DefaultPolicy{}
	for _, action := range []string{ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionList} {
		if p.Can(nil, action, nil) {
			t.Errorf("nil user should be denied %s", action)
		}
	}
}

func TestDefaultPolicy_EmptyUserIDDenied(t *testing.T) {
	p := &DefaultPolicy{}
	emptyUser := UserFrom("", RoleUser)
	for _, action := range []string{ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionList} {
		if p.Can(emptyUser, action, nil) {
			t.Errorf("empty user ID should be denied %s", action)
		}
	}
}

func TestDefaultPolicy_NonOwnableResourceDeniesUpdateDelete(t *testing.T) {
	p := &DefaultPolicy{}
	user := UserFrom("user-1", RoleUser)

	if p.Can(user, ActionUpdate, nil) {
		t.Error("non-ownable resource should deny update for regular user")
	}
	if p.Can(user, ActionDelete, nil) {
		t.Error("non-ownable resource should deny delete for regular user")
	}
}

func TestCan_UsesRegisteredPolicy(t *testing.T) {
	// Register a custom policy that denies everything
	Register("test_resource", &denyAllPolicy{})
	defer func() {
		registryMu.Lock()
		delete(registry, "test_resource")
		registryMu.Unlock()
	}()

	admin := UserFrom("admin-1", RoleAdmin)
	if Can(admin, ActionRead, "test_resource", nil) {
		t.Error("custom deny-all policy should deny admin")
	}
}

func TestCan_FallsBackToDefaultPolicy(t *testing.T) {
	user := UserFrom("user-1", RoleUser)
	if !Can(user, ActionCreate, "unregistered_type", nil) {
		t.Error("unregistered type should use DefaultPolicy, which allows create")
	}
}

func TestIsAdmin(t *testing.T) {
	if !IsAdmin(UserFrom("1", RoleAdmin)) {
		t.Error("admin should be admin")
	}
	if IsAdmin(UserFrom("1", RoleUser)) {
		t.Error("user should not be admin")
	}
	if IsAdmin(nil) {
		t.Error("nil should not be admin")
	}
}

func TestUserFrom(t *testing.T) {
	u := UserFrom("abc", "editor")
	if u.GetID() != "abc" {
		t.Errorf("GetID() = %q, want %q", u.GetID(), "abc")
	}
	if u.GetRole() != "editor" {
		t.Errorf("GetRole() = %q, want %q", u.GetRole(), "editor")
	}
}

type denyAllPolicy struct{}

func (p *denyAllPolicy) Can(User, string, any) bool { return false }
