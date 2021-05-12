{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  pname = "timewsync-server";
  version = "1.0.0";
  src = ./.;
  vendorSha256 = "0wbd4cpswgbr839sk8qwly8gjq4lqmq448m624akll192mzm9wj7";
}
