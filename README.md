# Deliveroo Take-Home Project

This project is a Go-based implementation of a cron expression parser. It expands cron expressions into a list of execution times. With `make coverage` you can see that there is 94.9% test coverage. You can use the command line tool like this:
```shell
$ git clone https://github.com/dantheman0207/deliveroo-take-home.git
$ cd deliveroo-take-home
$ make
$ ./deliveroo-take-home "*/15" 0 1,15 "*" 1-5 /usr/bin/find
minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```


I have adapted this from the requirements slightly to pass the arguments with quotes (`"*/15"`) around them to work on `zsh` as well as `bash`. This has been tested on macOS.

## Getting Started

### Prerequisites

- Go 1.16 or later

### Installation

Clone the repository:
```shell
git clone https://github.com/dantheman0207/deliveroo-take-home.git
cd deliveroo-take-home
```

## Usage

This project includes a Makefile to simplify common development tasks. Here are the available commands:

### Run all default tasks

To run tests and build the project:
```shell
make all
# Or just run
make
```

### Build the project

To compile the project:
```shell
make build
```

This will create an executable named `deliveroo-take-home` in the project root.

### Run tests

To run all tests:
```shell
make test
```

### Generate test coverage report

To run tests with coverage and view the report in your default web browser:
```shell
make coverage
```

### Clean up

To remove compiled binaries and coverage files:
```shell
make clean
```

## Project Structure

- `main.go`: Contains the main logic for the cron expression parser
- `main_test.go`: Contains unit tests for the parser
- `Makefile`: Defines commands for building, testing, and running the project
