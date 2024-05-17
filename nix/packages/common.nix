# This package is need for testing
{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "autobrowser-common";
  version = "0";
  vendorHash = null;
  src = import ../src.nix {inherit lib;};

  modRoot = "common";
}
