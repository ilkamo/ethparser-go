# Ethereum Parser

I had a lot of fun working on this project. My approach was to create a parser that could be easily extended with new
repositories and ethereum clients. Most of the code is tested with unit tests and I tried to keep it as clean as possible.
The most interesting parts are commented in the code so that reviewers can understand my thought process. The parser processes 
blocks one by one (for simplicity) as this is not a perfect production ready service. In a prod environment I would 
extend it to process multiple blocks at once by adding multiple workers. This way, the caller could control the batch size 
and the parser could be more efficient.

## Packages

I created different packages to have a separation of concerns and I used the `internal` package for the components that I don't want to expose to the public.

- `internal`
  - `e2e` e2e test suite.
  - `ethereum` logic to interact with the needed methods of the ethereum client. It relies on a generic `RPCClient` interface that can be implemented by any client. Inside the package, there are some utility functions to deal with ethereum hex numbers.
  - `jsonrpc` logic to interact with any JSON-RPC server. It is used by my ethereum client.
    Inside, there are some **transport-layer** types. The `HTTPRequestBuilder` is a really simple builder, and it could be
    replaced with a more generic one.
  - `storage` implementation of an in-memory concurrency safe `TransactionsRepository`
    and `AddressesRepository`.
  - `mock` mocks for the tests.
  - `testdata` test data used by the tests. Here I used **go:embed** to simply load a json file with real ethereum block
    data. 


- `parser` logic to parse and observe the ethereum blocks and transactions. The parser accepts
  different options to enhance the default implementation with more sophisticated components.


- `types` types used by the main components.


## Usage

The default parser can be initialized in the following way:

```go
ctx := context.Background()
log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

p, err := parser.NewParser(
  "https://cloudflare-eth.com",
  log,
)
// handle the error

err = p.Run(ctx)
```

The parser can be extended with different components by passing options to the constructor. For example, to use a different repository, you can pass the `WithTransactionsRepository` option:

```go
repo := NewTransactionsRepository()

p, err := parser.NewParser(
  "https://cloudflare-eth.com",
  log,
  parser.WithTransactionsRepository(repo),
)
```

All the available options are defined in the [parser/options.go](parser/options.go) file.


## Testing

All the unit tests can be run with the following command:

```bash
make test
```

To display the tests coverage in an interactive html page, run the following command:

```bash
make display_coverage
```

In addition, I created a simple e2e test suite that runs the parser with a real ethereum client to show that it works as
expected. The test is located in the `tests/e2e` folder. Usually I would run this with a docker-compose file that would start an ethereum emulator, but wanted to keep it simple for this project. To
start the test, just run the following command:

```bash
make test_e2e
```

## Linting

To lint the code, run the following command:

```bash
make lint
```

To fix linting issues, run the following command:

```bash
make fmt
```

## CI

I added a simple GitHub Actions workflow that runs the tests and linter on every pull request.

## Considerations

```go
type Parser interface {
    GetCurrentBlock() int
    Subscribe(address string) bool
    GetTransactions(address string) []Transaction
}
```

During the implementation, I thought about some improvements that could be made to the **public interface** of the **parser**.

- `context` is not being used in the parser methods. I would definitely add it in a real-world scenario to allow the
  caller to cancel the operation and pass trace information.
- The `Subscribe` and `GetTransactions` methods could be improved to return an error. This way, the caller could know if
  the subscription failed and why (really useful it there is a cache or a network error).
- The parser doesn't have a method to start parsing. I assumed that it would start parsing when the parser is created. I
  really don't like this approach this is why I added a Run method to the parser. This way, the caller can start the parser
  and control it with a context.
- The `GetCurrentBlock` returns an int. I would change it to return an **uint64** to avoid possible future overflow issues.
