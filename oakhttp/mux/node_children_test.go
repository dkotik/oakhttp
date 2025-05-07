package mux

import (
	"strconv"
	"testing"
)

func BenchmarkNodeChildren(b *testing.B) {
	key := "articles"
	children := []string{
		"*",
		"cmd.html",
		"code.html",
		"contrib.html",
		"contribute.html",
		"debugging_with_gdb.html",
		"docs.html",
		"effective_go.html",
		"files.log",
		"gccgo_contribute.html",
		"gccgo_install.html",
		"go-logo-black.png",
		"go-logo-blue.png",
		"go-logo-white.png",
		"go1.1.html",
		"go1.2.html",
		"go1.html",
		"go1compat.html",
		"go_faq.html",
		"go_mem.html",
		"go_spec.html",
		"help.html",
		"ie.css",
		"install-source.html",
		"install.html",
		"logo-153x55.png",
		"Makefile",
		"root.html",
		"share.png",
		"sieve.gif",
		"tos.html",
		"articles",
	}
	if len(children) != 32 {
		panic("bad len")
	}
	for _, n := range []int{2, 4, 8, 12, 16, 32} {
		list := children[:n]
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.Run("linear", func(b *testing.B) {
				var entries listOfChildren
				for _, c := range list {
					entries = append(entries, child{c, nil})
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = entries.get(key)
				}
			})
			b.Run("map", func(b *testing.B) {
				entries := make(mapOfChildren)
				for _, c := range list {
					entries[c] = nil
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = entries.get(key)
				}
			})
			b.Run("legacy-linear", func(b *testing.B) {
				var entries []entry
				for _, c := range list {
					entries = append(entries, entry{c, nil})
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					findChildLinear(key, entries)
				}
			})
			b.Run("legacy-map", func(b *testing.B) {
				m := map[string]*node{}
				for _, c := range list {
					m[c] = nil
				}
				var x *node
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					x = m[key]
				}
				_ = x
			})
		})
	}
}
