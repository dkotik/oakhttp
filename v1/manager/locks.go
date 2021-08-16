package oakmanager

import (
	"context"
	"fmt"

	"github.com/rs/xid"
)

type ErrIntegrityLockDenied struct {
	UUID xid.ID
}

func (e *ErrIntegrityLockDenied) Error() string {
	return fmt.Sprintf("failed to aquire resource lock: %s", e.UUID)
}

// Lock prevents objects from being altered until unlocked. If any of the objects are already locked, returns an error.
func (m *Manager) Lock(ctx context.Context, xid.ID...) (err error) {
    total := len(xid.ID)
    if total == 0 {
        return errors.New("cannot lock 0 objects")
    }
    return m.repo.Lock(ctx, xid.ID...)
}

// Unlock releases the lock.
func (m *Manager) Unlock(ctx context.Context, xid.ID...) (err error) {
    total := len(xid.ID)
    if total == 0 {
        return errors.New("cannot lock 0 objects")
    }
    return m.repo.Unlock(ctx, xid.ID...)
}

// PurgeLocks releases all the locks.
func (m *Manager) PurgeLocks(ctx context.Context) (err error) {
	if err = m.acs.Authorize(ctx, ACSService, DomainUniversal, "locks", "purge"); err != nil {
		return
	}
    m.acs.Broadcast(oakacs.Event{
        Context: ctx,
        Type: EventTypeMaintenance,
    })
    return m.repo.PurgeLocks(ctx)
}
