# Services

A Go-based microservices project with a modular structure. This project includes various components to demonstrate the implementation of a scalable microservices architecture.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Folder Structure](#folder-structure)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/dartey25/services.git
    cd services
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Build the project:
    ```bash
    make build
    ```

## Configuration

Configure the application using the `config.yaml` file. Here’s an example configuration:

```yaml
server:
  port: 8080
database:
  host: localhost
  port: 5432
  user: username
  password: password
  dbname: servicename
