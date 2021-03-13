{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoPackage rec {
  pname = "timewsync-server";
  version = "0.0.1";
  goPackagePath = "git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server";
  src = ./.;
  goDeps = ./deps.nix;
}
