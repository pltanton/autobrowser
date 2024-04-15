#include "browser.h"

@implementation BrowseAppDelegate
- (void)applicationWillFinishLaunching:(NSNotification *)aNotification
{
  NSAppleEventManager *appleEventManager = [NSAppleEventManager sharedAppleEventManager];
  [appleEventManager setEventHandler:self
                         andSelector:@selector(handleGetURLEvent:withReplyEvent:)
                         forEventClass:kInternetEventClass andEventID:kAEGetURL];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender
{
  return NSTerminateNow;
}

- (void)handleGetURLEvent:(NSAppleEventDescriptor *)event
           withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
  handleURL((char*)[[[event paramDescriptorForKeyword:keyDirectObject] stringValue] UTF8String],
   (int) [[event attributeDescriptorForKeyword:keySenderPIDAttr] int32Value]);
}
@end

void RunApp(void) {
  [NSAutoreleasePool new];
  [NSApplication sharedApplication];
  BrowseAppDelegate *app = [BrowseAppDelegate alloc];
  [NSApp setDelegate:app];
  [NSApp run];
}

struct AppInfo GetById(int pid) {
  NSRunningApplication* runningApp = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];

  struct AppInfo result;
  result.LocalizedName = [[runningApp localizedName] UTF8String];
  result.BundleID = [[runningApp bundleIdentifier] UTF8String];
  result.BundleURL = [[[runningApp bundleURL] absoluteString] UTF8String];
  result.ExecutableURL = [[[runningApp executableURL] absoluteString] UTF8String];

  return result;
}
char* GetLocalizedName(NSRunningApplication* runningApp) {
  return [[runningApp localizedName] UTF8String];
}
char* GetBundleIdentifier(NSRunningApplication* runningApp) {
  return [[runningApp bundleIdentifier] UTF8String];
}
char* GetBundleURL(NSRunningApplication* runningApp) {
  return [[[runningApp bundleURL] absoluteString] UTF8String];
}
char* GetExecutableURL(NSRunningApplication* runningApp) {
  return [[[runningApp executableURL] absoluteString] UTF8String];
}