# GREQ

### About

Greq is a super simple http request package that utilizes builder pattern and
go generics to simplify making http requests

### Examples

Send a request to an endpoint and unmarshall the response

```go
type Response struct {
    Status string
    UserID int
}
resp, err := greq.NewRequest[Response]().
    MustPost(endpoint).
    BaseType()

// handle err

fmt.Println(resp.Status, resp.UserID)
```
