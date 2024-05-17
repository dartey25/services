test12343

# Services

A Go-based microservices project with a modular structure. This project includes various components to demonstrate the implementation of a scalable microservices architecture.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Folder Structure](#folder-structure)

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
```

## Usage

To run the application:

```bash
make run
```

For live reloading during development, use Air:

### Running with Air

Install Air:

```bash
curl -fLo air https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air
chmod +x air
sudo mv air /usr/local/bin
```

Run the application:

```bash
make air
```

## Folder Structure

- `assets`: Static files
- `cmd`: Main applications
- `config`: Configuration files
- `internal`: Private application and library code
- `pkg`: Public library code
- `web`: Front-end code
