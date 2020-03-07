FROM node:13.4.0-alpine3.10 as frontend
COPY . argovue
RUN cd argovue/ui && yarn install
RUN cd argovue/ui && yarn build
RUN apk add git && \
	cd argovue && \
	export VERSION=$(git describe --tags) && sed -i "s/_VERSION_/$VERSION/" ui/dist/config.js && \
	export COMMIT=$(git rev-parse --short HEAD) && sed -i "s/_COMMIT_/$COMMIT/" ui/dist/config.js && \
	export BUILDDATE=$(date +%Y%m%d%H%M%S) && sed -i "s/_BUILDDATE_/$BUILDDATE/" ui/dist/config.js

FROM golang:1.13-alpine as backend
COPY . /home/argovue
RUN apk add git && \
	cd /home/argovue && \
	export VERSION=$(git describe --tags) && \
	export COMMIT=$(git rev-parse --short HEAD) && \
	export BUILDDATE=$(date +%Y%m%d%H%M%S) && \
	cd src && go build -ldflags="-X main.version=$VERSION -X main.builddate=$BUILDDATE -X main.commit=$COMMIT"

FROM alpine:3.10
RUN apk update
COPY --from=backend /home/argovue/src/argovue argovue
COPY --from=frontend argovue/ui/dist ui/dist
CMD ["./argovue"]
