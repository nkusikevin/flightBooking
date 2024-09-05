#!/bin/bash

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <search_directory_base> <working_directory>"
  exit 1
fi

search_directory_base="$1"
working_directory="$2"

max_parallel=8
running_processes=0
error_file=$(mktemp)

build_and_zip() {
  dir="$1"
  function=$(basename "$dir")

  if [ -f "${dir}/main.go" ]; then
    mkdir -p "${working_directory}/bin/$function"
    echo "Building $function"
    (cd "$dir" &&
      {
        env GOARCH=arm64 GOOS=linux go build -o "${working_directory}/bin/$function/bootstrap" *.go &&
        pushd ${working_directory}/bin/$function/ > /dev/null && zip -j -q "${working_directory}/bin/$function.zip" bootstrap && popd > /dev/null
      } > >(tee "${working_directory}/bin/$function/build_output.log") 2>&1) || { echo "Error building $function" >> "$error_file"; cat "${working_directory}/bin/$function/build_output.log" >> "$error_file"; }
  fi
}

cd "$search_directory_base"

for dir in $(find . -type d -not -path '*/node_modules/*'); do
  build_and_zip "$dir" &
  ((running_processes++))

  if [ "$running_processes" -ge "$max_parallel" ]; then
    while [ $(jobs -p | wc -l) -ge "$max_parallel" ]; do
      sleep 0.01
    done
    ((running_processes--))
  fi
done

cd "$working_directory"

wait

if [ -s "$error_file" ]; then
  echo "Errors occurred during the build:"
  cat "$error_file"
else
  echo "Build completed successfully."
fi

rm -rf "$error_file"