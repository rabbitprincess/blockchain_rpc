// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

abstract contract Batch {
    string public constant name = "";
    string public constant symbol = "";
    uint8 public constant decimals = 0;

    function totalSupply() virtual public view returns (uint);
    function balanceOf(address tokenOwner) virtual public view returns (uint balance);
    function allowance(address tokenOwner, address spender) virtual public view returns (uint remaining);
    function transfer(address to, uint tokens) virtual public view returns (bool success);
    function approve(address spender, uint tokens) virtual public view returns (bool success);
    function transferFrom(address from, address to, uint tokens) virtual public view returns (bool success);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed tokenOwner, address indexed spender, uint256 value);

    function transferEth(address payable[] memory recipient, uint256 amount) public payable {
        for(uint i = 0; i < recipient.length; i++) {
            recipient[i].transfer(amount);
        }
    }
    function transferEthWithDifferentValue(address payable[] memory recipient, uint256[] memory amount) public payable {
        for(uint i = 0; i < recipient.length; i++) {
            recipient[i].transfer(amount[i]);
        }
    }
}
