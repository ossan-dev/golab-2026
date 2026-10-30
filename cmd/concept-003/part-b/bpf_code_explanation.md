# Explanation of eBPF Connect Probe (`bpf/probe.c`)

This guide explains the eBPF helper functions and the program logic inside [bpf/probe.c](file:///home/ossan/Projects/golab-2026/cmd/concept-003/part-b/bpf/probe.c) to help you understand and explain it to others.

---

## 1. eBPF Helper Functions Explained

In eBPF, the code runs inside a sandboxed environment in the Linux kernel. Because of this, it cannot access system memory or resources directly. It must use special kernel-provided functions called **eBPF Helpers**. Here is exactly what the ones in your code do:

### `bpf_get_current_pid_tgid()`

* **What it does**: Returns a 64-bit unsigned integer containing the current thread and process identifiers.
  * **Top 32 bits**: Thread Group ID (TGID), which is what user space calls the **Process ID (PID)**.
  * **Bottom 32 bits**: Thread ID (TID), which is the unique identifier for the specific thread.
* **Why it's used**: To grab the PID of the process initiating the connection. Since we only want the PID, we shift the 64-bit result right by 32 bits (`>> 32`).

### `bpf_probe_read_user(void *dst, u32 size, const void *user_ptr)`

* **What it does**: Safely copies `size` bytes from the user-space address `user_ptr` to the kernel-space/stack buffer `dst`.
* **Why it's used**: The kernel cannot directly dereference pointers pointing to user-space memory (like the `struct sockaddr *` argument passed to the `connect` syscall). User-space memory might be swapped out, unmapped, or invalid. If we dereference it directly, the kernel would crash. `bpf_probe_read_user` checks that the address is safe to read before copying it.

### `bpf_ringbuf_reserve(void *ringbuf_map, u64 size, u64 flags)`

* **What it does**: Reserves a block of memory of size `size` in the shared ring buffer map. It returns a pointer to this reserved space inside the ring buffer.
* **Why it's used**: Instead of allocating memory on the kernel stack (which has a very tight 512-byte limit), we reserve memory directly inside the shared queue (ring buffer). This memory is automatically zero-initialized in modern Linux kernels.

### `bpf_get_current_comm(void *buf, u32 size_of_buf)`

* **What it does**: Copies the command name (the name of the executable, e.g. `curl`, `chrome`) of the current process into the buffer `buf`. The maximum name length is 16 bytes (defined by `TASK_COMM_LEN` in the kernel), including the null terminator.
* **Why it's used**: To associate the connection with the name of the process that initiated it.

### `bpf_ntohs(u16 netshort)`

* **What it does**: Converts a 16-bit integer from **Network Byte Order** (Big-Endian) to **Host Byte Order** (usually Little-Endian on x86 and ARM).
* **Why it's used**: IP ports are sent over the network in Big-Endian format. Since computers usually run on Little-Endian, we must convert the port number so it displays correctly (e.g., converting `0x901F` to `8080`).

### `bpf_ringbuf_submit(void *event, u64 flags)`

* **What it does**: Submits the reserved memory `event` to the ring buffer, making it immediately visible and readable to the user-space Go program.
* **Why it's used**: Once we finish writing all the connection info (PID, IP, Port, Comm) to the reserved memory, we submit it to send it to the Go side.

### `__builtin_memcpy(void *dst, const void *src, u32 size)`

* **What it does**: Copies `size` bytes from the memory area `src` directly to the memory area `dst`. It is a compiler-builtin version of the standard C library function `memcpy`.
* **Why it's used**: Standard C library functions (like `memcpy`) are not available inside the eBPF kernel sandbox. Instead, compiler builtins (provided by Clang/LLVM) are used. The compiler optimizes and inlines this function directly into efficient load/store instructions that are fully compatible with and verifiable by the eBPF verifier. In this probe, it is used to copy the IPv6 address bytes into our `connect_event` struct.

---

## 2. Review of Your Comments in the Code

The comments in [bpf/probe.c](file:///home/ossan/Projects/golab-2026/cmd/concept-003/part-b/bpf/probe.c) have been simplified alongside the code to maintain a clean, step-by-step walkthrough suitable for presentation.

### Refined Step-by-Step Flow

1. **Step 1: Safely read the address family** (Line 24)
   `// 1. Safely read the address family (IPv6 = 10)`
2. **Step 2: Filter non-IPv6 connections** (Line 28)
   `// 2. Filter: only trace IPv6 connections (AF_INET6)`
3. **Step 3: Reserve space in the Ring Buffer** (Line 31)
   `// 3. Reserve space in the Ring Buffer`
4. **Step 4: Populate event metadata** (Line 35)
   `// 4. Populate event metadata`
5. **Step 5: Read IPv6 address and port** (Line 39)
   `// 5. Read IPv6 address and port`
6. **Step 6: Submit the event to user space** (Line 45)
   `// 6. Submit the event to user space`
