if [ "$CIRCLE_CI" == "true" ]; then
  sudo wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
  sudo mv jq-linux64 /usr/bin/jq
  sudo chmod a+x /usr/bin/jq
  sudo rm -rf jq-linux64
  jq --version
fi

version=`cat .goxc.json | jq -r ".PackageVersion"`
has_version=`curl https://api.github.com/repos/minodisk/qiitactl/tags | jq "\
  any( \
    .name == \"v$version\" \
  ) \
"`
if [ "$has_version" == "true" ]; then
  echo "v$version is already released"
  exit 1
fi

goxc -wlc compile package publish-github -apikey=$GITHUB_TOKEN
goxc

git config --global user.email "daisuke.mino@gmail.com"
git config --global user.name "minodisk"

git clone https://minodisk:$GITHUB_TOKEN@github.com/minodisk/homebrew-qiitactl.git
go run .brew/generate.go
cd homebrew-qiitactl
git commit -m "Update formula" qiitactl.rb
git push
