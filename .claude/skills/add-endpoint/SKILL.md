# Add Endpoint Skill

When adding a new API endpoint:
1. Define the handler function in `handler/`
2. Define the business logic method on the relevant service interface
3. Implement the method in the service struct
4. Add the repository method to the repository interface if DB access is needed
5. Implement the repository method
6. Register the route in main.go
7. Write a unit test for the service method

Always follow the interface → implementation pattern. Never skip the interface step.
