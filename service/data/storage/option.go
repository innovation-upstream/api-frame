package storage

// Option is for customizing storage implementations
type Option func(Storage)

// WithOwnerUIDField sets the `ownerUIDField` used to lookup data by owner uid
func WithOwnerUIDField(f string) Option {
	return func(r Storage) {
		r.setOwnerUIDFieldName(f)
	}
}

// WithInternalUIDField sets the `dataUIDField` used to lookup data by data uid
func WithInternalUIDField(f string) Option {
	return func(r Storage) {
		r.setInternalUIDFieldName(f)
	}
}

// WithExternalUIDField sets the `externalUIDField` used to lookup data by data uid
func WithExternalUIDField(f string) Option {
	return func(r Storage) {
		r.setExternalUIDFieldName(f)
	}
}
