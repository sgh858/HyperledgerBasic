## Running the test network

You can use the `./test-network/startchaincode.sh` script to stand up a simple Fabric test network. The test network has two peer organizations with one peer each and a single node raft ordering service. A simple chaincode (basicuser) is deployed on the created channel between two peer organizations.

This chaincode has following functions:
- InitLedger(): InitLedger adds a base set of users to the ledger (3 basic users).
- RegisterUser(ID, accessLevel): RegisterUser add a new user to the ledger with given details (Name as ID, Access Level as accessLevel).
- GetUser(ID): GetUser returns the user stored in the ledger with given id.
- UpdateUser(ID, accessLevel) : UpdateUser updates an existing user in the ledger with provided parameters. UpdateUser will add a new user if existing user is not found.
- DeleteUser(ID): DeleteUser deletes an given user from the ledger.
- GetAllUsers(): GetAllUsers returns all users found in the ledger.


Once you finished using the ledger, you shall use the command `./test-network/stopnetwork.sh` to stop the blockchain, delete the docker images and volumes. You may see docker displaying error while deleting image and volume. That's expected and not a cause for alarm.

Acknowledgement: This code is developed from examples in github.com/hyperledger/fabric-samples/.

