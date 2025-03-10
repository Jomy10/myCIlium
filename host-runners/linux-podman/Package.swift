// swift-tools-version: 6.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
  name: "linux-podman",
  platforms: [.macOS(.v10_15)],
  dependencies: [
    .package(url: "https://github.com/vapor/console-kit", from: "4.15.2"),
  ],
  targets: [
    .executableTarget(
      name: "linux-podman",
      dependencies: [
        .product(name: "ConsoleKit", package: "console-kit")
      ]
    ),
  ]
)
