//
//  GDSNSManager.m
//  kuaiji
//
//  Created by gaodun on 15/11/18.
//  Copyright © 2015年 idea. All rights reserved.
//
#define kSinaRedirectURL @"http://www.gaodun.com/oauth/weibo/callback.php"

#import "GDSNSManager.h"
#import "WXApi.h"
#import "WXApiObject.h"
#import <WeiboSDK.h>
#import <TencentOpenAPI/TencentOAuth.h>
#import <TencentOpenAPI/QQApiInterface.h>

@interface GDSNSManager()<TencentSessionDelegate,QQApiInterfaceDelegate,WXApiDelegate,WeiboSDKDelegate>

@property (nonatomic,strong) TencentOAuth *tencentOAuth;
@property (nonatomic,copy) loginHandler loginBlock;
//@property (nonatomic,strong) TencentOAuth *tencentOAuth;
@end
@implementation GDSNSManager
+ (instancetype)shareGDSNSManager{
    static GDSNSManager *gdsnsmanager;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        gdsnsmanager = [[self alloc] init];
    });
    return gdsnsmanager;
}

+ (void)registerApp
{
    GDSNSManager* __self = [GDSNSManager shareGDSNSManager];
    //微信注册
    [WXApi registerApp:@"wx2d4a1bb5e5c25332"];
    //qq注册
    __self.tencentOAuth  = [[TencentOAuth alloc]initWithAppId:@"1104945900"andDelegate:__self];
    [__self.tencentOAuth openSDKWebViewQQShareEnable];
    __self.tencentOAuth.redirectURI = @"www.qq.com";
    __self.tencentOAuth.sessionDelegate = __self;
    //微博注册
    [WeiboSDK enableDebugMode:YES];
    [WeiboSDK registerApp:@"3995371278"];
    
    
    NSLog(@"GDSNSManager <TencentQQ->%@> <wechat->%@> <weibo->%@>",[TencentOAuth sdkVersion],[WXApi getApiVersion],[WeiboSDK getSDKVersion]);
}




#pragma mark --handleURL

+(BOOL)handleOpenURL:(NSURL *)url{
    GDSNSManager *__self = [GDSNSManager shareGDSNSManager];
    NSString *urlString = url.absoluteString;
    if ([urlString rangeOfString:@"wb"].length>0) {
        return [ WeiboSDK handleOpenURL:url delegate:__self ];
    }else if([urlString rangeOfString:@"tencent"].length>0) {
        return [TencentOAuth HandleOpenURL:url];
    }else if([urlString rangeOfString:@"QQ"].length>0) {
        return [TencentOAuth HandleOpenURL:url];
    }else if([urlString rangeOfString:@"wx"].length>0) {
        return [WXApi handleOpenURL:url delegate:__self];
    }
    return YES;
}
#pragma mark - SinaWeibo
#pragma mark -
- (void)didReceiveWeiboRequest:(WBBaseRequest *)request
{
    
}

- (void)didReceiveWeiboResponse:(WBBaseResponse *)response
{
    if ([response isKindOfClass:WBAuthorizeResponse.class])
    {
        WBAuthorizeResponse * authResponse = (WBAuthorizeResponse *)response;
        [self loginresult:authResponse.userID errCode:authResponse.statusCode];
        
    } else if ([response isKindOfClass:WBSendMessageToWeiboResponse.class]){
        
        if (response.statusCode == 0) {

        } else {

        }

    }
}
#pragma mark - wxdelegate

-(void) onReq:(BaseReq*)req
{
    
}
//onReq是微信终端向第三方程序发起请求，要求第三方程序响应。第三方程序响应完后必须调用sendRsp返回。在调用sendRsp返回时，会切回到微信终端程序界面。
-(void) onResp:(BaseResp*)resp
{
    NSString* code;
    NSString* kWeiXinAppKey = @"wx2d4a1bb5e5c25332";
    NSString* kWeiXinAppSecret = @"efce4b7ce9acc34e9b4ce1639a680e23";
    if ([resp isKindOfClass:[SendAuthResp class]]) {
        
        SendAuthResp *response = (SendAuthResp *)resp;
        if (response.errCode < 0) {
            [self loginresult:nil errCode:response.errCode];
            return ;
        }else if (response.errCode == 0){
            code = [response code];
            NSString *WXURL = @"https://api.weixin.qq.com/sns/oauth2/access_token";
            NSString * url = [NSString stringWithFormat:@"%@?%@=%@&%@=%@&%@=%@&%@=%@",WXURL,@"appid",kWeiXinAppKey,@"secret",kWeiXinAppSecret,@"code",code, @"grant_type",@"authorization_code"];
            NSURL * nsurl = [NSURL URLWithString:url];
            NSURLRequest * request = [NSURLRequest requestWithURL:nsurl];
            
            [NSURLConnection sendAsynchronousRequest:request queue:[[NSOperationQueue alloc] init] completionHandler:^(NSURLResponse * _Nullable response, NSData * _Nullable data, NSError * _Nullable connectionError) {
                if (connectionError == nil) {
                    NSError *error = nil;
                    id jsonObject = [NSJSONSerialization JSONObjectWithData:data options:NSJSONReadingAllowFragments error:&error];
                    if ([jsonObject isKindOfClass:[NSDictionary class]]) {
                        [self loginresult:[jsonObject objectForKey:@"openid"] errCode:0];
                        return;
                    }
                    [self loginresult:nil errCode:-5];//格式错误
                } else {
                    [self loginresult:nil errCode:-4];//请求错误
                }

            }];
        }
    }
}
#pragma mark - tecent delegate
- (void)tencentDidLogin
{
    if ([self.tencentOAuth getUserInfo]){
        NSString * userid = self.tencentOAuth.openId;
        [self loginresult:userid errCode:0];
    }
}
- (void)tencentDidNotLogin:(BOOL)cancelled
{
    [self loginresult:nil errCode:-6];
}
- (void)tencentDidNotNetWork
{
    [self loginresult:nil errCode:-7];
}
- (void)isOnlineResponse:(NSDictionary *)response
{

}
#pragma mark - Share
- (void)share:(id)object type:(GDShareType)type
{
    if (type == GDShareTypeSinaWeibo) {
        
        WBAuthorizeRequest *authRequest = [WBAuthorizeRequest request];
        authRequest.redirectURI = kSinaRedirectURL;
        authRequest.scope = @"all";
        
        WBMessageObject *message = [WBMessageObject message];
        message.text = @"测试通过WeiboSDK发送文字到微博!";
        WBWebpageObject *webpage = [WBWebpageObject object];
        webpage.objectID = @"identifier1";
        webpage.title = NSLocalizedString(@"分享网页标题", nil);
        webpage.description = [NSString stringWithFormat:NSLocalizedString(@"分享网页内容简介-%.0f", nil), [[NSDate date] timeIntervalSince1970]];
        webpage.thumbnailData = [NSData dataWithContentsOfFile:[[NSBundle mainBundle] pathForResource:@"image_1" ofType:@"jpg"]];
        webpage.webpageUrl = @"http://sina.cn?a=1";
        message.mediaObject = webpage;
        
        WBSendMessageToWeiboRequest *request = [WBSendMessageToWeiboRequest requestWithMessage:message authInfo:nil access_token:nil];
        request.userInfo = @{@"ShareMessageFrom": @"SendMessageToWeiboViewController",
                             @"Other_Info_1": [NSNumber numberWithInt:123],
                             @"Other_Info_2": @[@"obj1", @"obj2"],
                             @"Other_Info_3": @{@"key1": @"obj1", @"key2": @"obj2"}};

        [WeiboSDK sendRequest:request];
    }
    if (type == GDShareTypeQQ || type == GDShareTypeQQSpace) {
        NSString *url = @"http://www.baidu.com/";
        //分享图预览图URL地址
        NSString *previewImageUrl = @"preImageUrl.png";
        QQApiNewsObject *newsObj = [QQApiNewsObject
                                    objectWithURL:[NSURL URLWithString:url]
                                    title: @"title"
                                    description:@"description"
                                    previewImageURL:[NSURL URLWithString:previewImageUrl]];
        SendMessageToQQReq *req = [SendMessageToQQReq reqWithContent:newsObj];
        //将内容分享到qzone
        if (type == GDShareTypeQQ) {
            [QQApiInterface sendReq:req];
        }
        if (type == GDShareTypeQQSpace) {
            [QQApiInterface SendReqToQZone:req];
        }
        
    }
    if (type == GDShareTypeWeixiSession || type == GDShareTypeWeixiTimeline){
        WXMediaMessage *message = [WXMediaMessage message];
        message.title = @"专访张小龙：产品之上的世界观";
        message.description = @"微信的平台化发展方向是否真的会让这个原本简洁的产品变得臃肿？在国际化发展方向上，微信面临的问题真的是文化差异壁垒吗？腾讯高级副总裁、微信产品负责人张小龙给出了自己的回复。";
        [message setThumbImage:[UIImage imageNamed:@"res2.png"]];
        
        WXWebpageObject *ext = [WXWebpageObject object];
        ext.webpageUrl = @"http://tech.qq.com/zt2012/tmtdecode/252.htm";
        
        message.mediaObject = ext;
        
        SendMessageToWXReq* req = [[SendMessageToWXReq alloc] init];
        req.bText = NO;
        req.message = message;
        if (type == GDShareTypeWeixiSession) {
            req.scene = WXSceneSession;
        }
        if (type == GDShareTypeWeixiTimeline) {
            req.scene = WXSceneTimeline;
        }
        
        
        [WXApi sendReq:req];
    }
}

- (void)loginWithType:(GDSNSLoginType)type completionHandler:(loginHandler)handler
{
    self.loginBlock = handler;
    if (type == GDSNSLoginTypeQQ) {
        [self.tencentOAuth authorize:@[kOPEN_PERMISSION_GET_USER_INFO,
                                       kOPEN_PERMISSION_GET_SIMPLE_USER_INFO,
                                       kOPEN_PERMISSION_GET_INFO,
                                       kOPEN_PERMISSION_GET_INFO
                                       ]];
    }
    if (type == GDSNSLoginTypeWX) {
        SendAuthReq* req =[[SendAuthReq alloc ] init ];
        req.scope = @"snsapi_userinfo" ;
        req.state = @"123" ;
        //第三方向微信终端发送一个SendAuthReq消息结构
        [WXApi sendReq:req];
    }
    if (type == GDSNSLoginTypeWeibo) {
        WBAuthorizeRequest *request = [WBAuthorizeRequest request];
        request.redirectURI = kSinaRedirectURL;
        request.scope = @"email";
        //request.userInfo = @{};
        [WeiboSDK sendRequest:request];
    }
}

- (void)loginresult:(NSString*)str errCode:(NSInteger)code
{
    if (self.loginBlock)
    {
        self.loginBlock(str,code);
        self.loginBlock = nil;
    }
}
+(BOOL)isWXAppInstalled
{
    return [WXApi isWXAppInstalled];
}
@end
