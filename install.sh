#!/bin/sh
set -e

# Code generated by godownloader on 2018-06-07T07:36:49Z.
# Manually modified to download from S3 for confluentinc/cli goreleaser config
#
# The major modifications include:
# * reworked `github_release` into `s3_release` function
# * updated `TARBALL_URL` and `CHECKSUM_URL` to point to S3 instead of GitHub API
# * added a new `-l` flag to list versions from S3, since we can't link to our (private) GitHub repo
# * extracted a `BINARY` variable instead of having binary names hardcoded in `execute`
# * updated version/tag handling of the `v` prefix; it's expected in GitHub and inconsistently used in S3
# * updated the usage message, logging, and file comments a bit

S3_URL=https://s3-us-west-2.amazonaws.com/confluent.cloud
S3_CONTENT_CHECK_URL="${S3_URL}?prefix="
if [ -n "$OVERRIDE_S3_FOLDER" ]
then
	S3_CONTENT_CHECK_URL="${S3_URL}?prefix=${OVERRIDE_S3_FOLDER}/"
	S3_URL=${S3_URL}/${OVERRIDE_S3_FOLDER}
fi

usage() {
  this=$1
  cat <<EOF
$this: download binaries for confluentinc/cli

Usage: $this [-b] bindir [-d] [tag] | -l
  -b sets bindir or installation directory, Defaults to ./bin
  -d turns on debug logging
   [tag] is a valid version tag shown with -l
   If tag is missing, then the latest will be used.

  -l returns a list of all available tags/versions

EOF
  exit 2
}

parse_args() {
  #BINDIR is ./bin unless set be ENV
  # over-ridden by flag below

  BINDIR=${BINDIR:-./bin}
  while getopts "b:ldh?" arg; do
    case "$arg" in
      b) BINDIR="$OPTARG" ;;
      l) s3_releases ; exit 0 ;;
      d) log_set_priority 10 ;;
      h | \?) usage "$0" ;;
    esac
  done
  shift $((OPTIND - 1))
  TAG=$1
}
# this function wraps all the destructive operations
# if a curl|bash cuts off the end of the script due to
# network, either nothing will happen or will syntax error
# out preventing half-done work
execute() {
  tmpdir=$(mktmpdir)
  log_debug "downloading files into ${tmpdir}"
  http_download "${tmpdir}/${TARBALL}" "${TARBALL_URL}" "Accept:application/octet-stream"
  http_download "${tmpdir}/${CHECKSUM}" "${CHECKSUM_URL}" "Accept:application/octet-stream"
  hash_sha256_verify "${tmpdir}/${TARBALL}" "${tmpdir}/${CHECKSUM}"
  srcdir="${tmpdir}/${BINARY}"
  rm -rf "${srcdir}"
  (cd "${tmpdir}" && untar "${TARBALL}")
  install -d "${BINDIR}"
  for binexe in "${BINARY}" ; do
    if [ "$OS" = "windows" ]; then
      binexe="${binexe}.exe"
    fi
    install "${srcdir}/${binexe}" "${BINDIR}/"
    log_info "NOTICE: see licenses located in ${tmpdir}/${BINARY}"
    log_info "installed ${BINDIR}/${binexe}"
    log_info "please ensure ${BINDIR} is in your PATH"
  done
}
is_supported_platform() {
  platform=$1
  found=1
  case "$platform" in
    alpine/amd64) found=0 ;;
    linux/amd64) found=0 ;;
    darwin/amd64) found=0 ;;
    darwin/arm64) found=0 ;;
    windows/amd64) found=0 ;;
  esac
  case "$platform" in
    alpine/386) found=1 ;;
    linux/386) found=1 ;;
    darwin/386) found=1 ;;
    windows/386) found=1 ;;
  esac
  return $found
}
check_platform() {
  if is_supported_platform "$PLATFORM"; then
    # optional logging goes here
    true
  else
    log_crit "platform $PLATFORM is not supported.  Please contact Confluent support if you believe this is a mistake."
    exit 1
  fi
}
tag_to_version() {
  if [ -z "${TAG}" ]; then
    log_info "checking S3 for latest tag"
  else
    log_info "checking S3 for tag '${TAG}'"
  fi
  REALTAG=$(s3_release "${TAG}") && true
  if test -z "$REALTAG"; then
    log_crit "unable to find '${TAG}' - use 'latest' or see https://docs.confluent.io/${PROJECT_NAME}/current/release-notes.html for avaialble versions."
    exit 1
  fi
  # if version starts with 'v', don't remove it
  TAG="$REALTAG"
  VERSION=${TAG}
}
adjust_format() {
  # change format (tar.gz or zip) based on ARCH
  case ${ARCH} in
    windows) FORMAT=zip ;;
  esac
  true
}
adjust_os() {
  # adjust archive name based on OS
  case ${OS} in
    amd64) OS=x86_64 ;;
    darwin) OS=darwin ;;
    linux) OS=linux ;;
    alpine) OS=alpine ;;
    windows) OS=windows ;;
  esac
  true
}
s3_releases() {
  s3url="${S3_CONTENT_CHECK_URL}${PROJECT_NAME}/archives/&delimiter=/"
  xml=$(http_copy "$s3url")
  versions=$(echo "$xml" | sed -n 's/</\
</gp' | sed -n "s/<Prefix>${PROJECT_NAME}\/archives\/\(.*\)\//\1/p") || return 1
  test -z "$versions" && return 1
  echo "$versions" | sort --version-sort
}
s3_release() {
  version=$1
  test -z "$version" && version="latest"
  s3url="${S3_CONTENT_CHECK_URL}${PROJECT_NAME}/archives/${version#v}/&delimiter=/"
  xml=$(http_copy "$s3url")
  exists=$(echo "$xml" | grep "<Key>") || return 1
  test -z "$version" && return 1
  echo "$version"
}

cat /dev/null <<EOF
------------------------------------------------------------------------
https://github.com/client9/shlib - portable posix shell functions
Public domain - http://unlicense.org
https://github.com/client9/shlib/blob/master/LICENSE.md
but credit (and pull requests) appreciated.
------------------------------------------------------------------------
EOF
is_command() {
  command -v "$1" >/dev/null
}
echoerr() {
  echo "$@" 1>&2
}
log_prefix() {
  echo "$0"
}
_logp=6
log_set_priority() {
  _logp="$1"
}
log_priority() {
  if test -z "$1"; then
    echo "$_logp"
    return
  fi
  [ "$1" -le "$_logp" ]
}
log_tag() {
  case $1 in
    0) echo "emerg" ;;
    1) echo "alert" ;;
    2) echo "crit" ;;
    3) echo "err" ;;
    4) echo "warning" ;;
    5) echo "notice" ;;
    6) echo "info" ;;
    7) echo "debug" ;;
    *) echo "$1" ;;
  esac
}
log_debug() {
  log_priority 7 || return 0
  echoerr "$(log_prefix)" "$(log_tag 7)" "$@"
}
log_info() {
  log_priority 6 || return 0
  echoerr "$(log_prefix)" "$(log_tag 6)" "$@"
}
log_err() {
  log_priority 3 || return 0
  echoerr "$(log_prefix)" "$(log_tag 3)" "$@"
}
log_crit() {
  log_priority 2 || return 0
  echoerr "$(log_prefix)" "$(log_tag 2)" "$@"
}
uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  osid=$(awk -F= '$1=="ID" { print $2 ;}' /etc/os-release 2>/dev/null || true)
  case "$os" in
    msys*) os="windows" ;;
    mingw*) os="windows" ;;
    cygwin*) os="windows" ;;
  esac
  case "$osid" in
    alpine*) os="linux" ;;
  esac
  echo "$os"
}
uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    arm64) arch="arm64" ;;
    aarch64) arch="arm64" ;;
    armv5*) arch="armv5" ;;
    armv6*) arch="armv6" ;;
    armv7*) arch="armv7" ;;
  esac
  echo ${arch}
}
uname_os_check() {
  os=$(uname_os)
  case "$os" in
    darwin) return 0 ;;
    dragonfly) return 0 ;;
    freebsd) return 0 ;;
    linux) return 0 ;;
    alpine) return 0 ;;
    android) return 0 ;;
    nacl) return 0 ;;
    netbsd) return 0 ;;
    openbsd) return 0 ;;
    plan9) return 0 ;;
    solaris) return 0 ;;
    windows) return 0 ;;
  esac
  log_crit "uname_os_check '$(uname -s)' got converted to '$os' which is not a GOOS value. Please file bug at https://github.com/client9/shlib"
  return 1
}
uname_arch_check() {
  arch=$(uname_arch)
  case "$arch" in
    amd64) return 0 ;;
    arm64) return 0 ;;
    armv5) return 0 ;;
    armv6) return 0 ;;
    armv7) return 0 ;;
    ppc64) return 0 ;;
    ppc64le) return 0 ;;
    mips) return 0 ;;
    mipsle) return 0 ;;
    mips64) return 0 ;;
    mips64le) return 0 ;;
    s390x) return 0 ;;
    amd64p32) return 0 ;;
  esac
  log_crit "uname_arch_check '$(uname -m)' got converted to '$arch' which is not a GOARCH value.  Please file bug report at https://github.com/client9/shlib"
  return 1
}
untar() {
  tarball=$1
  case "${tarball}" in
    *.tar.gz | *.tgz) tar -xzf "${tarball}" ;;
    *.tar) tar -xf "${tarball}" ;;
    *.zip) unzip "${tarball}" ;;
    *)
      log_err "untar unknown archive format for ${tarball}"
      return 1
      ;;
  esac
}
mktmpdir() {
  test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
  mkdir -p "${TMPDIR}"
  echo "${TMPDIR}"
}
http_download_curl() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    code=$(curl -n -w '%{http_code}' -sL -o "$local_file" "$source_url")
  else
    code=$(curl -n -w '%{http_code}' -sL -H "$header" -o "$local_file" "$source_url")
  fi
  if [ "$code" != "200" ]; then
    log_debug "http_download_curl received HTTP status $code"
    return 1
  fi
  return 0
}
http_download_wget() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    wget -q -O "$local_file" "$source_url"
  else
    wget -q --header "$header" -O "$local_file" "$source_url"
  fi
}
http_download() {
  log_debug "http_download $2"
  if is_command curl; then
    http_download_curl "$@"
    return
  elif is_command wget; then
    http_download_wget "$@"
    return
  fi
  log_crit "http_download unable to find wget or curl"
  return 1
}
http_copy() {
  tmp=$(mktemp)
  http_download "${tmp}" "$1" "$2" || return 1
  body=$(cat "$tmp")
  rm -f "${tmp}"
  echo "$body"
}
hash_sha256() {
  TARGET=${1:-/dev/stdin}
  if is_command gsha256sum; then
    hash=$(gsha256sum "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command sha256sum; then
    hash=$(sha256sum "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command shasum; then
    hash=$(shasum -a 256 "$TARGET" 2>/dev/null) || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command openssl; then
    hash=$(openssl -dst openssl dgst -sha256 "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f a
  else
    log_crit "hash_sha256 unable to find command to compute sha-256 hash"
    return 1
  fi
}
hash_sha256_verify() {
  TARGET=$1
  checksums=$2
  if [ -z "$checksums" ]; then
    log_err "hash_sha256_verify checksum file not specified in arg2"
    return 1
  fi
  BASENAME=${TARGET##*/}
  want=$(grep "${BASENAME}" "${checksums}" 2>/dev/null | tr '\t' ' ' | cut -d ' ' -f 1)
  if [ -z "$want" ]; then
    log_err "hash_sha256_verify unable to find checksum for '${TARGET}' in '${checksums}'"
    return 1
  fi
  got=$(hash_sha256 "$TARGET")
  if [ "$want" != "$got" ]; then
    log_err "hash_sha256_verify checksum for '$TARGET' did not verify ${want} vs $got"
    return 1
  fi
}
cat /dev/null <<EOF
------------------------------------------------------------------------
End of functions from https://github.com/client9/shlib
------------------------------------------------------------------------
EOF

PROJECT_NAME="confluent-cli"
OWNER=confluentinc
REPO="cli"
BINARY=confluent
FORMAT=tar.gz
OS=$(uname_os)
ARCH=$(uname_arch)
PREFIX="$OWNER/$REPO"

# use in logging routines
log_prefix() {
	echo "$PREFIX"
}
PLATFORM="${OS}/${ARCH}"

main() {
  uname_os_check
  uname_arch_check

  parse_args "$@"

  check_platform

  tag_to_version

  adjust_format

  adjust_os

  if [ ${OS} = "windows" ]; then
    FORMAT=zip ;
  fi

  log_info "found version: ${VERSION} for ${TAG}/${OS}/${ARCH}"

  S3_ARCHIVES_URL=${S3_URL}/${PROJECT_NAME}/archives/${VERSION#v}
  NAME=${BINARY}_${VERSION}_${OS}_${ARCH}
  TARBALL=${NAME}.${FORMAT}
  TARBALL_URL=${S3_ARCHIVES_URL}/${TARBALL}
  CHECKSUM=${BINARY}_${VERSION}_checksums.txt
  CHECKSUM_URL=${S3_ARCHIVES_URL}/${CHECKSUM}

  execute
}

if [ "${TEST}" != "true" ]; then
  main "$@"
fi
