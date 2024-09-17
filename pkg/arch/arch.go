package arch

import "runtime"

// IsArm returns whether the current CPU architecture is ARM.
// IsArm 返回当前 CPU 架构是否为 ARM。
func IsArm() bool {
	return runtime.GOARCH == "arm" || runtime.GOARCH == "arm64"
}

// IsX86 returns whether the current CPU architecture is X86.
// IsX86 返回当前 CPU 架构是否为 X86。
func IsX86() bool {
	return runtime.GOARCH == "386" || runtime.GOARCH == "amd64"
}

// Is64Bit returns whether the current CPU architecture is 64-bit.
// Is64Bit 返回当前 CPU 架构是否为 64 位。
func Is64Bit() bool {
	return runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
}
