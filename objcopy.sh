#!/bin/bash
set -uo pipefail

src="${1:-}"
dst="${2:-}"

if [[ -z "${src}" ]] || [[ ! -f "${src}" ]] || [[ ! -d "${dst}" ]]; then
  echo "Usage $0 [object file] [destination]"
  echo "Script to copy dynamically linked windows binary with their import dependencies"
  exit 1
fi

deps=( )

stack=( "${src}" )
analyzed=( )

ignore=(
  "KERNEL32.dll"
  "msvcrt.dll"
  "ADVAPI32.dll"
  "SHELL32.dll"
  "USER32.dll"
  "WS2_32.dll"
  "ole32.dll"
  "GDI32.dll"
  "MSIMG32.dll"
  "DNSAPI.dll"
  "IPHLPAPI.DLL"
  "SHLWAPI.dll"
  "USERENV.dll"
  "RPCRT4.dll"
  "USP10.dll"
  "gdiplus.dll"
  "SETUPAPI.dll"
  "WINMM.dll"
  "MFPlat.DLL"
  "bcrypt.dll"
  "CRYPT32.dll"
  "ncrypt.dll"
  "WSOCK32.dll"
  "CFGMGR32.dll"
  "DWrite.dll"
)

paths=(
  "."
  "/mingw64/bin/"
)

while [[ ${#stack[@]} -gt 0 ]]; do
  dep="${stack[0]}"
  analyzed+=( "${dep}" )
  
  found=false
  for p in "${paths[@]}"; do
    if [[ -f "${p}/${dep}" ]]; then
      dep="$( readlink -f "${p}/${dep}" )"
      found=true
      break
    fi
  done

  if ! $found; then
    echo "Error: Could not locate dependency \"${dep}\"."
    exit 2
  fi

  echo "Analyze \"${dep}\""
  while IFS='' read -r line; do stack+=("${line}"); done < \
    <( comm -13 \
      <( printf "%s\n" "${ignore[@]}" "${analyzed[@]}" "${stack[@]}" | uniq | sort  ) \
      <( objdump --private-headers "${dep}" \
        | sed -n 's/^\s*DLL Name: \(.*\.dll\)$/\1/Ip' \
        | uniq | sort ) )

  deps+=( "${dep}" )
  stack=( "${stack[@]:1}" )
done

cp -uv "${deps[@]}" "${dst}"

echo -e "\nCopied ${#deps[@]} dependencies"
