package document

type TestCaseDefinition struct {
	Action         string
	ResourcePath   []string
	MockPredicates map[string]bool
	ExpectAllow    bool
	ExpectDeny     bool
	ExpectOmit     bool
}
