@apptests
Feature: Application Tests

  Scenario: Application Registration
    Given a developer registered with the portal
    And they have a new application they wish to register
    Then the application should be successfully registered

  Scenario: Retrieve Application Details
    Given a registed application
    Then the details assocaited with the application can be retrieved

  Scenario: Duplicate Application Registration
    Given an application has already been registered
    And a developer attempts to register an application with the same name
    Then an error is returned with status code StatusConflict
    And the error message indicates a duplicate registration was attempted