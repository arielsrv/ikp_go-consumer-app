# Compile stage
FROM  golang:1.20.1-bullseye AS build

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

ADD . /app
WORKDIR /app

RUN task build

EXPOSE ${PORT} ${PORT}

CMD ["./build/program"]
