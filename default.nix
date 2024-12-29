{ pkgs, ... }:
pkgs.buildGoModule {
  pname = "dns-firewall-controller";
  version = "0.1.0";
  src = ./.;
  vendorHash = "sha256-moCBoEjkhGE1UgGb9Pk894RgxGMImZXJ4u9rMYNtWzY=";
}
