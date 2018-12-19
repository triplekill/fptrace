//+build ignore

package main

import (
	"fmt"
	"os"

	"github.com/orivej/e"
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
)

const (
	SECCOMP_RET_TRACE = 0x7ff00000
	SECCOMP_RET_ALLOW = 0x7fff0000
)

var syscalls = []uint32{
	unix.SYS_CHDIR,
	unix.SYS_CLOSE,
	unix.SYS_DUP,
	unix.SYS_DUP2,
	unix.SYS_DUP3,
	unix.SYS_EXECVE,
	unix.SYS_EXECVEAT,
	unix.SYS_FCHDIR,
	unix.SYS_FCNTL,
	unix.SYS_LINK,
	unix.SYS_LINKAT,
	unix.SYS_OPEN,
	unix.SYS_OPENAT,
	unix.SYS_PIPE,
	unix.SYS_PREAD64,
	unix.SYS_PREADV,
	unix.SYS_PREADV2,
	unix.SYS_PWRITE64,
	unix.SYS_PWRITEV,
	unix.SYS_PWRITEV2,
	unix.SYS_READ,
	unix.SYS_READV,
	unix.SYS_RENAME,
	unix.SYS_RENAMEAT,
	unix.SYS_RENAMEAT2,
	unix.SYS_RMDIR,
	unix.SYS_UNLINK,
	unix.SYS_UNLINKAT,
	unix.SYS_WRITE,
	unix.SYS_WRITEV,
}

func main() {
	n := len(syscalls)
	p := make([]bpf.Instruction, n+3)
	p[0] = bpf.LoadAbsolute{Off: 0, Size: 4}
	p[n+1] = bpf.RetConstant{Val: SECCOMP_RET_ALLOW}
	p[n+2] = bpf.RetConstant{Val: SECCOMP_RET_TRACE}
	for i, sc := range syscalls {
		p[i+1] = bpf.JumpIf{Cond: bpf.JumpEqual, Val: sc, SkipTrue: uint8(n - i)}
	}
	ins, err := bpf.Assemble(p)
	e.Exit(err)

	os.Stdout, err = os.Create("seccomp.h")
	e.Exit(err)
	fmt.Println("// Code generated by ./seccomp.go. DO NOT EDIT.\n")
	fmt.Println("#include <linux/filter.h>\n")
	fmt.Println("struct sock_filter seccomp_filter[] = {")
	for _, in := range ins {
		fmt.Printf("\t{%#x, %d, %d, %#x},\n", in.Op, in.Jt, in.Jf, in.K)
	}
	fmt.Println("};\n")
	fmt.Printf("struct sock_fprog seccomp_program = {%d, seccomp_filter};\n", len(ins))
}
