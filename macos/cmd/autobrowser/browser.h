#import <Cocoa/Cocoa.h>

extern void HandleURL(char*, int);

@interface BrowseAppDelegate: NSObject<NSApplicationDelegate>
  - (void)handleGetURLEvent:(NSAppleEventDescriptor *) event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void RunApp();