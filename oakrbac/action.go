package oakrbac

type (
	// An Action specifies the verb of an [Intention]. OakRBAC comes with a set of most frequently occuring actions. Specify custom actions as constants.
	Action string
)

// In returns true if the [Action] matches one of the provided set.
func (a Action) In(set ...Action) bool {
	for _, action := range set {
		if action == a {
			return true
		}
	}
	return false
}

const (
	ActionCreate   Action = "create"
	ActionRetrieve Action = "retrieve"
	ActionUpdate   Action = "update"
	ActionDelete   Action = "delete"
	ActionQuery    Action = "query"

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
