# Go Assembly - Add Function

This program demonstrates how to implement a simple addition function in Go Assembly.

## Files

- `main.go` - Go program that declares and calls the assembly function
- `add_amd64.s` - AMD64 (x86-64) assembly implementation for Ubuntu/Linux
- `add_arm64.s` - ARM64 assembly implementation for Apple Silicon Macs

## Running

```bash
go run .
```

Output: `Add(10, 32) = 42`

## How It Works

The Go file (main.go:7) declares the `Add` function without a body - this tells the Go compiler to look for an assembly implementation. The assembly file provides the actual implementation.

---

## Detailed Line-by-Line Explanation of AMD64 Assembly (add_amd64.s)

### Line 1: `// func Add(a, b int64) int64`

Comment showing the function signature - this helps document what the assembly function implements

### Line 2: `TEXT ·Add(SB), $0-24`

- `TEXT` - Directive that defines a function
- `·Add` - The function name (the middle dot `·` links it to the current package)
- `(SB)` - Static Base pointer, used for global symbols
- `$0` - Stack frame size (0 bytes, we don't need local variables)
- `-24` - Argument and return size: 8 bytes (a) + 8 bytes (b) + 8 bytes (return) = 24 bytes

### Line 3: `MOVQ a+0(FP), AX`

- `MOVQ` - Move Quadword (64 bits / 8 bytes)
- `a+0(FP)` - First parameter `a` at offset 0 from the Frame Pointer
- `AX` - Load the value into the AX register (accumulator register)

### Line 4: `MOVQ b+8(FP), BX`

- `MOVQ` - Move Quadword
- `b+8(FP)` - Second parameter `b` at offset 8 bytes from Frame Pointer (after the first int64)
- `BX` - Load the value into the BX register (base register)

### Line 5: `ADDQ BX, AX`

- `ADDQ` - Add Quadword
- `BX, AX` - Add BX to AX, storing result in AX (AX = AX + BX)

### Line 6: `MOVQ AX, ret+16(FP)`

- `MOVQ` - Move Quadword
- `AX` - Source: the result of our addition
- `ret+16(FP)` - Destination: return value at offset 16 bytes from Frame Pointer (after both 8-byte parameters)

### Line 7: `RET`

Return from the function back to the Go caller

---

## Understanding Plan9 Assembly Syntax

Go uses **Plan9 assembly syntax**, which originated from the Plan 9 operating system. It's different from Intel or AT&T syntax you might see in traditional x86 assembly.

> AT&T syntax is one of the two main assembly language syntaxes used for x86/x86-64 processors.
> AT&T syntax comes from Unix history at Bell Labs, the same place where the Plan 9 OS (and Go's assembly syntax) originated. They're related but made different design choices!

### Key Characteristics

#### 1. Pseudo-Registers (Virtual Registers)

Plan9 uses four pseudo-registers that don't map directly to hardware:

- **FP (Frame Pointer)**: Points to function arguments and return values
  - Access arguments: `name+offset(FP)`
  - Example: `a+0(FP)` = first argument at offset 0

- **SB (Static Base)**: Points to global memory (like global variables and functions)
  - Used for: `symbol(SB)` or `symbol<>(SB)` for file-local symbols
  - Example: `·Add(SB)` = the Add function in current package

- **SP (Stack Pointer)**: Points to the local stack frame
  - Access local variables: `name+offset(SP)`
  - The real hardware SP is shadowed by this
  - There are two SP registers:
  - the pseudo one is accessed via an offset: `name+offset(SP)`
  - the hardware one who's pointing to the top of the stack in the CPU register: `(SP)`

- **PC (Program Counter)**: Points to the current instruction
  - Rarely used directly in Go assembly

#### 2. Symbol Naming

- **Middle dot (·)**: Links a symbol to the current package
  - `·Add` means the `Add` function in the current Go package
  - Without it, the symbol is considered external

- **Angle brackets (<>)**: Make a symbol file-local (private)
  - `message<>(SB)` = private to this assembly file
  - `message(SB)` = globally visible

#### 3. Operand Order

Plan9 uses **left-to-right** operand order:

```text
MOVQ source, destination
ADDQ value_to_add, accumulator
```

This is like Intel syntax, but different from AT&T syntax (which reverses them).

#### 4. Instruction Suffixes

Instructions have suffixes indicating operand size:

- **B** - Byte (8 bits / 1 byte)
- **W** - Word (16 bits / 2 bytes)
- **L** - Long (32 bits / 4 bytes)
- **Q** - Quadword (64 bits / 8 bytes)

Examples: `MOVB`, `MOVW`, `MOVL`, `MOVQ`

On ARM64, it's slightly different:

- **B** - Byte (8 bits)
- **H** - Halfword (16 bits)
- **W** - Word (32 bits)
- **D** - Doubleword (64 bits)

#### 5. Function Declaration Format

```text
TEXT symbol(SB), flags, $framesize-argsize
```

- `symbol(SB)` - Function name with static base
- `flags` - Like NOSPLIT (don't insert stack-split checks)
- `$framesize` - Size of local variables on stack
- `-argsize` - Total size of arguments + return values

Example: `TEXT ·Add(SB), $0-24`

- No local variables ($0)
- 24 bytes for args/returns (-24)

#### 6. Memory Addressing

Format: `symbol+offset(register)`

Examples:

- `a+0(FP)` - Variable `a` at offset 0 from Frame Pointer
- `message<>+8(SB)` - 8 bytes into the `message` data

### Why Plan9?

Go chose Plan9 assembly because:

1. **Architecture-independent**: Same syntax across different CPUs
2. **Simpler**: Easier to generate code programmatically
3. **Integrated**: Works seamlessly with Go's toolchain
4. **Heritage**: Designed by the same people who created Go (Rob Pike, Ken Thompson)

### Further Reading

- Official Go Assembler Guide: <https://go.dev/doc/asm>
- Plan9 Assembly Manual: <https://9p.io/sys/doc/asm.html>

---

## Key Concepts

1. **Frame Pointer (FP)**: Points to the function arguments and return values
2. **Static Base (SB)**: Used for global symbols and references
3. **Register naming**: Go assembly uses Plan9 syntax with pseudo-registers
   - `AX` - Accumulator register
   - `BX` - Base register
4. **Suffix Q**: Means Quadword (64-bit operations for int64)
5. **File naming**: `add_amd64.s` tells Go to use this file only when building for amd64 architecture
6. **Memory layout**:

   ```text
   FP+0:  parameter a (8 bytes)
   FP+8:  parameter b (8 bytes)
   FP+16: return value (8 bytes)
   ```

## ARM64 Differences

The ARM64 version (add_arm64.s) is nearly identical but uses:

- `MOVD` instead of `MOVQ` (Move Doubleword)
- `ADD` instead of `ADDQ`
- Different register names: `R0`, `R1` instead of `AX`, `BX`
