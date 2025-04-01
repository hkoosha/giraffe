package env

const (
	// Valid values for a switch (case-insensitive):
	//
	// - yes
	// - true
	// - on
	// - 1
	// - enabled
	// - en
	//
	// Anything else sets the switch off.

	Giraffe = "GIRAFFE"

	Debug    = Giraffe + "_DEBUG"
	ToString = Debug + "_TO_STRING"
	Tracer   = Debug + "_TRACER"

	// UnsafeErrors giraffe will spit its guts out to whoever is connecting to
	// its http server.
	UnsafeErrors = Giraffe + "_UNSAFE_ERRORS"
)
