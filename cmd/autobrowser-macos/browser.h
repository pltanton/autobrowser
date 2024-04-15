#import <Cocoa/Cocoa.h>

extern void HandleURL(char*, int);

@interface BrowseAppDelegate: NSObject<NSApplicationDelegate>
  - (void)handleGetURLEvent:(NSAppleEventDescriptor *) event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void RunApp();
NSRunningApplication * GetById(int pid);
char* GetLocalizedName(NSRunningApplication* runningApp);
char* GetBundleIdentifier(NSRunningApplication* runningApp);
char* GetBundleURL(NSRunningApplication* runningApp);
char* GetExecutableURL(NSRunningApplication* runningApp);
