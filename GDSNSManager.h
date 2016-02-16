//
//  GDSNSManager.h
//  kuaiji
//
//  Created by gaodun on 15/11/18.
//  Copyright © 2015年 idea. All rights reserved.
//

#import <Foundation/Foundation.h>
@protocol TencentSessionDelegate;

typedef enum _GDSNSLoginType{
    GDSNSLoginTypeQQ = 2,
    GDSNSLoginTypeWX = 3,
    GDSNSLoginTypeWeibo = 1
}GDSNSLoginType;

typedef enum _GDShareType{
    GDShareTypeWeixiSession = 1111,
    GDShareTypeWeixiTimeline,
    GDShareTypeQQ,
    GDShareTypeQQSpace,
    GDShareTypeSinaWeibo,
}GDShareType;

typedef void(^loginHandler)(NSString * openid,NSInteger errCode);

@interface GDSNSManager : NSObject
+ (instancetype)shareGDSNSManager;
+ (void)registerApp;
+(BOOL) isWXAppInstalled;
+(BOOL)handleOpenURL:(NSURL *)url;
- (void)share:(id)object type:(GDShareType)type;
- (void)loginWithType:(GDSNSLoginType)type completionHandler:(loginHandler)handler;
@end

