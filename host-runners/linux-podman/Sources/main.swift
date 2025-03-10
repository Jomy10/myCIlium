import Foundation

func printUsage() {
  print("USAGE: \(ProcessInfo.processInfo.arguments[0]) [mode]")
  print("modes:")
  print("\tinit\tInitialize the runner by setting up a systemd timer")
  print("\tpoll\tPolls the API to check if a build should be started and attempts to start building")
}

if ProcessInfo.processInfo.arguments.count != 2 {
  printUsage()
  exit(1)
}

let mode = ProcessInfo.processInfo.arguments[1]

switch (mode) {
  case "init":
    runInit()
  case "poll":
    break
  default:
    print("ERROR: invalid mode \(mode)")
    printUsage()
}
