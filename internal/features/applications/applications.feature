@apptests
Feature: Application Tests

  Scenario: Application Registration
    Given a developer registered with the portal
    And they have a new application they wish to register
    Then the application should be successfully registered

  Scenario: Retrieve Application Details
    Given a registed application
    Then the details assocaited with the application can be retrieved