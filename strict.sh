#!/usr/bin/env bash

baseDirectory="$(bazel info bazel-genfiles)/cmd/strict"
if ! ls "$baseDirectory" &> /dev/null; then
  bazel build //cmd/strict:strict
fi

commandDirectories=$(ls "$baseDirectory")
directoryCount=${#commandDirectories[@]}
if [ "$directoryCount" -eq 0 ]
then
  bazel build //cmd/strict:strict
  commandDirectories=$(ls "$baseDirectory")
fi

for name in $commandDirectories
do
  export binaryOutputDirectory="$baseDirectory/$name"
done

case "$OSTYPE" in
  msys*)
    exec "$binaryOutputDirectory/strict.exe" "$*"
    ;;
  *)
    echo "This strict doesn't currently support your OS."
    ;;
esac