package mmrtesting

import (
	"fmt"
	"math/rand"

	"github.com/datatrails/go-datatrails-merklelog/massifs/storage"
	"github.com/google/uuid"
)

// Option is a generic option type used for storage implementation testing.
// Implementations type assert to Options target record and if that fails the
// expectation they ignore the options
type Option func(any)

// TestOptions holds options generic for all storage implementations.
type TestOptions struct {
	// We seed the RNG of the provided StartTimeMS. It is normal to force it to
	// some fixed value so that the generated data is the same from run to run.
	StartTimeMS     int64
	EventRate       int
	TestLabelPrefix string
	LogID           storage.LogID // can be nil, defaults to TestLabelPrefix
	Rand            *rand.Rand
	WordList        []string // used for generating random words, defaults to bip32WordList
	LeafGenerator   LeafGenerator
	// Container       string        // can be "" defaults to TestLablePrefix
	// DebugLevel      string        // defaults to INFO

}

// WithDefaults sets the default values for TestOptions.
// The options are ap[plied in order]
// If you want to pre-empt it's choices preceeed it with the specific option.
// If you want to derive from a default, add your option after it.
// Typically WithStartTimeMS() would be the first option set, as that seeds the RNG.
func WithDefaults() Option {
	return func(o any) {
		options, ok := o.(*TestOptions)
		if !ok {
			return
		}
		if options.StartTimeMS == 0 {
			options.StartTimeMS = (1698342521) * 1000
		}

		if options.Rand == nil {
			options.Rand = rand.New(rand.NewSource(options.StartTimeMS / 1000))
		}
		if options.WordList == nil {
			options.WordList = bip32WordList()
		}

		if options.EventRate == 0 {
			options.EventRate = 500 // arbitrary default
		}
		if options.TestLabelPrefix == "" {
			a := options.WordList[options.Rand.Intn(len(options.WordList))]
			b := options.WordList[options.Rand.Intn(len(options.WordList))]
			options.TestLabelPrefix = fmt.Sprintf("mmrtesting.%s-%s", a, b)
		}
		if options.LogID == nil {
			id, err := uuid.NewRandomFromReader(options.Rand)
			if err != nil {
				panic("failed to generate random LogID: " + err.Error())
			}
			options.LogID = id[:]
		}
		if options.LeafGenerator == nil {
			options.LeafGenerator = MMRTestingGenerateNumberedLeaf
		}
	}
}

// WithStartTimeMS sets the StartTimeMS option for TestOptions.  This option
// determines the seed for the random number generator used in tests.  As with
// any option that should pre-empt the defaults,it must be placed before
// WithDefaults to take effect.
func WithStartTimeMS(startTimeMS int64) Option {
	return func(o any) {
		options, ok := o.(*TestOptions)
		if !ok {
			return
		}
		options.StartTimeMS = startTimeMS
	}
}

func WithLeafGenerator(leafGenerator LeafGenerator) Option {
	return func(o any) {
		options, ok := o.(*TestOptions)
		if !ok {
			return
		}
		options.LeafGenerator = leafGenerator
	}
}

// WithTestLebelPrefix pre-empts how the tests are identified. it is also
// typically used to isolate storage for integration tests
func WithTestLabelPrefix(prefix string) Option {
	return func(o any) {
		options, ok := o.(*TestOptions)
		if !ok {
			return
		}
		options.TestLabelPrefix = prefix
	}
}

func WithLogID(logID storage.LogID) Option {
	return func(o any) {
		options, ok := o.(*TestOptions)
		if !ok {
			return
		}
		options.LogID = logID
	}
}
