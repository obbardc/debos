/*
Apt Action

Install packages and their dependencies to the target rootfs with 'apt'.

Yaml syntax:
 - action: apt
   recommends: bool
   unauthenticated: bool
   packages:
     - package1
     - package2

Mandatory properties:

- packages -- list of packages to install

Optional properties:

- recommends -- boolean indicating if suggested packages will be installed

- unauthenticated -- boolean indicating if unauthenticated packages can be installed
*/
package actions

import (
	"github.com/go-debos/debos"

	"fmt"
	"net/url"
)

type AptAction struct {
	debos.BaseAction `yaml:",inline"`
	Recommends       bool
	Unauthenticated  bool
	Packages         []string
}

func (apt *AptAction) Run(context *debos.DebosContext) error {
	apt.LogStart()

	aptOptions := []string{"apt-get", "-y"}

	if !apt.Recommends {
		aptOptions = append(aptOptions, "--no-install-recommends")
	}

	if apt.Unauthenticated {
		aptOptions = append(aptOptions, "--allow-unauthenticated")
	}

	aptOptions = append(aptOptions, "install")

	// aptOptions = append(aptOptions, apt.Packages...)

/*

use like:-

  - action: apt
    packages:
      # install packages like normal, so to not break current recipes:
      - mixxx
      - xwax

      # Install packages from URL:
      - http://ftp.us.debian.org/debian/pool/main/b/bmap-tools/bmap-tools_3.5-2_all.deb
      - https://ftp.us.debian.org/debian/pool/main/b/bmap-tools/bmap-tools_3.5-2_all.deb
      - ftp://ftp.debian.org/debian/pool/main/r/rauc/rauc_1.3-1_amd64.deb

      # install packages from file (FIRST CHOICE):
      - file://origin/recipe/packages/test.deb           # installs "packages/test.deb" from "recipe" origin
      - file://origin/filesystem/test.deb    # installs "test.deb" from "filesystem" origin

      # install packages from file (SECOND CHOICE):
      - file://packages/test.deb                 # installs "packages/test.deb" from "recipe" origin
      - origin://recipe/packages/test.deb    # installs "packages/test.deb" from "recipe" origin
      - origin://filesystem/test.deb                  # Installs "test.deb" from "filesystem" origin

*/

	// create list of packages to install by parsing the URI of each
	for _, pkg := range apt.Packages {
		uri, err := url.Parse(pkg)
		if err != nil {
			return err
		}

		// lovely debugging
		fmt.Printf("APT package\n")
		fmt.Printf("\tpkg='%s'\n", pkg)
		fmt.Printf("\tisabs=%t\n", uri.IsAbs())
		fmt.Printf("\tscheme=%s\n", uri.Scheme)
		fmt.Printf("\thost=%s\n", uri.Host)
		fmt.Printf("\trequest uri=%s\n", uri.RequestURI())

		// pkg is a package name
		if !uri.IsAbs() {
			aptOptions = append(aptOptions, pkg)
			continue
		}

		// http, https
		// TODO attempt to support ftp ?
		if uri.Scheme == "http" || uri.Scheme == "https" {
			// TODO download file
			fmt.Printf("\tDownload package over HTTP/HTTPS\n")
		} else if uri.Scheme == "ftp" {
			fmt.Printf("\tDownload package from FTP")
		} else if uri.Scheme == "file" {
			fmt.Printf("\tFILE\n")
		} else {
			return fmt.Errorf("Package URI scheme %s not supported", uri.Scheme)
		}
	}

	return nil


	c := debos.NewChrootCommandForContext(*context)
	c.AddEnv("DEBIAN_FRONTEND=noninteractive")

	err := c.Run("apt", "apt-get", "update")
	if err != nil {
		return err
	}
	err = c.Run("apt", aptOptions...)
	if err != nil {
		return err
	}
	err = c.Run("apt", "apt-get", "clean")
	if err != nil {
		return err
	}

	return nil
}
