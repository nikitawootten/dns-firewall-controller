{ pkgs, ... }:
let
  pname = "coredns";
  version = "1.11.1";
  repo = "github.com/nikitawootten/dns-firewall-controller";
  plugin = "${repo}/src/coredns_plugin";
  plugin-name = "squawker";
  coredns-src = pkgs.fetchFromGitHub {
    owner = "coredns";
    repo = "coredns";
    rev = "v${version}";
    sha256 = "sha256-XZoRN907PXNKV2iMn51H/lt8yPxhPupNfJ49Pymdm9Y=";
  };
  plugin-src = pkgs.lib.fileset.toSource {
    root=./.;
    fileset = pkgs.lib.fileset.unions [
      ./src
      ./go.mod
      ./go.sum
    ];
  };
in
pkgs.buildGoModule {
  inherit pname version;

  src = coredns-src;

  outputs = [ "out" "man" ];
  
  nativeBuildInputs = [ pkgs.installShellFiles ];

  vendorHash = "sha256-OLFgY5O4RO6iizUbRino6eaY9vPakEA8doDI4jiegsY=";
  # vendorHash = pkgs.lib.fakeHash;

  # VERY hacky way to add a plugin to the coredns build
  modBuildPhase = ''
    # Add our plugin to the go.mod file using the replace directive
    go mod edit -replace '${repo}=${plugin-src}'
    go get ${plugin}
    echo "${plugin-name}:${plugin}" >> plugin.cfg

    GOOS= GOARCH= go generate
    go mod vendor
    # Vendoring only copies the relevant files from our source derivation (symlink to Nix store no longer maintained).
    # This is a problem because go.mod and modules.txt still reference the Nix store, and Nix gets very upset at random references to the Nix store
    
    # After vendoring we need to surgically remove all unused references to the Nix store
    go mod edit -dropreplace '${repo}'
    sed -i 's/ => \/nix\/store.*//g' vendor/modules.txt
  '';

  # Verbatim copy of the nixpkgs coredns derivation (https://github.com/NixOS/nixpkgs/blob/nixos-unstable/pkgs/servers/dns/coredns/default.nix)
  modInstallPhase = ''
    mv -t vendor go.mod go.sum plugin.cfg
    cp -r --reflink=auto vendor "$out"
  '';

  preBuild = ''
    chmod -R u+w vendor
    mv -t . vendor/go.{mod,sum} vendor/plugin.cfg

    GOOS= GOARCH= go generate
  '';

  postPatch = ''
    substituteInPlace test/file_cname_proxy_test.go \
      --replace "TestZoneExternalCNAMELookupWithProxy" \
                "SkipZoneExternalCNAMELookupWithProxy"

    substituteInPlace test/readme_test.go \
      --replace "TestReadme" "SkipReadme"

    # this test fails if any external plugins were imported.
    # it's a lint rather than a test of functionality, so it's safe to disable.
    substituteInPlace test/presubmit_test.go \
      --replace "TestImportOrdering" "SkipImportOrdering"
  '' + pkgs.lib.optionalString pkgs.stdenv.isDarwin ''
    # loopback interface is lo0 on macos
    sed -E -i 's/\blo\b/lo0/' plugin/bind/setup_test.go
  '';

  postInstall = ''
    installManPage man/*
  '';

  # Additional checks to ensure that the plugin was properly added to the binary
  postCheck = ''
    # Sanity check: was the plugin included at all?
    $GOPATH/bin/coredns -plugins | grep dns.${plugin-name} || { echo "Plugin not registered in output binary"; exit 1;}

    pushd vendor/${repo}

    # Sanity check all vendored plugin files against the source derivation
    # Currently we must update the vendor hash every time a go file changes
    find . -type f -name '*.go' -print0 | while IFS= read -r -d $'\0' file; do
      vendorSum=$(sha256sum "$file" | cut -d' ' -f1)
      srcSum=$(sha256sum "${plugin-src}/$file" | cut -d' ' -f1)
      if [ "$vendorSum" != "$srcSum" ]; then
        echo "File $file does not match source derivation"
        exit 1
      fi
    done

    popd
  '';
}
