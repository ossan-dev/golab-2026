
/* WARNING: Removing unreachable block (ram,0x006317a2) */
/* Golang function info: Flags: []
   Golang source: /home/ossan/Projects/golab-2026/cmd/plain/main.go:12
   Golang stacktrace signature: func main.main() ??? */

multireturn{int;error} main::main.main(void)

{
  os/exec.Cmd *self;
  bytes.Reader *pbVar1;
  unsafe.Pointer extraout_RCX;
  int extraout_RBX;
  string name;
  io.Writer w;
  any e;
  any e_00;
  multireturn{int;error} mVar2;
  multireturn{[]uint8;error} mVar3;
  []string arg;
  []any a;
  internal/abi.Type *local_30;
  undefined **local_28;
  char *local_20;
  undefined8 local_18;
  net/http.Client *local_10;
  
  while (&stack0x00000000 <= CURRENT_G.stackguard0) {
    runtime::runtime.morestack_noctxt();
  }
  local_18 = 0xf;
  local_20 = "/etc/os-release";
  name.len = 3;
  name.str = "cat";
  arg.len = 1;
  arg.array = (string *)&local_20;
  arg.cap = 1;
  self = os/exec::os/exec.Command(name,arg);
  mVar3 = os/exec::os/exec.(*Cmd).Output(self);
  if (mVar3.~r1.tab != (error_itab *)0x0) {
    e_00.data = mVar3.~r1.data;
    e_00.tab = (interface {}_itab *)(mVar3.~r1.tab)->Type;
                    /* WARNING: Subroutine does not return */
    runtime::runtime.gopanic(e_00);
  }
  local_10 = runtime::runtime.newobject(&net/http::net/http.Client___Struct_type.Type);
  pbVar1 = runtime::runtime.newobject(&bytes::bytes.Reader___Struct_type.Type);
  (pbVar1->s).len = mVar3.~r0.len;
  (pbVar1->s).cap = mVar3.~r0.cap;
  (pbVar1->s).array = mVar3.~r0.array;
  pbVar1->i = 0;
  pbVar1->prevRune = -1;
  net/http::net/http.(*Client).Post(local_10);
  if (extraout_RBX != 0) {
    e.data = extraout_RCX;
    e.tab = *(interface {}_itab **)(extraout_RBX + 8);
                    /* WARNING: Subroutine does not return */
    runtime::runtime.gopanic(e);
  }
  local_30 = &string___String_type;
  local_28 = &PTR_DAT_00783f80;
  w.data = os.Stdout;
  w.tab = (io.Writer_itab *)&os::*os.File__implements__io.Writer__itab;
  a.len = 1;
  a.array = (any *)&local_30;
  a.cap = 1;
  mVar2 = fmt::fmt.Fprintln(w,a);
  return mVar2;
}
