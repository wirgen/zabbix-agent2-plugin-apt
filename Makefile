PACKAGE=zabbix-agent2-plugin-apt

DISTFILES = \
	go.mod \
	go.sum \
	LICENSE \
	main.go \
	Makefile \
	apt.conf \
	README.md

DIST_SUBDIRS = \
	plugin \
	vendor

build:
	go build -o "$(PACKAGE)"

clean:
	rm -rf ./vendor
	rm -rf ./$(PACKAGE)*
	go clean ./...

style:
	golangci-lint run --new-from-rev=$(NEW_FROM_REV) ./...

format:
	go fmt ./...

dist:
	go mod vendor; \
	major_version=$$(grep 'MajorVersion = ' ./vendor/git.zabbix.com/ap/plugin-support/plugin/comms/version.go | awk '{ print $$3 }'); \
	minor_version=$$(grep 'MinorVersion = ' ./vendor/git.zabbix.com/ap/plugin-support/plugin/comms/version.go | awk '{ print $$3 }'); \
	plugin_version=$$(grep 'pluginVersion =' ./main.go | awk '{ print $$4 }'); \
	distdir="$(PACKAGE)-$${major_version}.$${minor_version}.$${plugin_version}"; \
	dist_archive="$${distdir}.tar.gz"; \
	mkdir -p $${distdir}; \
	for distfile in '$(DISTFILES)'; do \
		cp -fp $${distfile} $${distdir}/; \
	done; \
	for subdir in '$(DIST_SUBDIRS)'; do \
		cp -fpR $${subdir} $${distdir}; \
	done; \
	tar -czvf $${dist_archive} $${distdir}; \
	rm -rf $${distdir}
