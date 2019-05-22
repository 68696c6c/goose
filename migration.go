package goose

type MigrationInterface interface {
	GetUpSQL() string
	GetDownSQL() string
}

type Migration struct {
	Name string
	Up   string
	Down string
}

func (m *Migration) GetUpSQL() string {
	return m.Up
}

func (m *Migration) GetDownSQL() string {
	return m.Down
}

