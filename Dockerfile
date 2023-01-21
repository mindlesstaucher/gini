FROM golang:bullseye

WORKDIR /myworkdir

COPY . /myworkdir/

ENV PORT=8080

EXPOSE 8080

CMD ["go", "run", "main.go"]