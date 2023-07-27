package store

type Store interface {
	getUserRepo() *UserRepository
}
