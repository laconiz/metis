package session

import "sync"

func NewManager() *Manager {
	return &Manager{sessions: map[uint64]*Session{}}
}

type Manager struct {
	sessions map[uint64]*Session
	mutex    sync.RWMutex
}

func (mgr *Manager) Load(id uint64) *Session {

	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	return mgr.sessions[id]
}

func (mgr *Manager) Count() int64 {

	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	return int64(len(mgr.sessions))
}

func (mgr *Manager) Insert(session *Session) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.sessions[session.ID()] = session
}

func (mgr *Manager) Remove(session *Session) {

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	delete(mgr.sessions, session.ID())
}

func (mgr *Manager) Range(handler func(*Session) bool) {

	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	for _, ses := range mgr.sessions {
		if !handler(ses) {
			return
		}
	}
}
