{ pkgs, ... }:
pkgs.buildGoModule {
  pname = "firewall-controller";
  version = "0.0.1";
  src = pkgs.lib.fileset.toSource {
    root=./.;
    fileset = pkgs.lib.fileset.unions [
      ./src
      ./go.mod
      ./go.sum
    ];
  };
  vendorHash = "sha256-sF8RFUEIy3mip/EyJDn0+mRfFbeBbn18rqsWtfsAOqo=";
  # vendorHash = pkgs.lib.fakeHash;
}
