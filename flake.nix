{
  description = "A Prometheus exporter for cheap 2.5G switches";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , gomod2nix
    ,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Simple test check added to nix flake check
        go-test = pkgs.stdenvNoCC.mkDerivation {
          name = "go-test";
          dontBuild = true;
          src = ./.;
          doCheck = true;
          nativeBuildInputs = with pkgs; [
            go
            writableTmpDirAsHomeHook
            (lib.optional stdenv.hostPlatform.isDarwin pkgs.apple_sdk_11_0)
          ];
          checkPhase = ''
            go test -v ./...
          '';
          installPhase = ''
            mkdir "$out"
          '';
        };
        # Simple lint check added to nix flake check
        go-lint = pkgs.stdenvNoCC.mkDerivation {
          name = "go-lint";
          dontBuild = true;
          src = ./.;
          doCheck = true;
          nativeBuildInputs = with pkgs; [
            golangci-lint
            go
            writableTmpDirAsHomeHook
          ];
          checkPhase = ''
            golangci-lint run
          '';
          installPhase = ''
            mkdir "$out"
          '';
        };
      in
      {
        checks = {
          inherit go-test go-lint;
        };
        packages.default = pkgs.callPackage ./. {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
        };
        devShells.default = pkgs.callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        };
      }
    ));
}
