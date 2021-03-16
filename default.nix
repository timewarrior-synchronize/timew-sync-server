{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoPackage rec {
  pname = "timewsync-server";
  version = "1.0.0";
  goPackagePath = "github.com/timewarrior-synchronize/timew-sync-server";
  src = ./.;
  goDeps = ./deps.nix;
}
