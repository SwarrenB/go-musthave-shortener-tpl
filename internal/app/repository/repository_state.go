package repository

type URLRepositoryState struct {
	state map[string]Record
}

func CreateURLRepositoryState(state map[string]Record) *URLRepositoryState {
	return &URLRepositoryState{
		state: state,
	}
}

func (m *URLRepositoryState) GetURLRepositoryState() map[string]Record {
	return m.state
}
