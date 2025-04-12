{
  description = "abi-testdata";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    # this tool allows us use nix-shell and nix shell
    # and is used for our shell.nix
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.abi-testdata = pkgs.buildGoModule {
          pname = "abi-testdata";
          version = "0.1.0-prerelease";

          src = ./.;

          vendorHash = "sha256-BwWWvDjOiUFfRqZYnY6auhhhXgi/EaqN5lJRSABKM9g=";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.git
            pkgs.go
            pkgs.gotools
            pkgs.nixfmt-rfc-style
          ];
        };

        defaultPackage = self.packages.${system}.abi-testdata;
        defaultDevShell = self.devShells.${system}.default;
      }
    );
}
