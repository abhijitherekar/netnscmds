# netnscmds

This repo is for some useful tips on container networking, it displays
the following:

The repo will show how to do it programmatically and its equivalent CLI.

1. Creating veth pairs.
2. Getting the network namespaces of the respective containers.
3. Getting into the netns and running cmds.
4. Attaching the veth pairs to containers.

**After building the binary run the following:**
_sudo mv netnscmds /usr/local/bin/_
_sudo chown root:root /usr/local/bin/netnscmds_
_sudo chmod 4755 /usr/local/bin/netnscmds_

