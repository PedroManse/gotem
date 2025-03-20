{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShellNoCC {
  name = "dev-shell";
  packages = with pkgs; [ go ];
  env = {
    GOEXPERIMENT="aliastypeparams";
  };
}

