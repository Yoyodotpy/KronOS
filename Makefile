KERNEL_BZIMAGE = boot/bzImage
INITRAMFS_DIR = initramfs
INITRAMFS_CPIO = initramfs.cpio
UTILS_SRC = userland/bin
INIT = ./userland/init

GO_FLAGS = CGO_ENABLED=0 GOOS=linux
LD_FLAGS = -a -ldflags '-extldflags "-static"'


.PHONY: all build clean run

# Default target
all: build run

# 1. Build the Go binary and package the initramfs
build:
	@echo "making init"
	$(GO_FLAGS) go build -a $(LD_FLAGS) -o $(INITRAMFS_DIR)/init $(INIT)

	@echo "making utils"
	mkdir -p $(INITRAMFS_DIR)/bin
	cp -r boot/root/. $(INITRAMFS_DIR)
	@for dir in $(shell ls -d $(UTILS_SRC)/*/); do \
		utilname=$$(basename $$dir); \
		echo "making $$utilname"; \
		$(GO_FLAGS) go build -a $(LD_FLAGS) -o $(INITRAMFS_DIR)/bin/$$utilname ./$(UTILS_SRC)/$$utilname; \
	done

	@echo "--- Packaging Initramfs ---"
	cd $(INITRAMFS_DIR) && find . -print0 | cpio --null -ov --format=newc > ../$(INITRAMFS_CPIO)

# 2. Launch QEMU with the new window settings
run:
	@echo "starting qemu"
	qemu-system-x86_64 \
		-kernel ./boot/bzImage \
		-initrd initramfs.cpio \
		-drive file=./boot/disk.img,format=raw,if=ide \
		-cpu host \
		-enable-kvm \
		-vga std \
		-m 8G \
		-display sdl \
		-append "console=tty1 video=1024x768-32 rdinit=/init"

# 3. Clean up build artifacts
clean:
	rm -f $(INITRAMFS_CPIO)
	rm -rf $(INITRAMFS_DIR)/*
