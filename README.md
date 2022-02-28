Pando Blockchain Ledger Protocol
The Pando Blockchain Ledger is a Proof-of-Stake decentralized ledger designed for the video streaming industry. It powers the Pando token economy which incentives end users to share their redundant bandwidth and storage resources, and encourage them to engage more actively with video platforms and content creators. The ledger employs a novel multi-level BFT consensus engine, which supports high transaction throughput, fast block confirmation, and allows mass participation in the consensus process. Off-chain payment support is built directly into the ledger through the resource-oriented micropayment pool, which is designed specifically to achieve the “pay-per-byte” granularity for streaming use cases. Moreover, the ledger storage system leverages the microservice architecture and reference counting based history pruning techniques, and is thus able to adapt to different computing environments, ranging from high-end data center server clusters to commodity PCs and laptops. The ledger also supports Turing-Complete smart contracts, which enables rich user experiences for DApps built on top of the Pando Ledger. 

To learn more about the Pando Network in general, please visit the Pando Documentation site: https://doc.pandoproject.org/pando-network-testament/introduction.

Table of Contents
Setup

Setup
Intall Go
Install Go and set environment variables GOPATH , GOBIN, and PATH. The current code base should compile with Go 1.14.2. On macOS, install Go with the following command

brew install go@1.14.1
brew link go@1.14.1 --force
Build and Install
Next, clone this repo into your $GOPATH. The path should look like this: $GOPATH/src/github.com/pandotoken/pando

git clone https://github.com/logan-pandolab/Pando-Mainnet-Blockchain.git $GOPATH/src/github.com/pandotoken/pando
export PANDO_HOME=$GOPATH/src/github.com/pandotoken/pando
cd $PANDO_HOME
Now, execute the following commands to build the Pando binaries under $GOPATH/bin. Two binaries pando and pandocli are generated. pando can be regarded as the launcher of the Pando Ledger node, and pandocli is a wallet with command line tools to interact with the ledger.

export GO111MODULE=on
make install
Notes for Linux binary compilation
The build and install process on Linux is similar, but note that Ubuntu 18.04.4 LTS / Centos 8 or higher version is required for the compilation.

Notes for Windows binary compilation
The Windows binary can be cross-compiled from macOS. To cross-compile a Windows binary, first make sure mingw64 is installed (brew install mingw-w64) on your macOS. Then you can cross-compile the Windows binary with the following command:

make exe
You'll also need to place three .dll files libgcc_s_seh-1.dll, libstdc++-6.dll, libwinpthread-1.dll under the same folder as pando.exe and pandocli.exe.

Run Unit Tests
Run unit tests with the command below

make test_unit
