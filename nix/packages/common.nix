# This package is need for testing
{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "autobrowser-common";
  version = "0";
  vendorHash = "sha256-CVycV7wxo7nOHm7qjZKfJrIkNcIApUNzN1mSIIwQN0g=";
  src = import ../src.nix {inherit lib;};

  modRoot = "common";
}
