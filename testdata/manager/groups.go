package oakmanager

const groupResource = "group"

// // CreateGroup creates a new group.
// func (m *Manager) CreateGroup(ctx context.Context, name string) (identity.Group, error) {
// 	// if err := m.acs.Authorize(ctx, ACSService, DomainUniversal, groupResource, "create"); err != nil {
// 	// 	return
// 	// }
// 	// id, err := m.persistent.Groups.Create(ctx, name)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// group, err := m.persistent.Groups.Retrieve(ctx, id)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// return group, nil
// 	return nil, nil
// }

// // RetrieveGroup fetches a group.
// func (m *Manager) RetrieveGroup(ctx context.Context, uuid xid.ID) (identity.Group, error) {
// 	err := m.acs.Authorize(ctx, service, domain, groupResource, "retrieve")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m.persistent.Groups.Retrieve(ctx, uuid)
// }

// // DeleteGroup removes the group from the backend.
// func (m *Manager) DeleteGroup(ctx context.Context, uuid xid.ID) (err error) {
// 	err = m.acs.Authorize(ctx, service, domain, groupResource, "delete")
// 	if err != nil {
// 		return
// 	}
// 	// members, err := m.persistent.Groups.ListMembers(ctx, &oakquery.Query{
// 	// 	PerPage: 5000,
// 	// })
// 	// if err != nil {
// 	// 	return
// 	// }
// 	// if l := len(members); l > 0 {
// 	// 	return fmt.Errorf("cannot delete a group because it has %d members", l)
// 	// }
// 	return m.persistent.Groups.Delete(ctx, uuid)
// }
