# DeviceSpoof
The original file was given to me by BoomYouBang#0001, as I can not code GO. For now it uses a fork of GopherTunnel to allow NBT with Non-UTF8 characters. I might add more features in the future. Open an issue for any bugs or any features you would like to see.

# Running
In order to compile your own copy make sure you have GOLang installed (https://golang.org/doc/install). CD to the folder then use go get github.com/sandertv/gophertunnel. After running the command it should download a fork of GopherTunnel (You can confirm this is there is a go.sum file). Once done use go build. This will give you a new exe file to run. To run it from source use go run main.go (Make sure to do the steps before go build!). You will be promped to authenticate a microsoft account at a link. After authenticating, launch minecraft then sign in with the same account. Lastly, add a server with the ip 127.0.0.1 and port 19132, then connect. You should be connected to your desired server.

# Configuring 
To configure, go to main.go. At line 15 you will see a server address to add, you can change this to whatever server you like. Under that you can configure what you would like your DeviceID to be (It can be anything, but ID's that don't look like a real device id may get you banned). 

