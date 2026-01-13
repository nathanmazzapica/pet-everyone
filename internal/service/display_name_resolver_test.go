package service

import (
	"errors"
	"testing"
	"time"

	"pet-everyone/internal/data/model"

	"github.com/google/uuid"
)

type fakeUserModel struct {
	users   map[uuid.UUID]string
	lookups int
}

func (f *fakeUserModel) Get(id *uuid.UUID) (*model.User, error) {
	f.lookups++
	if name, ok := f.users[*id]; ok {
		return &model.User{DisplayName: name}, nil
	}
	return nil, errors.New("not found")
}

type fakeVisitorModel struct {
	visitors map[uuid.UUID]string
	lookups  int
}

func (f *fakeVisitorModel) Get(id uuid.UUID) (*model.Visitor, error) {
	f.lookups++
	if name, ok := f.visitors[id]; ok {
		return &model.Visitor{DisplayName: name}, nil
	}
	return nil, errors.New("not found")
}

func TestDisplayNameResolverCachesResults(t *testing.T) {
	userID := uuid.New()
	userModel := &fakeUserModel{users: map[uuid.UUID]string{userID: "UserOne"}}
	visitorModel := &fakeVisitorModel{visitors: map[uuid.UUID]string{}}

	resolver := NewDisplayNameResolver(userModel, visitorModel, DisplayNameResolverConfig{TTL: 50 * time.Millisecond})

	actor := Actor{ID: userID.String(), IsGuest: false}

	name1, err := resolver.Resolve(actor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name1 != "UserOne" {
		t.Fatalf("expected UserOne, got %s", name1)
	}

	name2, err := resolver.Resolve(actor)
	if err != nil {
		t.Fatalf("unexpected error on cached resolve: %v", err)
	}
	if name2 != "UserOne" {
		t.Fatalf("expected cached UserOne, got %s", name2)
	}
	if userModel.lookups != 1 {
		t.Fatalf("expected single lookup, got %d", userModel.lookups)
	}
}

func TestDisplayNameResolverExpiresCache(t *testing.T) {
	userID := uuid.New()
	userModel := &fakeUserModel{users: map[uuid.UUID]string{userID: "UserOne"}}
	visitorModel := &fakeVisitorModel{visitors: map[uuid.UUID]string{}}

	resolver := NewDisplayNameResolver(userModel, visitorModel, DisplayNameResolverConfig{TTL: 5 * time.Millisecond})

	actor := Actor{ID: userID.String(), IsGuest: false}

	_, err := resolver.Resolve(actor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = resolver.Resolve(actor)
	if err != nil {
		t.Fatalf("unexpected error after expiry: %v", err)
	}

	if userModel.lookups < 2 {
		t.Fatalf("expected cache expiry to trigger new lookup")
	}
}
