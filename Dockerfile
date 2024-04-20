FROM golang:1.21
WORKDIR /myapp
COPY . .
ENV GOPRIVATE=gitlab.com/nina8884807/mail
RUN echo "machine gitlab.com login ninamusatova90 password glpat-61DSd-F9qJ4H9sZnqwwp" > $HOME/.netrc
RUN go mod download

RUN go build -o taskManager .
EXPOSE 8021
CMD ./taskManager