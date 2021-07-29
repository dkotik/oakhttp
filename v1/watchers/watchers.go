/*

Package watchers contains consumers of events that can react to them.

*/
package watchers

import "github.com/dkotik/oakacs/v1"

type Watcher func() (chan (oakacs.Event), error)
