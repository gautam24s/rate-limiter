## Rate Limiter Middleware

This project demonstrates a simple HTTP server with rate limiting middleware implemented in Go. The middleware restricts the number of requests based on either IP addresses or endpoints.

### Overview

The project consists of the following files:

- `/middleware`: Contains the implementation of the rate limiting middleware.
- `main.go`: Implements a basic HTTP server using the rate limiting middleware.

### Rate Limiting Rules

The rate limiting middleware supports two types of rules:

1. **IP Rule**: Limits requests based on client IP addresses.
2. **Endpoint Rule**: Limits requests based on specific endpoints.

### How to Use

To use the rate limiting middleware in your project, follow these steps:

1. Define your rate limiting rules in the `main.go` file by creating an instance of `LimitRules`.
2. Create a new instance of the middleware using the defined rules with the `New` function.
3. Define your HTTP request handler and attach the rate limiting middleware to the desired endpoint using `RateLimitingMiddleware`.
4. Start the HTTP server using `http.ListenAndServe`.

Example:

```go
// Define rate limiting rules
limiters := []middleware.LimitRules{
    {
        Limit:    5,
        Window:   10 * time.Second,
        IP:       []string{"0.0.0.0"},
        RuleType: middleware.IPRule,
    },
    {
        Limit:    10,
        Window:   10 * time.Second,
        Endpoint: []string{"/"},
        RuleType: middleware.EndpointRule,
    },
}

// Create a new instance of middleware with the defined rate limiting rules
mw := middleware.New(ctx, limiters)

// Define a handler function
handler := func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Request processed for %s", r.URL.Path)
}

// Attach the rate limiting middleware to the root endpoint
http.HandleFunc("/", mw.RateLimitingMiddleware(handler))

// Start the HTTP server
http.ListenAndServe(":8080", nil)
```

### Testing

To test the rate limiting middleware, run the following command:
```bash
go test ./...
```

To run the sample server, run the following commands:
```bash
go run main.go
```


### License

This project is licensed under the MIT License - see the LICENSE file for details.


### Authors

Gautam Sharma
