# urlshortener

This project provides an API for shortening the URL. It uses the standard http router, Go's standard library, and Docker for easy containerization.

## Prerequisites

- Docker installed on your system.
- Go 1.23+ installed on your local machine (if you want to run the code locally).

## Getting Started

1. **Clone the repository**:

    ```bash
    git clone https://github.com/rishabh625/urlshortener.git
    cd urlshortener
    ```

2. **Running with Docker**:

    - Build and run the Docker container:

        ```bash
        docker build -t urlshortener-service .
        docker run -p 8080:8080 urlshortener-service
        ```

    - The server will be available on `http://localhost:8080`.

3. **Running locally** (without Docker):

    - Install dependencies:

        ```bash
        go mod vendor
        go mod tidy
        ```

    - Run the application:

        ```bash
        go run cmd/main.go 
        ```

    - The server will be available on `http://localhost:8080`.

## API Endpoints

### `POST /shortURL`
Shorten's the provided longURL

#### Request Body:

```json
{
  "longURL":"https://zh.wikipedia.org/wiki/%E7%99%BE%E5%BA%A6"
}
```
### Curl Call

Shorten's the provided longURL
```
curl --location 'http://localhost:8080/shortURL' \
--header 'Content-Type: application/json' \
--data '{
    "longURL":"https://zh.wikipedia.org/wiki/%E7%99%BE%E5%BA%A6"
}'
```


### GET `/{id}`
Redirects User to original URL if id is present in DB and is shortened using /shortURL API call
### Curl Call
```
curl --location --request GET 'http://localhost:8080/YFGmAXAvns4D4C2dO6FtJXqw8xHgfltm' \
--header 'Content-Type: application/json' \
--data '{
    "longURL":"https://www.baeldung.com/cs/redirection-status-codes"
}'
```

### GET `/metrics`

This endpoint returns top 3 domain , for which shorten URL service was used

### Curl Call
```
curl --location --request GET 'http://localhost:8080/metrics' \
--header 'Content-Type: application/json' \
--data '{
    "longURL":"https://www.baeldung.com/cs/redirection-status-codes"
}'
```
