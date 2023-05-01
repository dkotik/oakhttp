package document

import (
	"errors"
	"strings"
)

type PolicyDefinition struct {
	// Name must begin with either "Allow" or "Deny" prefix. The prefix determines if the returned value when a policy matches will be either Allow or Deny.
	Name string
	// Description expresses the comment that will be added above the generated policy function.
	Description       string
	MatchAction       []string
	MatchResourcePath [][]string
	MatchPredicates   []PredicateDefinition
}

type PredicateDefinition struct {
	Name          string
	DesiredValues []string
	Each          []PredicateDefinition
	Any           []PredicateDefinition
}

func (p *PolicyDefinition) Validate() (err error) {
	if !strings.HasPrefix(p.Name, "Allow") && !strings.HasPrefix(p.Name, "Deny") {
		return errors.New("a policy name must begin with either \"Allow\" or \"Deny\" prefix")
	}

	return
}
