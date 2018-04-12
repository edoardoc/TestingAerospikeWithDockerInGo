#!/bin/bash
cnt=0
echo cont = $cnt
for i in {1..25}
do
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]')"
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[11, 21], \"groupId\":22}, {\"interestIds\":[3, 45, 32], \"groupId\":3}]')"
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[23], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]')"
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[4, 5], \"groupId\":2}]')"
done

