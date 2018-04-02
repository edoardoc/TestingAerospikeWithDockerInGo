# Testing aerospike in GO

## Install aerospike with the following namespace in its configuration

```SQL
namespace cibucks {
        replication-factor 2
        memory-size 1G
        default-ttl 0d
        storage-engine device {
                file /opt/cibucks.dat
                filesize 1G
                data-in-memory true
        }
}
```

## Insert the following record

```SQL
insert into cibucks.userProfiles (PK,key,profile) values(1,1,'JSON[{"interestIds":[1, 2], "groupId":1}, {"interestIds":[3], "groupId":2}]')
```

## Insert the following records

```SQL
insert into cibucks.campaigns (PK,key,profile) values(1,1,'JSON[{"interestIds":[1, 2], "groupId":1}, {"interestIds":[3], "groupId":2}]')
insert into cibucks.campaigns (PK,key,profile) values(2,2,'JSON[{interestIds":[11, 21], "groupId":22}, {"interestIds":[3, 45, 32], "groupId":3}]')
insert into cibucks.campaigns (PK,key,profile) values(3,3,'JSON[{"interestIds":[23], "groupId":1}, {"interestIds":[3], "groupId":2}]')
insert into cibucks.campaigns (PK,key,profile) values(4,4,'JSON[{"interestIds":[1, 2], "groupId":1}, {"interestIds":[4,5], "groupId":2}]')
```

## Create a server in GoLang which expects a http parameter named ```userId``` and does the following

1. Gets the user data by querying aerospike on the set cibucks.userProfiles (for this exercise it will only work for userId =1 )
2. Queries the aerospike for all campaigns and checks which ones are valid for this user according to the following logic:
   * User should have all the groupIds that the campaign targets
   * User should have at least one same interestId per groupId   (In the dataset given only the first campaign is valid for this user)
   * Returns a list with the ids of the valid campaigns
