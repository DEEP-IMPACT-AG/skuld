{
  description =
    "CLI tool for AWS Multi-Factor Authentication ";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        lib = pkgs.lib;
      in rec {
        packages = { skuld = pkgs.callPackage ./default.nix { }; };
        defaultPackage = packages.skuld;
        apps.skuld =
          flake-utils.lib.mkApp { drv = packages.skuld; };
        defaultApp = apps.skuld;
      });
}
