"""
Code for building Openshift v3
"""

import sys

from tito.common import (get_latest_commit, get_latest_tagged_version, check_tag_exists,
        run_command, get_script_path, find_spec_file, get_spec_version_and_release)
from tito.builder import Builder

class OpenshiftBuilder(Builder):
  """
  builder which defines 'commit' as the git hash prior to building

  Used For:
    - Packages that want to know the commit in all situations
  """

  def _get_rpmbuild_dir_options(self):
      git_hash = get_latest_commit()
      cmd = '. ./hack/common.sh ; echo $(os::build::ldflags)'
      ldflags = run_command('bash -c \'%s\''  % (cmd) )

      return ('--define "_topdir %s" --define "_sourcedir %s" --define "_builddir %s" '
            '--define "_srcrpmdir %s" --define "_rpmdir %s" --define "ldflags %s" '
            '--define "commit %s" ' % (
                self.rpmbuild_dir,
                self.rpmbuild_sourcedir, self.rpmbuild_builddir,
                self.rpmbuild_basedir, self.rpmbuild_basedir,
                ldflags, git_hash))

  def _get_build_version(self):
      """
      Figure out the git tag and version-release we're building.
      """
      # Determine which package version we should build:
      build_version = None
      if self.build_tag:
          build_version = self.build_tag[len(self.project_name + "-"):]
      else:
          build_version = get_latest_tagged_version(self.project_name)
          if build_version is None:
              if not self.test:
                  error_out(["Unable to lookup latest package info.",
                          "Perhaps you need to tag first?"])
              sys.stderr.write("WARNING: unable to lookup latest package "
                  "tag, building untagged test project\n")
              build_version = get_spec_version_and_release(self.start_dir,
                  find_spec_file(in_dir=self.start_dir))
          self.build_tag = "v%s" % (build_version)

      if not self.test:
          check_tag_exists(self.build_tag, offline=self.offline)
      return build_version
