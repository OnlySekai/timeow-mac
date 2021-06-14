package main

// #cgo LDFLAGS: -framework CoreFoundation -framework IOKit
// #include <CoreFoundation/CoreFoundation.h>
// #include <IOKit/IOKitLib.h>
//
// const char *
// getSerialNumber()
// {
//     CFMutableDictionaryRef matching = IOServiceMatching("IOPlatformExpertDevice");
//     io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, matching);
//     CFStringRef serialNumber = IORegistryEntryCreateCFProperty(service,
//         CFSTR("IOPlatformSerialNumber"), kCFAllocatorDefault, 0);
//     const char *str = CFStringGetCStringPtr(serialNumber, kCFStringEncodingUTF8);
//     IOObjectRelease(service);
//
//     return str;
// }
import "C"

func getSerialNumber() string {
	return C.GoString(C.getSerialNumber())
}
