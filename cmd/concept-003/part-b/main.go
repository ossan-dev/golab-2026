// this runs in the user space
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang Bpf ./bpf/probe.c

// --- how to run C and Go together ---
// sudo apt update
// sudo apt install clang llvm libbpf-dev linux-headers-$(uname -r)
// sudo ln -s /usr/include/x86_64-linux-gnu/asm /usr/include/asm
// go generate ./...
// go build -o observer .
// sudo -E ./observer --observed-process "secure_binary"

// ConnectEvent must perfectly match the updated C struct (IPv6 only)
type ConnectEvent struct {
	PID  uint32
	IP   [16]byte
	Port uint16
	Comm [16]byte
}

func main() {
	observedProcess := flag.String("observed-process", "", "The process name to trace (empty traces all)")
	flag.Parse()

	// 1. Allow the Go program to load eBPF into the kernel
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("removing memlock: %v", err)
	}

	// 2. Load our compiled C probe into the kernel
	objs := BpfObjects{}
	if err := LoadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// 3. attach the eBPF program to the tracepoint!
	kp, err := link.Tracepoint("syscalls", "sys_enter_connect", objs.DetectConnect, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	defer kp.Close()

	// 4. Read from the Ring Buffer...
	rd, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %s", err)
	}
	defer rd.Close()

	log.Println("🚨 eBPF Go Spy is running! Watching for network connections...")
	log.Println("👉 Try running './secure_binary' in another terminal.")

	// 6. handle Ctrl+C gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Received signal, exiting...")
		rd.Close()
	}()

	// 5. Loop and read events sent from the kernel
	var event ConnectEvent
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("Reader closed, exiting...")
				return
			}
			log.Printf("reading from reader: %s", err)
			continue
		}

		// The C code pushed a __u32 (4 bytes). Let's read it.
		// record.RawSample contains the exact bytes pushed by bpf_ringbuf_output
		reader := bytes.NewReader(record.RawSample)
		if err := binary.Read(reader, binary.LittleEndian, &event); err != nil {
			log.Printf("parsing event: %v", err)
			continue
		}

		ipAddr := net.IP(event.IP[:])

		n := bytes.IndexByte(event.Comm[:], 0)
		if n == -1 {
			n = len(event.Comm)
		}
		comm := string(event.Comm[:n])

		if *observedProcess != "" && comm != *observedProcess {
			continue
		}

		fmt.Printf("🚀 PID: %-6d | Comm: %-16s | Destination: %-39s:%d\n",
			event.PID, comm, ipAddr.String(), event.Port)
	}
}
