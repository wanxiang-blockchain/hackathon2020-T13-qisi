pragma solidity ^0.6.0;
pragma experimental ABIEncoderV2;


contract device {
    struct Dev {
        bytes location;
        address factory;
        bytes category;
    }

    struct Power {
        uint256 createdAt;
        uint256 consumption;
    }
    // deviceId => deviceLogs
    mapping(address => Power[]) devLogs; // 设备运行日志
    mapping(address => mapping(bytes => Dev)) devInfos; // 设备详情

    event EvtRegisterDevice(bytes location, address factory, bytes category);

    function registerDevice(bytes memory location, address deviceId, address factory, bytes memory category) public payable {
        // msg.sender is user
        Dev memory myDev = devInfos[msg.sender][category];
        require(myDev.factory == address(0), "[Device](registerDevice) device exists.");

        devInfos[msg.sender][category] = Dev(location, factory, category);

        devLogs[deviceId].push(Power(now, 0)); // 写入注册时间
        emit EvtRegisterDevice(location, factory, category);
    }

    function uploadLogs(uint256 consumption) public {
        // msg.sender is a device
        require(devLogs[msg.sender][0].createdAt > 0, "[Device](uploadLogs) device doesn't exist.");
        devLogs[msg.sender].push(Power(now, consumption));
    }

    function isDeviceAuthorized(address deviceOwner, bytes memory location, uint8 curStatus) public view returns(bool) {
        // curStatus是房屋的当前的可用状态
//        devInfos[deviceOwner]
        return false;
    }
}