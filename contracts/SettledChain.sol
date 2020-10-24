pragma solidity >=0.4.0 <0.7.0;

contract SettledChain {

    struct Room {
        address owner;     //房东Id
        bytes location;    //地址+楼层
        bool isAvailable;  //房屋是否可用
        int64 price;       //房源定价   
        bytes area;      //面积
        bytes decorateStatus;      //装修状态
        bytes description;     //描述
    }

    mapping(bytes => Room)  public rooms;  //房屋信息

    event RecoredRoomRegisterEvent(address owner, bytes location, bool isAvailable, int64 price,
        bytes area, bytes decorateStatus, bytes description);
    // Room Register Function
    function roomRegister(address owner, bytes memory location, bool isAvailable, int64 price,
        bytes memory area, bytes memory decorateStatus, bytes memory description) public  {
        //存储房东入驻的房屋信息
        rooms[location] = Room(owner, location, isAvailable, price, area, decorateStatus, description);
        //记录放我注册事件到链上
        emit RecoredRoomRegisterEvent(owner, location, isAvailable, price, area, decorateStatus, description);

    }
    // Get RoomInfo Function
    function getRoom(bytes memory location) public view returns (address owner, bytes memory, bool isAvailable, int64 price,
        bytes memory area, bytes memory decorateStatus, bytes memory description) {
        Room memory room = rooms[location];
        owner = room.owner;
        //房东Id
        // location1 = room.location;    //地址+楼层
        isAvailable = room.isAvailable;
        //房屋是否可用
        price = room.price;
        //房源定价
        area = room.area;
        //面积
        decorateStatus = room.decorateStatus;
        //装修状态
        description = room.description;
        //描述
        return (owner, location, isAvailable, price, area, decorateStatus, description);
    }
}
