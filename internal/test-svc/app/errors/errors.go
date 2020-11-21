package errors

type AuthorizationError string

func (e AuthorizationError) Error() string                { return string(e) }
func NewAuthorizationError(msg string) AuthorizationError { return AuthorizationError(msg) }

type IdentificationError string

func (e IdentificationError) Error() string                 { return string(e) }
func NewIdentificationError(msg string) IdentificationError { return IdentificationError(msg) }

type RepositoryError string

func (e RepositoryError) Error() string             { return string(e) }
func NewRepositoryError(msg string) RepositoryError { return RepositoryError(msg) }

type DomainError string

func (e DomainError) Error() string         { return string(e) }
func NewDomainError(msg string) DomainError { return DomainError(msg) }
