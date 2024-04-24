# This package is need for testing
{ lib, stdenv, buildGoModule, makeWrapper, ddcutil }:
buildGoModule rec {
  pname = "autobrowser-common";
  version = "0";
  vendorHash = null;
  src = import ./src.nix { inherit lib; };

  modRoot = "common";
}
