package errors

type AuthorizationError string

func (e AuthorizationError) Error() string                { return string(e) }
func NewAuthorizationError(msg string) AuthorizationError { return AuthorizationError(msg) }

type TargetIdentificationError string

func (e TargetIdentificationError) Error() string { return string(e) }
func NewTargetIdentificationError(msg string) TargetIdentificationError {
	return TargetIdentificationError(msg)
}

type StorageLoadingError string

func (e StorageLoadingError) Error() string                 { return string(e) }
func NewStorageLoadingError(msg string) StorageLoadingError { return StorageLoadingError(msg) }

type StorageSavingError string

func (e StorageSavingError) Error() string                { return string(e) }
func NewStorageSavingError(msg string) StorageSavingError { return StorageSavingError(msg) }

type DomainError string

func (e DomainError) Error() string         { return string(e) }
func NewDomainError(msg string) DomainError { return DomainError(msg) }
