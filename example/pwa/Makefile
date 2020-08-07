DIST=./dist

deploy:
	spago deploy -tinygo $(DIST)
	cp -Rf assets manifest.json serviceworker.js offline.html $(DIST)
	icons -o $(DIST) --type favicon assets/favicon.png

clean:
	@-rm -rf $(DIST)
