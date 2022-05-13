# hey_grpc

`hey_grpc` is very simple cli tool for sending multiple requests on gRPC endpoints.   
Creating this tool was inspired by amazing [hey](https://github.com/rakyll/hey) project.


## Installation
```shell
go install github.com/mispon/hey_grpc@latest
```

--- *or* ---

1. `sudo chown -R $(whoami) /usr/local/bin` (*optional, if you have not permissions*)
2. `git clone https://github.com/mispon/hey_grpc`
3. `cd hey_grpc`
4. `make build`

--- *or* ---
1. [download](https://github.com/mispon/hey_grpc/releases/download/v0.2.0/hey_grpc_darwin_amd64) pre-compilled binary for MacOS m1
2. put it were you want

## Commands
```
  * ls      - prints services on specified host:port
  * lsm     - prints service's methods
  * call    - process gRPC unary calls
  * version - print current version
  * help    - just help
```

## Call usage
```
Usage: grpc_hey call [options...] <host:port> <Service>/<Method> <Message>
How to pass right args for grpc_cli see in [official documetation](https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md).

Options:
  -n  Number of requests to run. Default is 1.
  -w  Number of workers to run concurrently. Default is 1.
  -d  Duration of sending requests. When duration is reached,
      tool stops and exits. If duration is specified, n is ignored.
      Examples: -d 1h or -d 3m or -d 100s or -d 500ms.
  -t  Timeout for each request in seconds. Default is 0s.
      Examples: -t 1h or -t 2m or -t 10s or -t 500ms.
  -q  Limiter of calls per secons.
```

## Examples
Let's imagine that there is simple `ping.proto` contract for our server:
```protobuf
message PingRequest {
  string s = 1;
}

message PongResponse {
  string s = 1;
}

service PingService {
  rpc Ping(PingRequest) returns (PongResponse) {}
}
```

Then we can call this endpoint 50 times:
```shell
hey_grpc call -n 50 localhost:80 PingService/Ping 's: "ping!"'
```

or 100 times in 10 workers:
```shell
hey_grpc call -n 100 -w 10 localhost:80 PingService/Ping 's: "ping!"'
```

or during 5m in 3 workers with 10 seconds timeout after each call:
```shell
hey_grpc call -d 5m -w 3 -t 10s localhost:80 PingService/Ping 's: "ping!"'
```
