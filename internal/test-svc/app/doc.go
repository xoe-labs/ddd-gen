// Package app declares interfaces which the application layer either requires or offers.
//
// By convention, the following prefixes further qualify the interfaces:
// 	Offers*
// 	Requires*
//
// The name of the go file (e.g. `storage.go`) signifies the adapter or object of the interface.
//
// Names terminating in 'able' represent types for which app offers an implementation:
// Adapters or ports shall understand those interface types as common language comming from external services
// Hence, their implementation is part of the package's public api.
package app
