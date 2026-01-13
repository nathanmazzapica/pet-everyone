package service

import (
	"errors"
	"sync"
	"time"

	"pet-everyone/internal/data/model"

	"github.com/google/uuid"
)

// DisplayNameResolver resolves display names for users and guests.
type DisplayNameResolver interface {
	Resolve(actor Actor) (string, error)
}

// DisplayNameResolverConfig configures cache behavior.
type DisplayNameResolverConfig struct {
	TTL time.Duration
}

const DefaultDisplayNameCacheTTL = 10 * time.Minute

type displayNameResolver struct {
	userModel    userLookup
	visitorModel visitorLookup
	cache        map[string]cachedDisplayName
	mu           sync.RWMutex
	tl           time.Duration
}

type userLookup interface {
	Get(id *uuid.UUID) (*model.User, error)
}

type visitorLookup interface {
	Get(id uuid.UUID) (*model.Visitor, error)
}

type cachedDisplayName struct {
	name    string
	expires time.Time
}

func displayNameCacheKey(actor Actor) string {
	if actor.IsGuest {
		return "guest|" + actor.ID
	}
	return "user|" + actor.ID
}

// NewDisplayNameResolver constructs a resolver backed by user/visitor models with TTL cache.
func NewDisplayNameResolver(userModel userLookup, visitorModel visitorLookup, cfg DisplayNameResolverConfig) DisplayNameResolver {
	ttl := cfg.TTL
	if ttl == 0 {
		ttl = DefaultDisplayNameCacheTTL
	}

	return &displayNameResolver{
		userModel:    userModel,
		visitorModel: visitorModel,
		cache:        make(map[string]cachedDisplayName),
		tl:           ttl,
	}
}

func (r *displayNameResolver) Resolve(actor Actor) (string, error) {
	key := displayNameCacheKey(actor)

	if name, ok := r.lookupCache(key); ok {
		return name, nil
	}

	name, err := r.resolveFresh(actor)
	if err != nil {
		return "", err
	}

	r.saveCache(key, name)
	return name, nil
}

func (r *displayNameResolver) resolveFresh(actor Actor) (string, error) {
	if actor.ID == "" {
		return "", errors.New("missing actor id")
	}

	uuidVal, err := uuid.Parse(actor.ID)
	if err != nil {
		return "", err
	}

	if actor.IsGuest {
		visitor, verr := r.visitorModel.Get(uuidVal)
		if verr != nil {
			return "", verr
		}
		return visitor.DisplayName, nil
	}

	user, uerr := r.userModel.Get(&uuidVal)
	if uerr != nil {
		return "", uerr
	}
	return user.DisplayName, nil
}

func (r *displayNameResolver) lookupCache(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.cache[key]
	if !ok || time.Now().After(entry.expires) {
		return "", false
	}
	return entry.name, true
}

func (r *displayNameResolver) saveCache(key string, name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[key] = cachedDisplayName{name: name, expires: time.Now().Add(r.tl)}
}
