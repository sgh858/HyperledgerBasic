## Running the test network

You can use the `./startchaincode.sh` script to stand up a simple Fabric test network. The test network has two peer organizations with one peer each and a single node raft ordering service. A simple chaincode (basicuser) is deployed on the created channel between two peer organizations.

This chaincode has following functions:
- InitLedger(): InitLedger adds a base set of users to the ledger (3 basic users).
- RegisterUser(ID, accessLevel): RegisterUser add a new user to the ledger with given details (Name as ID, Access Level as accessLevel).
- GetUser(ID): GetUser returns the user stored in the ledger with given id.
- UpdateUser(ID, accessLevel) : UpdateUser updates an existing user in the ledger with provided parameters. UpdateUser will add a new user if existing user is not found.
- DeleteUser(ID): DeleteUser deletes an given user from the ledger.
- GetAllUsers(): GetAllUsers returns all users found in the ledger.



