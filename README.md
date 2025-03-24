# SPIN Kit
**S**andbox **P**rovisioning and **IN**spection Kit

## Features
1. Configure regularly used sandboxes for easy access.
2. List out local connections
3. Refresh sandboxes (based on the configured boxes in #1 or type in a unique sandbox name with the `ref` command)
4. Reconnect to a sandbox (if you've refreshed and need to re-authenticate for development)

# Setup
This is currently only availalbe for `bash` terminals. If I get enough requests, I may consider creating one for PowerShell for Windows. 

## Setup 1 (quick-start)
1. Clone this repo to a directory of your choice. 
2. You are Done!
3. Open your terminal > navigate to where you cloned the repo to > run `./spin-kit.sh`
   1. You may need to add execution permissions for this. You can run: `chmod +x ./spin-kit.sh`
4. This will launch the `spin-kit` app in your terminal. 

## Setup 2 (longterm)
Mint/Ubuntu setups below. The same idea exists though for whatever flavor. 
Including it in your PATH.  Personally, I have a number of scripts that I like to run from the terminal and have a specific scripts drirectory setup that is in my path. 
1. Clone the repo to a directory of your choice.
2. Move the `.sh` file to your directory of choice.  I personally use a directory at `~/scripts`
  1. navigate to the cloned repo directory
  2. if you don't have a scripts directory yet and want one: `mkdir -p $HOME/scripts`
  3. `sudo mv spin-kit.sh $HOME/scripts`
  4. `chmod +x ~/scripts/myscript.sh`
  5. at this point, I actually go in and remove the `.sh` extension for the script.
3. need to add it to the PATH now. Below is how I did it, YMMV.
  1. edit `.bashrc` 
  2. I  have `export PATH=$PATH:x/y/x:a/b/c` setup at the bottom of my `.bashrc` file. Add the route to your `scripts` directory here.
  3. `export PATH=$PATH:$HOME/scripts`
4. at this point, you should be able to open a new terminal window and run `spin-kit` and have the utility launch.

# Configuration
1. Edit the .sh file (you need to edit the one that you're going to use.  if you moved it from the repo dir to the scripts dir, you need to update the scripts dir version)
2. Find the "CUSTOMIZE HERE" header.
**Required**
3. You HAVE to set the `prod_alias` variable.
  1. this can be an Alias you have pre-configured as a connection locally. You can find this by using the `ls` command within **spin-kit**.
  2. or this can be the username for a production user that has the rights to refresh sandboxes.
**Optional**
4. If you have regularly used dev sandboxes, you can add them to the `sandbox_list` variable in order to make it a bit quicker.  


# Menu
![image](https://github.com/user-attachments/assets/a13963b1-6b1c-4456-bd62-f11bf76f04de)

# result of the `ls` command
![image](https://github.com/user-attachments/assets/498b0ef6-3082-4076-b316-54377ac72ac7)


