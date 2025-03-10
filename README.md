<div align="center">
  <h1 style="font-weight: normal">MY<b>CI</b>LIUM</h1>
  <p>
    distribute tasks to host machines for CI purposes
  </p>
  <!--
  ❰
  <a href="/tests">examples</a>
  |
  <a href="/main.go">endpoints</a>
  ❱
  -->
</div><br/>

## How it works

1. An event triggers triggers an action that instructs the orchestrator to make a build request for specific platforms. e.g. A push to a branch
   ```mermaid
   graph LR

   GH[Push to repo] --Build request--> Orchestrator[Orchestrator Server]
   ```

2. Host machines checks with the orchestrator server to determine if they should start a new build. When there is a new repo to build, then the host machine will confirm that it will start building
   ```mermaid
   sequenceDiagram

   macOS Host Machine -->> Orchestrator Server: poll
   Orchestrator Server -->> macOS Host Machine: list of things to build for macOS
   macOS Host Machine -->> Orchestrator Server: I will start building
   Orchestrator Server -->> macOS Host Machine: ok, start building
   ```

3. The host machine clones the repo and starts doing the steps specified in the configuration.

4. When finished, it notifies the orchestrator and will do any post-finish steps specified (like uploading artifacts)

   ```mermaid
   sequenceDiagram

   macOS Host Machine -->> Orchestrator Server: I have finished
   Orchestrator Server -->> macOS Host Machine: ok
   macOS Host Machine -->> GitHub: e.g. upload artifacts to latest release
   ```

## In this repository

- [server: server implementation "myCIlium orchestrator"](/server)
- [spec: specification of the yaml format to be used for specifying what a host machine must do](/spec)
