#!/bin/bash
base_branch="main"
version=""
version_flag=0
major=0
minor=0
patch=0

for arg in "$@"
do
  case $arg in
    --version)
      shift
      version=$1
      version_flag=1
      shift
      ;;
    --major)
      major=1
      shift
      ;;
    --minor)
      minor=1
      shift
      ;;
    --patch)
      patch=1
      shift
      ;;
  esac
done

if [[ -z "$version" ]] && [[ $major -eq 0 ]] && [[ $minor -eq 0 ]] && [[ $patch -eq 0 ]]; then
  echo "Either specify version or kind of version increment (major, minor or patch)."
  exit 1
fi

if [[ $major -eq 1 ]] && [[ $minor -eq 1 ]] && [[ $patch -eq 1 ]]; then
  echo "Either specifiy major, minor or patch."
  exit 1
fi

if [[ -z "$version" ]]; then
  version=$(git describe --abbrev=0 | sed -e "s/^v//g")
fi

if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "$version is not valid semver."
  exit 1
fi

if [[ $(git branch --show-current) != "$base_branch" ]]; then
  echo "Must be on $base_branch branch."
  exit 1
fi

ver=""
if [[ $version_flag -eq 1 ]]; then
  ver=$version
else
  semver=($(echo $version | sed -e "s/\./ /g"))
  if [[ $major -eq 1 ]]; then
    semver[0]=$((sevmer[1]+1))
    semver[1]=0
    semver[2]=0
  elif [[ $minor -eq 1 ]]; then
    semver[1]=$((semver[1]+1))
    semver[2]=0
  elif [[ $patch -eq 1 ]]; then
    semver[2]=$((semver[2]+1))
  else
    echo "Could not determine version."
    exit 1
  fi
  ver=$(echo ${semver[@]} | sed -e "s/ /./g")
fi


tag=v$ver
echo "Creating tag: $tag for version: $ver."
echo ""

echo "Pulling from $base_branch branch..."
git pull
echo ""

echo "Running tests..."
go test ./...
echo ""

echo "Creating and pushing tag..."
git tag -a $tag -m "Version $version"
git push origin $tag
