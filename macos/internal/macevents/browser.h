#import <Cocoa/Cocoa.h>

extern void handleURL(char*, int);

@interface BrowseAppDelegate: NSObject<NSApplicationDelegate>
  - (void)handleGetURLEvent:(NSAppleEventDescriptor *) event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void RunApp();

char* GetLocalizedName(NSRunningApplication* runningApp);
char* GetBundleIdentifier(NSRunningApplication* runningApp);
char* GetBundleURL(NSRunningApplication* runningApp);
char* GetExecutableURL(NSRunningApplication* runningApp);