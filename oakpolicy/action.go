package oakpolicy

// An Action specifies the verb of a [Policy]. Package comes with a set of most frequently occuring actions. Specify custom actions as constants.
type Action string

// In returns true if the [Action] is in a given set.
func (a Action) In(set ...Action) bool {
	for _, action := range set {
		if action == a {
			return true
		}
	}
	return false
}

// Matches returns are of this [Action] is the equals the match or the wildcard [ActionAny].
func (a Action) Matches(match Action) bool {
	return a == match || match == ActionAny
}

const (
	ActionCreate   Action = "create"
	ActionRetrieve Action = "retrieve"
	ActionUpdate   Action = "update"
	ActionDelete   Action = "delete"
	ActionQuery    Action = "query"

	ActionAny       Action = "*"
	ActionAssign    Action = "assign"
	ActionUnassign  Action = "unassign"
	ActionBlock     Action = "block"
	ActionUnblock   Action = "unblock"
	ActionReset     Action = "reset"
	ActionRecover   Action = "recover"
	ActionPromote   Action = "promote"
	ActionDemote    Action = "demote"
	ActionUpgrade   Action = "upgrade"
	ActionDowngrade Action = "downgrade"
	ActionCommit    Action = "commit"
	ActionClear     Action = "clear"
	ActionInstall   Action = "install"
)
