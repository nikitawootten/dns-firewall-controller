{ pkgs, ... }:
let
  common = import ./buildCommon.nix { inherit pkgs; };
in
pkgs.buildGoModule {
  pname = "firewall-controller";
  version = common.version;
  src = common.src;
  vendorHash = "sha256-sF8RFUEIy3mip/EyJDn0+mRfFbeBbn18rqsWtfsAOqo=";
  # vendorHash = pkgs.lib.fakeHash;
}
