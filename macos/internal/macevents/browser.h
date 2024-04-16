#import <Cocoa/Cocoa.h>

extern void handleURL(char*, int);

@interface BrowseAppDelegate: NSObject<NSApplicationDelegate>
  - (void)handleGetURLEvent:(NSAppleEventDescriptor *) event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void RunApp();
void StopApp();

struct AppInfo{
  const char* LocalizedName;
  const char* BundleID;  
  const char* BundleURL;
  const char* ExecutableURL;
};

struct AppInfo GetById(int pid);
char* GetLocalizedName(NSRunningApplication* runningApp);
char* GetBundleIdentifier(NSRunningApplication* runningApp);
char* GetBundleURL(NSRunningApplication* runningApp);
char* GetExecutableURL(NSRunningApplication* runningApp);