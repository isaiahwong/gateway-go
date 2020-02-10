package server

// Payload defines a generic type for poly data types that are not known
// to be encapsulated into the `Body` property. This is makes it convenient for
// forwarding payloads from third party services
type Payload struct {
	Body []byte `json:"body"`
}
