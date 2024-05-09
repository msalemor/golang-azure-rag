default:

.Phony : clean
clean:
	rm -rf src/backend/wwwroot
	mkdir src/backend/wwwroot

.Phony : build-ui
build-ui: clean
	cd src/frontend && bun run prod
	cp -r src/frontend/dist/* src/backend/wwwroot

.Phony : run-ui 
run-ui: build-ui
	cd src/backend && go run .

.Phony : docker-build
docker-build: build-ui
	cd src/backend && docker build . -t am8850/gorag:0.0.0

.Phony : docker-run
docker-run:
	cd src/backend && docker run --rm -p 8090:8080 --env-file=.env am8850/gorag:0.0.0