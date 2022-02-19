# streamer clock - desktop edition
Browser page is easier to share, but a limitation that is hit right away is what features may or may not work across different browsers when a tab is not in focus. This is an effort to rebuild the voiced clock with golang + fyne. 

Binaries for win64 will be available on the releases page in a future build

# building on windows

I used msys to get gcc and built with `go build -ldflags -H=windowsgui .` to not launch a console window

Fyne has CGO dependencies, so GCC and golang are both required to build this repo. 

- get GO from [the golang website](https://golang.org/)
- get GCC with a tool like MSYS, MSYS available [here](https://www.msys2.org/)

# MSYS and GCC
msys may not come with gcc (use `which gcc` to check), so get it with this command if necessary `pacman -S mingw-w64-x86_64-toolchain`. I may update that to a specific set of tools if I ever figure out which pieces of it are needed. If after restarting msys `which gcc` still gives nothing, you can try one of these solutions:

## first option

See if running the 'correct' msys resolves this for you. When you install msys, you get a handful of executables. Try navigating to that folder (probably C:\msys64) and running mingw64.exe directly. Then check `which gcc` again. If you get something like "/mingw64/bin/gcc" perfect, you're done here.

## second option

I'm just perusing old stack overflow questions at this point. One suggestion here is finding the gcc executable and adding it to your path. It's probably the same path as the example path from option 1, so you'll want to add that /mingw64/bin folder to your windows path. Before doing so, verify it's there with a command like `ls /mingw64/bin | grep gcc`. If it's not there, something went wrong with the pacman command. If it is there, great - add this to your windows path, relaunch msys if necessary, and run that 'which' check again. 

## golang in msys

If option #2 for gcc works for you, hopefully something similar will do for adding golang to your path. For whatever reason, that did not work for me. In my case I had to add go by editing the bash rc. This might be the scariest part of this guide - it's time to use vim. First, make sure you're home, use `cd ~`. Then really commit to this life experience and run `vim .bashrc`

Vim being scary is a bit of a meme, but if you're not familiar with it or terminals it certainly looks scary, especially since you can't ctrl+c your way out of it. For now...

1. scroll down to the bottom of the file, using your arrow keys to get the cursor to the end
2. press i - this will allow you to edit the text. Use arrow keys to get to the end of whatever line you're on, then hit enter to get a new line
3. from there, let the terminal know where go is installed, type out these 2 lines (maybe check windows explorer and see where your go files are):
- `export GOBIN=/c/Program\ Files/Go/bin`
- `export PATH=$PATH:$GOBIN`
4. ensure there are no '#' characters at the beginning of those lines and that should be it, press `escape`, then press `:`, then press `x` and hit enter
5. back at the terminal, type `source .bashrc`, following that type `which go` and you should get the path you entered

## verify

- to verify go is installed correctly, type `go version`
- to verify gcc is installed correctly, type `gcc --version`
- ideally, both of those statements will print back something that looks like a version number

If there was no trouble, now you're ready to build! For windows, see the earlier section that has the build command to use. That will build this application without also launching a terminal in the background. If you want the terminal to be there for whatever reason, use `go build .` to build the application. Apologies to the other platforms - I don't know what flags are needed to stop terminals from launching in the background nor do I know if one will, but hey, at least the clock will still work! 

Whatever your desired build may be, I'm also going to mention for the less technical that you need to run this command in the repo you downloaded or cloned. MSYS will start you off in its own little world, but you can get back to your regular system by typing `cd /c/` and then navigate your way to the folder via terminal by using tab to find the next folder! If you know the path to your destination you can go straight there. Let's say for example I download this repository to my desktop - from msys I can use `cd /c/Users/tom/Desktop/crump-clock` to get there, and then build. There may also be a way to add msys to your right click menu to open it directly at the folder for future convenience.

# future note

I'll be looking into github actions once some other projects are further along. That will automate the build instructions you see above and make an executable available to download from the releases page - but it will warn you that the executable is unsigned or unknown. If you manage to build this project yourself, there will be no annoying warning when you launch it. 
