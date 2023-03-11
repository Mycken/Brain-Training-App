package store

type Store interface {
	User() UserRepository
}

type ResStore interface {
	Result() ResultRepository
}
