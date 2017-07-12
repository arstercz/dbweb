echo "building ..."
rm -rf ./dbweb
mkdir ./dbweb
GOOS=linux GOARCH=amd64 go build -o ./dbweb/dbweb
cp usercfg.conf ./dbweb/
cp -r ./options ./dbweb/
cp -r ./public ./dbweb/
cp -r ./templates ./dbweb/
cp ./*.pem ./dbweb/
tar zcvf ./dbweb.tar.gz ./dbweb
echo "done."
