//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

struct connect_event {
  u32 pid;
  u8  ip[16];
  u16 port;
  char comm[16];
};

struct
{
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 256 * 1024);
} events SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_connect")
int detect_connect(struct trace_event_raw_sys_enter *ctx)
{
  void *addr = (void *)ctx->args[1];

  // 1. Safely read the address family (IPv6 = 10)
  u16 family = 0;
  bpf_probe_read_user(&family, sizeof(family), addr);

  // 2. Filter: only trace IPv6 connections (AF_INET6)
  if (family != 10) return 0;

  // 3. Reserve space in the Ring Buffer
  struct connect_event *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
  if (!e) return 0;

  // 4. Populate event metadata
  e->pid = bpf_get_current_pid_tgid() >> 32;
  bpf_get_current_comm(e->comm, sizeof(e->comm));

  // 5. Read IPv6 address and port
  struct sockaddr_in6 sa6;
  bpf_probe_read_user(&sa6, sizeof(sa6), addr);
  __builtin_memcpy(e->ip, &sa6.sin6_addr, 16);
  e->port = bpf_ntohs(sa6.sin6_port);

  // 6. Submit the event to user space
  bpf_ringbuf_submit(e, 0);
  return 0;
}

char _license[] SEC("license") = "GPL";