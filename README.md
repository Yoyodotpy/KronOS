Custom Operating System based on the Linux Kernel, everything else written from scratch

### HOW TO RUN THE RELEASE IMG:
qemu-system-x86_64 -m 512M -bios /path/to/OVMF.fd -drive format=raw,file=kronos.img -cpu host -enable-kvm

Replace /path/to/OVMF.fd with your actual path to OVMF.fd, which you can find on linux using:
find /usr/share -name "*OVMF*"

Once you get running, I recommend running "nfb" to start the custom wm alongside bad apple and doom, as well as edit and lua for a text editor and the lua interpretter. All the required and usual core utils will be available.


### TO BUILD AND TEST IN INITRAMFS

*recommended if disk image doesnt work*

To start operating system, you need qemu and go installed. Then, run "make" and it should build the operating system and start a vm.
