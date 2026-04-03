KERNEL_BZIMAGE = boot/bzImage
INITRAMFS_DIR = initramfs
INITRAMFS_CPIO = initramfs.cpio
UTILS_SRC = userland/coreutils
INIT = userland/init/init.go

.PHONY: all build clean run

# Default target
all: build run

# 1. Build the Go binary and package the initramfs
build:
	@echo "--- Building Go Userland ---"
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $(INITRAMFS_DIR)/init $(INIT)

	mkdir -p $(INITRAMFS_DIR)/bin
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $(INITRAMFS_DIR)/bin $(UTILS_SRC)/*.go

	@echo "--- Packaging Initramfs ---"
	cd $(INITRAMFS_DIR) && find . -print0 | cpio --null -ov --format=newc > ../$(INITRAMFS_CPIO)

# 2. Launch QEMU with the new window settings
run:
	@echo "--- Launching OS ---"
	qemu-system-x86_64 \
    	-kernel ./boot/bzImage \
    	-initrd initramfs.cpio \
    	-m 512M \
    	-append "console=tty1"

# 3. Clean up build artifacts
clean:
	rm -f $(INITRAMFS_CPIO)
	rm -rf $(INITRAMFS_DIR)/*
