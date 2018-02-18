pragma solidity ^0.4.17;
contract InBox {
    string public message;
    function InBox(string initialMessage) public {
        message = initialMessage; // same name as contract so constructor
    }
    function setMessage(string newMessage) public {
        message = newMessage;
    }
    function getMessage() public view returns (string){
        return message;
    }
    
}
