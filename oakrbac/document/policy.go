package document

import (
	"errors"
	"strings"

	multierr "github.com/hashicorp/go-multierror"
)

type PolicyDefinition struct {
	// Name must begin with either "Allow" or "Deny" prefix. The prefix determines if the returned value when a policy matches will be either Allow or Deny.
	Name string
	// Description expresses the comment that will be added above the generated policy function.
	Description              string
	MatchResourcePathExactly string
	MatchResourcePathByMask  string
	MatchAnyAction           []string
	MatchEachPredicate       []PredicateDefinition
	MatchAnyPredicate        []PredicateDefinition
}

type PredicateDefinition struct {
	Property      string
	DesiredValues []string
}

func (p *PolicyDefinition) Validate() (err error) {
	if !strings.HasPrefix(p.Name, "Allow") && !strings.HasPrefix(p.Name, "Deny") {
		err = multierr.Append(err,
			errors.New("a policy name must begin with either \"Allow\" or \"Deny\" prefix"))
	}

	return
}
