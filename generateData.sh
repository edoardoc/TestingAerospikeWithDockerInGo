#!/bin/bash
echo "insert into cibucks.userProfiles (PK,key,profile) values(1,1,'JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]')"
# -- drop index cibucks.interestIds id_groupId
# CREATE MAPKEYS INDEX mpk_id_groupId ON cibucks.campaigns (profile) string
# CREATE MAPVALUES INDEX mpv_id_groupId on cibucks.campaigns (profile) numeric
# CREATE LIST INDEX id_groupId on cibucks.campaigns (profile) numeric

cnt=0
for ((i=0; i <= $1; i++));
do
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]')"
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[11, 21], \"groupId\":22}, {\"interestIds\":[3, 45, 32], \"groupId\":3}]')"
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[23], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]')"
    echo "insert into cibucks.campaigns (PK,key,profile) values($((++cnt)),$((cnt)),'JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[4, 5], \"groupId\":2}]')"
done

