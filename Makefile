NAME=logger
VERSION=1.2
PORT_MAP=3021:3021

default: buildp

# production build
buildp: buildgop build push

# development build and run
buildd: buildgop build rund

# build go for production
buildgop:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o $(NAME) .

# build go for development
# buildgod:
	# go build -o $(NAME) .

# standard docker build
build:
	docker build --rm --tag=pulpfree/$(NAME):$(VERSION) .

# docker run for development
# docker run --name logger -p 3021:3021 --env environment=dev pulpfree/logger:1.2
rund:
	# docker run --name $(NAME) -p $(PORT_MAP) --env environment=dev --rm pulpfree/$(NAME):$(VERSION)
	docker run --name $(NAME) -p $(PORT_MAP) --env environment=dev -d pulpfree/$(NAME):$(VERSION)

# docker run for production
# sudo docker run --name pf-auth -p 3021:3021 --env environment=prod -d pulpfree/logger:1.2
runp:
	docker run --name $(NAME) -p $(PORT_MAP) --env environment=prod -d pulpfree/$(NAME):$(VERSION)

push:
	docker push pulpfree/$(NAME):$(VERSION)

start:
	docker start $(NAME)

stop:
	docker stop $(NAME)