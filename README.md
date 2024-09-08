# Deliveroo Take-Home Project

This project is a Go-based implementation of a cron expression parser. It expands cron expressions into a list of execution times. With `make coverage` you can see that there is 94.9% test coverage. You can use the command line tool like this:
```shell
make 
```

## Getting Started

### Prerequisites

- Go 1.16 or later

### Installation

1. Clone the repository:
   ```shell
   git clone https://github.com/dantheman0207/deliveroo-take-home.git
   cd deliveroo-take-home
   ```

2. Install dependencies:
   ```shell
   make deps
   ```

## Usage

This project includes a Makefile to simplify common development tasks. Here are the available commands:


### Run all default tasks

To run tests and build the project:
```shell
make all
```

### Build the project

To compile the project:
```shell
make build
```

This will create an executable named `deliveroo-take-home` in the project root.

### Run the project

To build and run the project:
```shell
make run
```

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
