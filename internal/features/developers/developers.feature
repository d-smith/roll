@devtests
Feature: Developer Tests

  Scenario: Developer Registration
    Given A developer who registers on the portal
    And They have not registered before
    And Their data is formatted correctly
    Then They are added to the portal successfully

  Scenario: Malformed Email
    Given a developer who registers on the portal
    And They provide a malformed email
    Then An error is returned with StatusBadRequest