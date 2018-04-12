#!/bin/bash
cnt=12
echo cont = $cnt
echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),1,'JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]')"
echo 'insert into cibucks.campaigns (PK,key,profile) values(2,2,'JSON[{"interestIds":[11, 21], "groupId":22}, {"interestIds":[3, 45, 32], "groupId":3}]')'
echo 'insert into cibucks.campaigns (PK,key,profile) values(3,3,'JSON[{"interestIds":[23], "groupId":1}, {"interestIds":[3], "groupId":2}]')'
echo 'insert into cibucks.campaigns (PK,key,profile) values(4,4,'JSON[{"interestIds":[1, 2], "groupId":1}, {"interestIds":[4,5], "groupId":2}]')'
