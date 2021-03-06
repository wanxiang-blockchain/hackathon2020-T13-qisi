pragma solidity ^0.6.0;
pragma experimental ABIEncoderV2;

import "./openzeppelin/access/Ownable.sol";
import "./device.sol";

contract leasehold is Ownable, device {

    struct Role {
        uint256 funds; // 余额
        address id;
        bytes description;

    }
    // tenantry 租户
    // landlord 房东
    // property 物业
    // factory 家电厂商
    // government 政府/街道

    struct Room {
        address landlord;
        address property; //物业
        address factory; // 家电厂商
        bytes location;    //地址+楼层
        uint256 price;       //房源定价
        bytes area;      //面积
        uint8 status;  // 0: 未装修, 1: 装修中, 2: 房屋可用, 3: 废弃
        bytes description;     //描述
        address[] deviceIds;
    }

    struct Order {
        address from;
        address property;
        uint256 createdAt;
        int64 startAt;
        int64 endAt;
        bytes location;
        uint256 funds; // 订单金额
        uint8 status; // 0: idle, 1: waiting confirm, 2: confirmed, 3: discard
    }

    mapping (bytes => mapping (address => bytes)) devices;
    mapping (address => Role) users;
    mapping (bytes => Room) rooms;

    Order[] allOrders; // 所有订单
    mapping(address => mapping(address => uint256[])) availableOrders;

    event EvtOrderMade(uint256 orderId, address from, address to, bytes location, uint256 amount);
    event EvtTransfer(address to, uint256 amount);
    event EvtRecordRoomRegister(address owner, address property, address factory, bytes location,
        uint256 price, bytes area, uint8 status, bytes description);

    function makeOrder(address to, bytes memory location, int64 startAt, int64 endAt, uint256 funds) public payable {
        Room memory myRoom = getRoom(location);
        require(myRoom.status == 2, "location room is not exists.");
        // getOrderStatus
        uint256[] memory orderIds = availableOrders[msg.sender][to];
        for (uint256 i = 0; i < orderIds.length; i++) {
            require(!isOrderConflict(orderIds[i], startAt, endAt), "exists unfinished order.");
        }

        Order memory myOrder = Order(
            msg.sender,
            myRoom.property,
            now,
            startAt,
            endAt,
            location,
            funds, // 资金存入订单
            1
        );
        allOrders.push(myOrder);
        availableOrders[msg.sender][to].push(allOrders.length - 1);
        emit EvtOrderMade(allOrders.length - 1, msg.sender, to, location, funds);
    }

    function getOrderStatus(uint256 orderId) public view returns (Order memory myOrder) {
        require(orderId < allOrders.length, "Invalid orderId.");
        return allOrders[orderId];
    }

    function isOrderConflict(uint256 orderId, int64 startAt, int64 endAt) public view returns(bool isConflict) {
        Order memory myOrder = getOrderStatus(orderId);
        if (myOrder.status == 1) {
            if ( (startAt >= myOrder.startAt && startAt <= myOrder.endAt) ||
                (endAt >= myOrder.endAt && endAt <= myOrder.endAt) ) {
                return true; // 存在订单
            }
        }
        
        return false;
    }

    // 退房后调用
    function confirmOrder(uint256 orderId) public payable {
        Order memory myOrder = getOrderStatus(orderId);
        require(myOrder.status == 1, "only status == 1 order can be confirmed.");
        require(msg.sender == myOrder.property, "only landlord can confirm this order.");

        shareMoney(myOrder.location, myOrder.funds);
        allOrders[orderId].status = 3; // 完成
    }

    function shareMoney(bytes memory location, uint256 funds) private { // 对订单资金进行分配
        Room memory orderRoom = rooms[location];
        transfer(orderRoom.landlord, funds * 7 / 10);
        transfer(orderRoom.property, funds * 2 / 10);
        transfer(orderRoom.factory, funds *  1 / 10);
    }

    function transfer(address to, uint256 amount) private {
        // TODO use safeMath
        users[to].funds += amount;
    }

    // Room Register Function
    function roomRegister(address property, address factory, bytes memory location,
        uint256 price, bytes memory area, uint8 status, bytes memory description, address[] memory deviceIds) public  {
        //Check room  Register
        require(location.length != 0, "room's location  must be not null");
        Room memory room = rooms[location];
        require(room.location.length == 0, "room has been regist");
        //存储房东入驻的房屋信息
        Room memory myRoom = rooms[location];
        rooms[location].factory = factory;
        rooms[location].factory = factory;
        rooms[location].location = location;
        rooms[location].price = price;
        rooms[location].area = area;
        rooms[location].status = status;
        rooms[location].description = description;
        for (uint256 i = 0; i < deviceIds.length; i++) {
            rooms[location].deviceIds.push(deviceIds[i]);
        }

        //记录放我注册事件到链上
        emit EvtRecordRoomRegister(msg.sender, property, factory, location, price, area, status, description);
    }

    // 装修时调用 registerDevice 将厂家设备进行注册
    function updateRoomInfo(address property, address factory, bytes memory location,
        uint256 price, bytes memory area, uint8 nextStatus, bytes memory description, address[] memory deviceIds) public payable {
        Room memory myRoom = rooms[location];
        require(msg.sender == myRoom.property, "只有物业能够修改");
        rooms[location].factory = factory;
        rooms[location].location = location;
        rooms[location].price = price;
        rooms[location].area = area;
        rooms[location].status = nextStatus;
        rooms[location].description = description;
        for (uint256 i = 0; i < deviceIds.length; i++) {
            rooms[location].deviceIds.push(deviceIds[i]);
        }
    }

    // Get RoomInfo Function
    function getRoom(bytes memory location) public view returns (Room memory roomInfo) {
        return rooms[location];
    }

    function balanceOf(address to) public view returns(uint256) {
        return users[to].funds;
    }

}

