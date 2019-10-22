# WSMatcher

## brief 
WSMtcher has many rooms that can make clients communication who has same roomid 

### useage: WSMatcher {listenAddr, default: ':3333'} 
etc: WSMatcher :3333 
      WSMatcher 127.0.0.1:3333 
      WSMatcher 0.0.0.0:3333 

### client useage: ws://{host:port}/wsmatcher/{roomid}?type={typeid} 
etc: ws://192.168.1.3:3333/wsmatcher/23f32g32g3g?type=1 

#### params:
* type: room type
  - 1: only 2 client, pure-tran data without rewrite and watch at WSMatcher, one fd close then close another
