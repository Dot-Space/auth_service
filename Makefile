# General
define set-default-container
	ifndef c
	c = exam-service
	else ifeq (${c},all)
	override c=
	endif
endef

# General
TAIL=100

set-container:
	$(eval $(call set-default-container))


build:
	docker compose -f docker-compose.local.yml build
up:
	docker compose -f docker-compose.local.yml up -d
down:
	docker compose -f docker-compose.local.yml down
logs: set-container
	docker compose -f docker-compose.local.yml logs --tail=$(TAIL) -f $(c)
restart: set-container
	docker compose -f docker-compose.local.yml restart $(c)
exec: set-container
	docker compose -f docker-compose.local.yml exec $(c) /bin/bash
startup: build up


# Proto build GO
build-proto:
	protoc -I proto gen-proto/proto/exam/*.proto --go_out=./gen-proto/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen-proto/gen/go/ --go-grpc_opt=paths=source_relative