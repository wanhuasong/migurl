build:
	CGO_ENABLE=0 go build -trimpath -o bin/migurl main.go

upload:
	rsync --progress -v bin/migurl new-marsdev:/tmp