package cuestore

import "cuelang.org/go/cue/cuecontext"

var scope = cuecontext.New().CompileString(`
    #role: {
        name: string
        deny: [...#permission]
        allow: [...#permission]
        timeouts?: {
            session?: int
            idle?: int
        }
    }

    #permission: [...#match]
    #match: close({
        attribute: string
        operation?: string
        pattern: string
    })
`)
