package types

// Config is the top level structure to read from the user
type Config struct {
	Builds []Build `json:"builds"` // list of builds to call
}
