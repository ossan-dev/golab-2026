# RLIMIT_MEMLOCK

**RLIMIT_MEMLOCK** stands for `Resource Limit - Memory Locked`. Under the hood, it is a kernel-enforced limit that restricts how many bytes of virtual memory a process can lock into physical RAM.

## 0. What does "Locking Memory" mean?

Normally, the operating system can move pages of process memory between physical RAM and swap space (on
disk) depending on memory pressure. When a process locks memory:

- *It forces the pages to stay in physical RAM*.
- *The OS is prohibited from swapping those pages out to disk.*
- This guarantees deterministic access times and prevents latency spikes.

In `main.go`, the call to `rlimit.RemoveMemlock()` raises the current process's `RLIMIT_MEMLOCK` (locked-in-memory size) limit to **infinity**.

Here is why this is necessary:

## 1. eBPF Memory Pinning

eBPF programs and maps must reside permanently in physical memory (RAM).  
They cannot be swapped/paged out
because the kernel needs guaranteed, instant access to them during execution (e.g., inside interrupt
handlers or critical syscall paths).

## 2. BPF Map Memory Allocation

When the Go loader initializes maps (like the events Ring Buffer), they are locked into memory.  
This
allocation counts against the user-space process's RLIMIT_MEMLOCK resource limit.

## 3. Evading Size Restrictions

On older Linux kernels (**pre-5.11**), the default RLIMIT_MEMLOCK limit is very small (**often 64KB**). Since
loading almost any real-world eBPF program or map exceeds this, the kernel will reject the bpf() system
call with EPERM or ENOMEM.

> In Linux **5.11 and newer**, the kernel uses cgroup-based memory accounting (**memcg**) instead of RLIMIT_MEMLOCK to track eBPF allocations.  
> Therefore, while no longer strictly necessary on modern kernels, calling `rlimit.RemoveMemlock()` remains a best practice to ensure backward compatibility.
