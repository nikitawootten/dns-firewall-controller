{
  description = "dns-firewall-controller";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    microvm.url = "github:astro/microvm.nix";
    microvm.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, microvm }: let
    forEachSystem = f: nixpkgs.lib.genAttrs ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"] f;
  in {
    packages = forEachSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      in {
        coredns = import ./buildCoreDnsWithPlugin.nix { inherit pkgs; };
        firewall-controller = import ./buildFirewallController.nix { inherit pkgs; };
      });
    devShells = forEachSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      default = pkgs.mkShell {
        packages = [
          microvm.packages.${system}.microvm

          pkgs.go
          pkgs.gopls
          pkgs.delve

          pkgs.protobuf
          pkgs.protoc-gen-go
          pkgs.protoc-gen-go-grpc
        ];
      };
    });
  };
}
