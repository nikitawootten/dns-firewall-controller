{ pkgs, ... }:
let
  gomod = builtins.readFile ../go.mod;
  gomod-lines =  pkgs.lib.strings.splitString "\n" gomod;
in
{
  version = "0.0.1";
  src = pkgs.lib.fileset.toSource {
    root=../.;
    fileset = pkgs.lib.fileset.unions [
      ../src
      ../go.mod
      ../go.sum
    ];
  };
  # Extract repository name from go.mod
  repo = let
    module-line = pkgs.lib.lists.findFirst
      (line: pkgs.lib.strings.hasPrefix "module " line)
      null
      gomod-lines;
  in pkgs.lib.lists.last
    (pkgs.lib.strings.splitString " " module-line);
  # Extract coredns version from go.mod
  coredns-version = let
    coredns-line = pkgs.lib.lists.findFirst
      (line: pkgs.lib.strings.hasInfix "github.com/coredns/coredns " line)
      null
      gomod-lines;
    # github.com/... vX.X.X -> vX.X.X
    raw-version = pkgs.lib.lists.last
      (pkgs.lib.strings.splitString " " coredns-line);
  in if (pkgs.lib.strings.hasPrefix "v" raw-version)
    then (builtins.substring 1 (-1) raw-version)
    else raw-version;
}
