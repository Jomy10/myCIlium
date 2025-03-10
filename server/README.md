<div align="center">
  <h1 style="font-weight: normal">MY<b>CI</b>LIUM ORCHESTRATOR</h1>
  <p>
    organizes tasks for other computers to pick up
  </p>
  ❰
  <a href="/server/tests">examples</a>
  |
  <a href="/server/main.go">endpoints</a>
  ❱
</div><br/>

## Building

```sh
go build -o mycilium-orchestrator -ldflags "-s -w"
```

## Running

```sh
mycilium-orchestrator [PORT]
```

## Testing

Testing of the server is done using a ruby script

```sh
bundler install
ruby test.rb
```
