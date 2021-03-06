## Running the HyperLedger Fabric test network with a basic chaincode and ABAC

This repo contains a simple Hyperledger Fabric chaincode used for basic user management. The blockchain contains 1 peer each from organization 1 (Org1) and organization 2 (Org2). A channel was created for these 2 peers in the blockchain. The chaincode was deployed on this channel. There is also an ordering service to support transaction ordering into blocks for distribution. All of them has their own Certificates with attributes for using ABAC (Attribute-Based Access Control).

Please see **Install.txt** for steps to install HyperLedger Fabric 2.2.

This chaincode was under directory HyperledgerBasic/user-basic. It's written in Go and has following functions:

- `InitLedger()`: InitLedger adds a base set of users to the ledger (3 basic users).
- `RegisterUser(ID, accessLevel)`: RegisterUser add a new user to the ledger with given details (Name as ID, Access Level as accessLevel).
- `GetUser(ID)`: GetUser returns the user stored in the ledger with given id.
- `UpdateUser(ID, accessLevel)` : UpdateUser updates an existing user in the ledger with provided parameters. UpdateUser will add a new user if existing user is not found.
- `DeleteUser(ID)`: DeleteUser deletes an given user from the ledger.
- `GetAllUsers()`: GetAllUsers returns all users found in the ledger.

To start the chaincode, please move to `HyperledgerBasic/test-network` directory before starting the blockchain. 
You can use the `.\startchaincode.sh` script to start up the blockchain network and deploy the chaincode.

1. cd `~/HyperledgerBasic/test-network`
2. Start the blockchain network and deploy the Chaincode using shellscript command `./startchaincode.sh`
3. It takes about 60sec for the blockchain to start and the Chaincode to be deployed. The chaincode initializes the ledger with 3 basic users. Afterthat, you can use the following shellscript to add/update/delete/modify users in the ledger.

   - `./getallusers.sh` 	: No parameter, display all user in the ledger
   - `./getusers.sh Name` 	: 1 parameter, display user with Name in the ledger
   - `./deleteuser.sh Name` 	: 1 parameter, delete user with Name in the ledger
   - `./registeruser.sh Name accessLevel` : 2 parameters, create a new user with given params in the ledger
   - `./updateuser.sh Authority Name accessLevel` : 3 parameters, update or create a new user using given authority (either *org1* or *org2*) and params  in the ledger. Need to use with `./approveuser.sh` and `./revokeuser.sh`
   - `./approveuser.sh Authority`: 1 parameter (either *org1* or *org2*), approve user of org1/org2 to use the chaincode `UpdateUser()` function. 
   - `./revokeuser.sh Authority`: 1 parameter (either *org1* or *org2*), revoke users of org1/org2 from using the chaincode `UpdateUser()` function.

   For example, for user from org1 to update the ledger, you should run the following command in order:
   - `./approveuser.sh org1		    // To authorize org1 users`
   - `./updateuser.sh org1 Dennis 10  // To update the ledger`
   - `./updateuser.sh org2 Dennis 20  // To return err as org2 is not authorised yet`
   - `./approveuser.sh org2		    // To authorize org2 users`
   - `./updateuser.sh org2 Dennis 20  // To update the ledger`

4. Once you finish the testing, please shutdown the blockchain using the following script
   - `./stopnetwork.sh`

This script will stop the blockchain network, delete the docker images and volumes. You may see docker displaying error while deleting images and volumes. That's expected and not a cause for alarm.

**Acknowledgement: This code is developed from examples in github.com/hyperledger/fabric-samples/.**

