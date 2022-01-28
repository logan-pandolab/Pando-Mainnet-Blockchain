run blockchain node
git clone https://github.com/chandanoodles/pandoblockchaincode.git usr/local/go/src/github.com/pandotoken/pando
export PANDO_HOME=/usr/local/go/src/github.com/pandotoken/pando
#sudo apt-get install build-essential
#sudo snap install go --classic
cd $PANDO_HOME
export GO111MODULE=on

make install

cd $PANDO_HOME

cp -r ./integration/pandonet ../pandonet
mkdir ~/.pandocli
cp -r ./integration/pandonet/pandocli/* ~/.pandocli/

Sudo chmod 700 ~/.pandocli/keys/encrypted

Sudo pando start --config=../pandonet/node
