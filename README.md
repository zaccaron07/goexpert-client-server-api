# Client-Server Exchange Rate

This project consists of two Go applications: `client.go` and `server.go`. The `server.go` application fetches the exchange rate between USD and BRL from an external API and stores it in a SQLite database. The `client.go` application requests the current exchange rate from the server and writes it to a file.

## How to Run

- Ensure all dependencies are installed:
    ```sh
    go mod tidy
    ```

### Server

1. Navigate to the `server` directory:
    ```sh
    cd server
    ```

2. Run the server:
    ```sh
    go run server.go
    ```

   The server will start on port 8080 and expose an endpoint `/cotacao`.

### Client

1. Navigate to the `client` directory:
    ```sh
    cd client
    ```

2. Run the client:
    ```sh
    go run client.go
    ```

   The client will request the current exchange rate from the server and write it to `cotacao.txt` in the format `Dólar: {valor}`.
