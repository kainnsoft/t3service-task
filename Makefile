docker.image.create:
	docker build -t task .

docker.container.run:
	docker run -dp 3000:3000 task

docker.compose.up:
	docker-compose up -d

swagger.init:
	swag init -g cmd/app/main.go -o internal/docs

test.unit:
	go test ./...

goose.create:
	goose create init2 sql

mock.gen:
	mockgen -source=".\internal\usecase\interfaces.go" -destination=".\mocks\usecase\moks.go" -package="repo_mocks"
	mockgen -source=".\internal\controller\http\v1\interfaces.go" -destination=".\mocks\controller\http\v1\moks.go" -package="ctrl_mocks"
