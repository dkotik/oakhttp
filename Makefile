default:
	cd v1 && for s in $$(go list ./...); do if ! go test -failfast -v -p 1 -timeout 60s $$s; then break; fi; done
