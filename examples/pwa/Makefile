DIST=./dist
ICONS=$(shell which icons)

deploy:
	$(if $(ICONS),\
	$(ICONS) -o $(DIST) --type favicon assets/favicon.png,\
	$(error you need icons tool. install: "pip3 install icons")\
	)
	spago deploy -tinygo $(DIST)
	cp -Rf assets manifest.json serviceworker.js offline.html $(DIST)

clean:
	@-rm -rf $(DIST)
