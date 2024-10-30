package repository

type URLRepositoryState struct {
	state map[string]string
}

func CreateURLRepositoryState(state map[string]string) *URLRepositoryState {
	return &URLRepositoryState{
		state: state,
	}
}

func (m *URLRepositoryState) GetURLRepositoryState() map[string]string {
	return m.state
}
