# Book API Application

This Application an API for the **BookStore manegment**. 


[BookStore](https://github.com/JuniorPaula/book-store)

## Setup  

To set up this project locally, you need to have `Golang` and `postgres` installed. Then, you can clone the repository and install the dependencies:

**Note**
Make sure the configurate correct database before start the application.

```bash
git clone
cd repository

go mod tidy
go build ./cmd/api/*.go
```

If you have a `Makefile` installed in yours machine, you can run:
```bash
make start
```
or
```bash
make stop
```
or

```bash
make restart
```

To start the application, run:

```bash
./main
```
The application will be available at `http://localhost:8081`